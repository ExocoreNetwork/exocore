package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ClearUnbondingInformation clears all information related to an operator's opt out
// or key replacement. This is done because the operator has opted back in or has
// replaced their key (again) with the original one.
func (k Keeper) ClearUnbondingInformation(
	ctx sdk.Context, addr sdk.AccAddress, pubKey *tmprotocrypto.PublicKey,
) {
	optOutEpoch := k.GetOperatorOptOutFinishEpoch(ctx, addr)
	k.DeleteOperatorOptOutFinishEpoch(ctx, addr)
	k.RemoveOptOutToFinish(ctx, optOutEpoch, addr)
	consAddress, err := operatortypes.TMCryptoPublicKeyToConsAddr(pubKey)
	if err != nil {
		return
	}
	k.DeleteConsensusAddrToPrune(ctx, optOutEpoch, consAddress)
}

// SetUnbondingInformation sets information related to an operator's opt out or key replacement.
func (k Keeper) SetUnbondingInformation(
	ctx sdk.Context, addr sdk.AccAddress, pubKey *tmprotocrypto.PublicKey, isOptingOut bool,
) {
	unbondingCompletionEpoch := k.GetUnbondingCompletionEpoch(ctx)
	if isOptingOut {
		k.AppendOptOutToFinish(ctx, unbondingCompletionEpoch, addr)
		k.SetOperatorOptOutFinishEpoch(ctx, addr, unbondingCompletionEpoch)
	}
	consAddress, err := operatortypes.TMCryptoPublicKeyToConsAddr(pubKey)
	if err != nil {
		return
	}
	k.AppendConsensusAddrToPrune(ctx, unbondingCompletionEpoch, consAddress)
}

// GetUnbondingCompletionEpoch returns the epoch at the end of which
// an unbonding triggered in this epoch will be completed.
func (k Keeper) GetUnbondingCompletionEpoch(
	ctx sdk.Context,
) int64 {
	unbondingEpochs := k.GetEpochsUntilUnbonded(ctx)
	epochInfo, _ := k.epochsKeeper.GetEpochInfo(
		ctx, k.GetEpochIdentifier(ctx),
	)
	// if i execute the transaction at epoch 5, the vote power change
	// goes into effect at the beginning of epoch 6. the information
	// should be held for 7 epochs, so it should be deleted at the
	// beginning of epoch 13 or the end of epoch 12.
	return epochInfo.CurrentEpoch + int64(unbondingEpochs) // #nosec G701
}

// AppendUndelegationsToMature stores that the undelegation with recordKey should be
// released at the end of the provided epoch.
func (k Keeper) AppendUndelegationToMature(
	ctx sdk.Context, epoch int64, recordKey []byte,
) {
	prev := k.GetUndelegationsToMature(ctx, epoch)
	next := types.UndelegationRecordKeys{
		List: append(prev, recordKey),
	}
	k.setUndelegationsToMature(ctx, epoch, next)
}

// GetUndelegationsToMature returns all undelegation entries that should be released
// at the end of the provided epoch.
func (k Keeper) GetUndelegationsToMature(
	ctx sdk.Context, epoch int64,
) [][]byte {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.UnbondingReleaseMaturityKey(epoch)
	bz := store.Get(key)
	if bz == nil {
		return [][]byte{}
	}
	var res types.UndelegationRecordKeys
	if err := res.Unmarshal(bz); err != nil {
		// should never happen
		panic(err)
	}
	return res.GetList()
}

// ClearUndelegationsToMature is a pruning method which is called after we mature
// the undelegation entries.
func (k Keeper) ClearUndelegationsToMature(
	ctx sdk.Context, epoch int64,
) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.UnbondingReleaseMaturityKey(epoch)
	store.Delete(key)
}

// setUndelegationsToMature sets all undelegation entries that should be released
// at the end of the provided epoch.
func (k Keeper) setUndelegationsToMature(
	ctx sdk.Context, epoch int64, undelegationRecords types.UndelegationRecordKeys,
) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.UnbondingReleaseMaturityKey(epoch)
	val, err := undelegationRecords.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, val)
}
