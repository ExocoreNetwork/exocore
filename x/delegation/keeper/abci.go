package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	delegationtype "github.com/exocore/x/delegation/types"
	"github.com/exocore/x/restaking_assets_manage/types"
)

// EndBlock : completed Undelegation events according to the canCompleted blockHeight
func (k Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	ctx.Logger().Info("the blockHeight is:", "height", ctx.BlockHeight())
	records, err := k.GetWaitCompleteUndelegationRecords(ctx, uint64(ctx.BlockHeight()))
	if err != nil {
		panic(err)
	}
	if len(records) == 0 {
		return []abci.ValidatorUpdate{}
	}
	for _, record := range records {
		// check if the operator has been slashed or frozen
		operatorAccAddress := sdk.MustAccAddressFromBech32(record.OperatorAddr)
		if k.slashKeeper.IsOperatorFrozen(ctx, operatorAccAddress) {
			// reSet the completed height if the operator is frozen
			record.CompleteBlockNumber = k.operatorOptedInKeeper.GetOperatorCanUndelegateHeight(ctx, record.AssetId, operatorAccAddress, record.BlockNumber)
			if record.CompleteBlockNumber <= uint64(ctx.BlockHeight()) {
				panic(fmt.Sprintf("the reset completedHeight isn't in future,setHeight:%v,curHeight:%v", record.CompleteBlockNumber, ctx.BlockHeight()))
			}
			_, err = k.SetSingleUndelegationRecord(ctx, record)
			if err != nil {
				panic(err)
			}
			continue
		}

		// get operator slashed proportion to calculate the actual canUndelegated asset amount
		proportion := k.slashKeeper.OperatorAssetSlashedProportion(ctx, operatorAccAddress, record.AssetId, record.BlockNumber, record.CompleteBlockNumber)
		if proportion.IsNil() || proportion.IsNegative() || proportion.GT(sdkmath.LegacyNewDec(1)) {
			panic(fmt.Sprintf("the proportion is invalid,it is:%v", proportion))
		}
		canUndelegateProportion := sdkmath.LegacyNewDec(1).Sub(proportion)
		actualCanUndelegateAmount := canUndelegateProportion.MulInt(record.Amount).TruncateInt()
		record.ActualCompletedAmount = actualCanUndelegateAmount
		recordAmountNeg := record.Amount.Neg()

		// update delegation state
		delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
		delegatorAndAmount[record.OperatorAddr] = &delegationtype.DelegationAmounts{
			WaitUndelegationAmount: recordAmountNeg,
		}
		err = k.UpdateDelegationState(ctx, record.StakerId, record.AssetId, delegatorAndAmount)
		if err != nil {
			panic(err)
		}

		// todo: if use recordAmount as an input parameter, the delegation total amount won't need to be subtracted when the related operator is slashed.
		err = k.UpdateStakerDelegationTotalAmount(ctx, record.StakerId, record.AssetId, recordAmountNeg)
		if err != nil {
			panic(err)
		}

		// update the staker state
		err := k.restakingStateKeeper.UpdateStakerAssetState(ctx, record.StakerId, record.AssetId, types.StakerSingleAssetOrChangeInfo{
			CanWithdrawAmountOrWantChangeValue:      actualCanUndelegateAmount,
			WaitUndelegationAmountOrWantChangeValue: recordAmountNeg,
		})
		if err != nil {
			panic(err)
		}

		// update the operator state
		err = k.restakingStateKeeper.UpdateOperatorAssetState(ctx, operatorAccAddress, record.AssetId, types.OperatorSingleAssetOrChangeInfo{
			TotalAmountOrWantChangeValue:            actualCanUndelegateAmount.Neg(),
			WaitUndelegationAmountOrWantChangeValue: recordAmountNeg,
		})
		if err != nil {
			panic(err)
		}

		// update Undelegation record
		record.IsPending = false
		_, err = k.SetSingleUndelegationRecord(ctx, record)
		if err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}
