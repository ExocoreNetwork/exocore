package keeper

import (
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlock : completed Undelegation events according to the canCompleted blockHeight
// This function will be triggered at the end of every block,it will query the undelegation state to get the records that need to be handled and try to complete the undelegation task.
func (k *Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
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
				record.CompleteBlockNumber = k.expectedOperatorInterface.GetUnbondingExpirationBlockNumber(ctx, operatorAccAddress, record.BlockNumber)
				if record.CompleteBlockNumber <= uint64(ctx.BlockHeight()) {
					panic(fmt.Sprintf("the reset completedHeight isn't in future,setHeight:%v,curHeight:%v", record.CompleteBlockNumber, ctx.BlockHeight()))
				}
				_, err = k.SetSingleUndelegationRecord(ctx, record)
				if err != nil {
					panic(err)
				}
				continue
			}*/

		recordID := delegationtype.GetUndelegationRecordKey(record.LzTxNonce, record.TxHash, record.OperatorAddr)
		if k.GetUndelegationHoldCount(ctx, recordID) > 0 {
			// store it again with the next block and move on
			record.CompleteBlockNumber = uint64(ctx.BlockHeight()) + 1
			// we need to store two things here: one is the updated record in itself
			recordKey, err := k.SetSingleUndelegationRecord(ctx, record)
			if err != nil {
				panic(err)
			}
			// and the other is the fact that it matures at the next block
			k.StoreWaitCompleteRecord(ctx, recordKey, record)
			continue
		}
		// operator opt out: since operators can not immediately withdraw their funds, that is,
		// even operator funds are not immediately available, operator opt out does not require
		// any special handling here. if an operator undelegates before they opt out, the undelegation
		// will be processed normally. if they undelegate after they opt out, the undelegation will
		// be released at the same time as opt out completion, provided there are no other chains that
		// the operator is still active on. the same applies to delegators too.
		// TODO(mike): ensure that operator is required to perform self delegation to match above.

		//calculate the actual canUndelegated asset amount
		delegationInfo, err := k.GetSingleDelegationInfo(ctx, record.StakerID, record.AssetID, record.OperatorAddr)
		if err != nil {
			panic(err)
		}
		if record.Amount.GT(delegationInfo.UndelegatableAmountAfterSlash) {
			record.ActualCompletedAmount = delegationInfo.UndelegatableAmountAfterSlash
		} else {
			record.ActualCompletedAmount = record.Amount
		}
		recordAmountNeg := record.Amount.Neg()

		// update delegation state
		delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
		delegatorAndAmount[record.OperatorAddr] = &delegationtype.DelegationAmounts{
			WaitUndelegationAmount:        recordAmountNeg,
			UndelegatableAmountAfterSlash: record.ActualCompletedAmount.Neg(),
		}
		err = k.UpdateDelegationState(ctx, record.StakerID, record.AssetID, delegatorAndAmount)
		if err != nil {
			panic(err)
		}

		// update the staker state
		err = k.restakingStateKeeper.UpdateStakerAssetState(ctx, record.StakerID, record.AssetID, types.StakerSingleAssetOrChangeInfo{
			CanWithdrawAmountOrWantChangeValue:   record.ActualCompletedAmount,
			WaitUnbondingAmountOrWantChangeValue: recordAmountNeg,
		})
		if err != nil {
			panic(err)
		}

		// update the operator state
		err = k.restakingStateKeeper.UpdateOperatorAssetState(ctx, operatorAccAddress, record.AssetID, types.OperatorSingleAssetOrChangeInfo{
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
