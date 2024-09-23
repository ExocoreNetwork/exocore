package keeper

import (
	"strings"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/delegation/types"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EndBlock : completed Undelegation events according to the canCompleted blockHeight
// This function will be triggered at the end of every block, it will query the undelegation state to get the records that need to be handled and try to complete the undelegation task.
func (k *Keeper) EndBlock(oCtx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// #nosec G703 // the error is always nil
	records, _ := k.GetPendingUndelegationRecords(oCtx, uint64(oCtx.BlockHeight()))
	if len(records) == 0 {
		return []abci.ValidatorUpdate{}
	}
	for _, record := range records {
		ctx, writeCache := oCtx.CacheContext()
		// check if the operator has been slashed or frozen
		operatorAccAddress := sdk.MustAccAddressFromBech32(record.OperatorAddr)
		// todo: don't think about freezing the operator in current implementation
		/*		if k.slashKeeper.IsOperatorFrozen(ctx, operatorAccAddress) {
				// reSet the completed height if the operator is frozen
				record.CompleteBlockNumber = k.operatorKeeper.GetUnbondingExpirationBlockNumber(ctx, operatorAccAddress, record.BlockNumber)
				if record.CompleteBlockNumber <= uint64(ctx.BlockHeight()) {
					panic(fmt.Sprintf("the reset completedHeight isn't in future,setHeight:%v,curHeight:%v", record.CompleteBlockNumber, ctx.BlockHeight()))
				}
				_, innerError = k.SetSingleUndelegationRecord(ctx, record)
				if innerError != nil {
					panic(innerError)
				}
				continue
			}*/

		recordID := types.GetUndelegationRecordKey(record.BlockNumber, record.LzTxNonce, record.TxHash, record.OperatorAddr)
		if k.GetUndelegationHoldCount(ctx, recordID) > 0 {
			// store it again with the next block and move on
			// #nosec G701
			record.CompleteBlockNumber = uint64(ctx.BlockHeight()) + 1
			// we need to store two things here: one is the updated record in itself
			// #nosec G703 // the error is always nil
			recordKey, _ := k.SetSingleUndelegationRecord(ctx, record)
			// and the other is the fact that it matures at the next block
			// #nosec G703 // the error is always nil
			_ = k.StorePendingUndelegationRecord(ctx, recordKey, record)
			writeCache()
			continue
		}
		// TODO(mike): ensure that operator is required to perform self delegation to match above.

		recordAmountNeg := record.Amount.Neg()
		// update delegation state
		deltaAmount := &types.DeltaDelegationAmounts{
			WaitUndelegationAmount: recordAmountNeg,
		}
		_, err := k.UpdateDelegationState(ctx, record.StakerID, record.AssetID, record.OperatorAddr, deltaAmount)
		if err != nil {
			// use oCtx so that the error is logged on the original context
			k.Logger(oCtx).Error("failed to update delegation state", "error", err)
			continue
		}

		// update the staker state
		if record.AssetID == assetstypes.NativeAssetID {
			parsedStakerID := strings.Split(record.StakerID, "_")
			stakerAddr := sdk.AccAddress(hexutil.MustDecode(parsedStakerID[0]))
			if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(ctx, types.DelegatedPoolName, stakerAddr, sdk.NewCoins(sdk.NewCoin(assetstypes.NativeAssetDenom, record.ActualCompletedAmount))); err != nil {
				k.Logger(oCtx).Error("failed to undelegate coins from module to account", "error", err)
				continue
			}
		} else {
			err = k.assetsKeeper.UpdateStakerAssetState(ctx, record.StakerID, record.AssetID, assetstypes.DeltaStakerSingleAsset{
				WithdrawableAmount:        record.ActualCompletedAmount,
				PendingUndelegationAmount: recordAmountNeg,
			})
			if err != nil {
				k.Logger(oCtx).Error("failed to update staker asset state", "error", err)
				continue
			}
		}

		// update the operator state
		err = k.assetsKeeper.UpdateOperatorAssetState(ctx, operatorAccAddress, record.AssetID, assetstypes.DeltaOperatorSingleAsset{
			PendingUndelegationAmount: recordAmountNeg,
		})
		if err != nil {
			k.Logger(oCtx).Error("failed to update operator asset state", "error", err)
			continue
		}

		// delete the Undelegation records that have been complemented
		// #nosec G703 // the error is always nil
		_ = k.DeleteUndelegationRecord(ctx, record)
		// when calling `writeCache`, events are automatically emitted on the parent context
		writeCache()
	}
	return []abci.ValidatorUpdate{}
}
