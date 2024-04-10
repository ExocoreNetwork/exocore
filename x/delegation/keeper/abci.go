package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
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
		// todo: don't think about freezing the operator in current implementation
		/*		if k.slashKeeper.IsOperatorFrozen(ctx, operatorAccAddress) {
				// reSet the completed height if the operator is frozen
				record.CompleteBlockNumber = k.operatorKeeper.GetUnbondingExpirationBlockNumber(ctx, operatorAccAddress, record.BlockNumber)
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
			// #nosec G701
			record.CompleteBlockNumber = uint64(ctx.BlockHeight()) + 1
			// we need to store two things here: one is the updated record in itself
			recordKey, err := k.SetSingleUndelegationRecord(ctx, record)
			if err != nil {
				panic(err)
			}
			// and the other is the fact that it matures at the next block
			err = k.StoreWaitCompleteRecord(ctx, recordKey, record)
			if err != nil {
				panic(err)
			}
			continue
		}
		// TODO(mike): ensure that operator is required to perform self delegation to match above.

		// calculate the actual canUndelegated asset amount
		delegationInfo, err := k.GetSingleDelegationInfo(ctx, record.StakerID, record.AssetID, record.OperatorAddr)
		if err != nil {
			panic(err)
		}
		if record.Amount.GT(delegationInfo.UndelegatableAfterSlash) {
			record.ActualCompletedAmount = delegationInfo.UndelegatableAfterSlash
		} else {
			record.ActualCompletedAmount = record.Amount
		}
		recordAmountNeg := record.Amount.Neg()

		// update delegation state
		delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts)
		delegatorAndAmount[record.OperatorAddr] = &delegationtype.DelegationAmounts{
			WaitUndelegationAmount:  recordAmountNeg,
			UndelegatableAfterSlash: record.ActualCompletedAmount.Neg(),
		}
		err = k.UpdateDelegationState(ctx, record.StakerID, record.AssetID, delegatorAndAmount)
		if err != nil {
			panic(err)
		}

		// update the staker state
		err = k.assetsKeeper.UpdateStakerAssetState(ctx, record.StakerID, record.AssetID, types.DeltaStakerSingleAsset{
			WithdrawableAmount:  record.ActualCompletedAmount,
			WaitUnbondingAmount: recordAmountNeg,
		})
		if err != nil {
			panic(err)
		}

		// update the operator state
		err = k.assetsKeeper.UpdateOperatorAssetState(ctx, operatorAccAddress, record.AssetID, types.DeltaOperatorSingleAsset{
			WaitUnbondingAmount: recordAmountNeg,
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
