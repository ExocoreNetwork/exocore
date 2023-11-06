// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package keeper

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/delegation/types"
	"github.com/exocore/x/restaking_assets_manage/types"
)

// EndBlock : completed unDelegation events according to the canCompleted blockHeight
func (k Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	records, err := k.GetWaitCompleteUnDelegationRecords(ctx, uint64(ctx.BlockHeight()))
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
			//reSet the completed height if the operator is frozen
			record.CompleteBlockNumber = k.operatorOptedInKeeper.GetOperatorCanUnDelegateHeight(ctx, record.AssetId, operatorAccAddress, record.BlockNumber)
			if record.CompleteBlockNumber <= uint64(ctx.BlockHeight()) {
				panic(fmt.Sprintf("the reset completedHeight isn't in future,setHeight:%v,curHeight:%v", record.CompleteBlockNumber, ctx.BlockHeight()))
			}
			_, err = k.SetSingleUnDelegationRecord(ctx, record)
			if err != nil {
				panic(err)
			}
			continue
		}

		//get operator slashed proportion to calculate the actual canUnDelegated asset amount
		proportion := k.slashKeeper.OperatorAssetSlashedProportion(ctx, operatorAccAddress, record.AssetId, record.BlockNumber, record.CompleteBlockNumber)
		if proportion.IsNil() || proportion.IsNegative() || proportion.GT(sdkmath.LegacyNewDec(1)) {
			panic(fmt.Sprintf("the proportion is invalid,it is:%v", proportion))
		}
		actualCanUnDelegateAmount := proportion.MulInt(record.Amount.Amount).TruncateInt()
		record.ActualCompletedAmount.Amount = actualCanUnDelegateAmount
		recordAmountNeg := record.Amount.Amount.Neg()

		//update delegation state
		delegatorAndAmount := make(map[string]*types2.DelegationAmounts)
		delegatorAndAmount[record.OperatorAddr] = &types2.DelegationAmounts{
			WaitUnDelegationAmount: &types2.ValueField{Amount: recordAmountNeg},
		}
		err = k.UpdateDelegationState(ctx, record.StakerId, record.AssetId, delegatorAndAmount)
		if err != nil {
			panic(err)
		}

		//todo: if use recordAmount as an input parameter, the delegation total amount won't need to be subtracted when the related operator is slashed.
		err = k.UpdateStakerDelegationTotalAmount(ctx, record.StakerId, record.AssetId, recordAmountNeg)
		if err != nil {
			panic(err)
		}

		//update the staker state
		err := k.retakingStateKeeper.UpdateStakerAssetState(ctx, record.StakerId, record.AssetId, types.StakerSingleAssetOrChangeInfo{
			CanWithdrawAmountOrWantChangeValue:      actualCanUnDelegateAmount,
			WaitUnDelegationAmountOrWantChangeValue: recordAmountNeg,
		})
		if err != nil {
			panic(err)
		}

		//update the operator state
		err = k.retakingStateKeeper.UpdateOperatorAssetState(ctx, operatorAccAddress, record.AssetId, types.OperatorSingleAssetOrChangeInfo{
			TotalAmountOrWantChangeValue:            actualCanUnDelegateAmount.Neg(),
			WaitUnDelegationAmountOrWantChangeValue: recordAmountNeg,
		})
		if err != nil {
			panic(err)
		}

		//update unDelegation record
		record.IsPending = false
		_, err = k.SetSingleUnDelegationRecord(ctx, record)
		if err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}
