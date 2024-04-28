package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetUnbondingInformation sets information related to an operator's opt out or key replacement.
func (k Keeper) SetUnbondingInformation(
	ctx sdk.Context, addr sdk.AccAddress, pubKey *tmprotocrypto.PublicKey,
) {
	unbondingCompletionEpoch := k.GetUnbondingCompletionEpoch(ctx)
	k.AppendOptOutToFinish(ctx, unbondingCompletionEpoch, addr)
	k.SetOperatorOptOutFinishEpoch(ctx, addr, unbondingCompletionEpoch)
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

// GetUndelegationMaturityEpoch gets the maturity epoch for the undelegation record.
func (k Keeper) GetUndelegationMaturityEpoch(
	ctx sdk.Context, recordKey []byte,
) (int64, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.UndelegationMaturityEpochKey(recordKey)
	bz := store.Get(key)
	if bz == nil {
		return 0, false
	}
	epoch := sdk.BigEndianToUint64(bz)
	return types.SafeUint64ToInt64(epoch)
}

// SetUndelegationMaturityEpoch sets the maturity epoch for the undelegation record.
func (k Keeper) SetUndelegationMaturityEpoch(
	ctx sdk.Context, recordKey []byte, epoch int64,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.UndelegationMaturityEpochKey(recordKey)
	uepoch, _ := types.SafeInt64ToUint64(epoch)
	bz := sdk.Uint64ToBigEndian(uepoch)
	store.Set(key, bz)
}

// ClearUndelegationMaturityEpoch clears the maturity epoch for the undelegation record.
func (k Keeper) ClearUndelegationMaturityEpoch(
	ctx sdk.Context, recordKey []byte,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.UndelegationMaturityEpochKey(recordKey)
	store.Delete(key)
}
