package keeper

import (
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/delegation/types"

	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// EndBlock : completed Undelegation events according to the canCompleted blockHeight
// This function will be triggered at the end of every block, it will query the undelegation state to get the records that need to be handled and try to complete the undelegation task.
func (k *Keeper) EndBlock(
	originalCtx sdk.Context, _ abci.RequestEndBlock,
) []abci.ValidatorUpdate {
	logger := k.Logger(originalCtx)
	records, err := k.GetPendingUndelegationRecords(
		originalCtx, uint64(originalCtx.BlockHeight()),
	)
	if err != nil {
		logger.Error("failed to get pending undelegation records", "error", err)
	}
	if len(records) == 0 {
		logger.Info("no pending undelegation records")
		return []abci.ValidatorUpdate{}
	}
	for i := range records {
		record := records[i] // avoid implicit memory aliasing
		ctx, writeCache := originalCtx.CacheContext()
		// we can use `Must` here because we stored this record ourselves.
		operatorAccAddress := sdk.MustAccAddressFromBech32(record.OperatorAddr)
		// TODO check if the operator has been slashed or frozen
		recordID := types.GetUndelegationRecordKey(
			record.BlockNumber, record.LzTxNonce, record.TxHash, record.OperatorAddr,
		)
		if k.GetUndelegationHoldCount(ctx, recordID) > 0 {
			// delete from all 3 states
			if err := k.DeleteUndelegationRecord(ctx, record); err != nil {
				logger.Error("failed to delete undelegation record", "error", err)
				continue
			}
			// add back to all 3 states, with the new block height
			// #nosec G701
			record.CompleteBlockNumber = uint64(ctx.BlockHeight()) + 1
			if err := k.SetUndelegationRecords(
				ctx, []types.UndelegationRecord{*record},
			); err != nil {
				logger.Error("failed to set undelegation records", "error", err)
				continue
			}
			writeCache()
			continue
		}

		recordAmountNeg := record.Amount.Neg()
		// update delegation state
		deltaAmount := &types.DeltaDelegationAmounts{
			WaitUndelegationAmount: recordAmountNeg,
		}
		if _, err := k.UpdateDelegationState(
			ctx, record.StakerID, record.AssetID, record.OperatorAddr, deltaAmount,
		); err != nil {
			logger.Error(
				"failed to update delegation state",
				"error", err,
			)
			continue
		}

		// update the staker state
		if record.AssetID == assetstypes.ExocoreAssetID {
			stakerAddrHex, _, err := assetstypes.ParseID(record.StakerID)
			if err != nil {
				logger.Error(
					"failed to parse staker ID",
					"error", err,
				)
				continue
			}
			stakerAddrBytes, err := hexutil.Decode(stakerAddrHex)
			if err != nil {
				logger.Error(
					"failed to decode staker address",
					"error", err,
				)
				continue
			}
			stakerAddr := sdk.AccAddress(stakerAddrBytes)
			if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
				ctx, types.DelegatedPoolName, stakerAddr,
				sdk.NewCoins(
					sdk.NewCoin(assetstypes.ExocoreAssetDenom, record.ActualCompletedAmount),
				),
			); err != nil {
				logger.Error(
					"failed to undelegate coins from module to account",
					"error", err,
				)
				continue
			}
		} else {
			if err := k.assetsKeeper.UpdateStakerAssetState(
				ctx, record.StakerID, record.AssetID,
				assetstypes.DeltaStakerSingleAsset{
					WithdrawableAmount:        record.ActualCompletedAmount,
					PendingUndelegationAmount: recordAmountNeg,
				},
			); err != nil {
				logger.Error(
					"failed to update staker asset state",
					"error", err,
				)
				continue
			}
		}

		// update the operator state
		if err := k.assetsKeeper.UpdateOperatorAssetState(
			ctx, operatorAccAddress, record.AssetID,
			assetstypes.DeltaOperatorSingleAsset{
				PendingUndelegationAmount: recordAmountNeg,
			},
		); err != nil {
			logger.Error(
				"failed to update operator asset state",
				"error", err,
			)
			continue
		}

		// delete the Undelegation records that have been complemented
		if err := k.DeleteUndelegationRecord(ctx, record); err != nil {
			logger.Error(
				"failed to delete undelegation record",
				"error", err,
			)
			continue
		}
		// when calling `writeCache`, events are automatically emitted on the parent context
		writeCache()
	}
	return []abci.ValidatorUpdate{}
}
