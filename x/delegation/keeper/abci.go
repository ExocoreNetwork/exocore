package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
)

// EndBlock : completed Undelegation events according to the canCompleted blockHeight
// This function will be triggered at the end of every block,it will query the undelegation state to get the records that need to be handled and try to complete the undelegation task.
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
		//todo: don't think about freezing the operator in current implementation
		/*		if k.slashKeeper.IsOperatorFrozen(ctx, operatorAccAddress) {
				// reSet the completed height if the operator is frozen
				record.CompleteBlockNumber = k.expectOperatorInterface.GetUnBondingExpirationBlockNumber(ctx, operatorAccAddress, record.BlockNumber)
				if record.CompleteBlockNumber <= uint64(ctx.BlockHeight()) {
					panic(fmt.Sprintf("the reset completedHeight isn't in future,setHeight:%v,curHeight:%v", record.CompleteBlockNumber, ctx.BlockHeight()))
				}
				_, err = k.SetSingleUndelegationRecord(ctx, record)
				if err != nil {
					panic(err)
				}
				continue
			}*/

		//calculate the actual canUndelegated asset amount
		delegationInfo, err := k.GetSingleDelegationInfo(ctx, record.StakerId, record.AssetId, record.OperatorAddr)
		if record.Amount.GT(delegationInfo.CanUndelegateAmountAfterSlash) {
			record.ActualCompletedAmount = delegationInfo.CanUndelegateAmountAfterSlash
		} else {
			record.ActualCompletedAmount = record.Amount
		}
		recordAmountNeg := record.Amount.Neg()

		// update delegation state
		delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
		delegatorAndAmount[record.OperatorAddr] = &delegationtype.DelegationAmounts{
			WaitUndelegationAmount:        recordAmountNeg,
			CanUndelegateAmountAfterSlash: record.ActualCompletedAmount.Neg(),
		}
		err = k.UpdateDelegationState(ctx, record.StakerId, record.AssetId, delegatorAndAmount)
		if err != nil {
			panic(err)
		}

		// update the staker state
		err = k.restakingStateKeeper.UpdateStakerAssetState(ctx, record.StakerId, record.AssetId, types.StakerSingleAssetOrChangeInfo{
			CanWithdrawAmountOrWantChangeValue:   record.ActualCompletedAmount,
			WaitUnbondingAmountOrWantChangeValue: recordAmountNeg,
		})
		if err != nil {
			panic(err)
		}

		// update the operator state
		err = k.restakingStateKeeper.UpdateOperatorAssetState(ctx, operatorAccAddress, record.AssetId, types.OperatorSingleAssetOrChangeInfo{
			WaitUnbondingAmountOrWantChangeValue: recordAmountNeg,
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
