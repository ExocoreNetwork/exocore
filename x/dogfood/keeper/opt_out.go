package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AppendOptOutToFinish appends an operator address to the list of operator addresses that have
// opted out and will be finished at the end of the provided epoch.
func (k Keeper) AppendOptOutToFinish(
	ctx sdk.Context, epoch int64, operatorAddr sdk.AccAddress,
) {
	prev := k.GetOptOutsToFinish(ctx, epoch)
	next := types.AccountAddresses{List: append(prev, operatorAddr)}
	k.setOptOutsToFinish(ctx, epoch, next)
}

// GetOptOutsToFinish returns the list of operator addresses that have opted out and will be
// finished at the end of the provided epoch.
func (k Keeper) GetOptOutsToFinish(
	ctx sdk.Context, epoch int64,
) [][]byte {
	store := ctx.KVStore(k.storeKey)
	// the epochs module validates at genesis that epoch is non-negative.
	key, _ := types.OptOutsToFinishKey(epoch)
	bz := store.Get(key)
	if bz == nil {
		return [][]byte{}
	}
	var res types.AccountAddresses
	if err := res.Unmarshal(bz); err != nil {
		panic(err)
	}
	return res.GetList()
}

// setOptOutsToFinish sets the list of operator addresses that have opted out and will be
// finished at the end of the provided epoch.
func (k Keeper) setOptOutsToFinish(
	ctx sdk.Context, epoch int64, addrs types.AccountAddresses,
) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.OptOutsToFinishKey(epoch)
	bz, err := addrs.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

// ClearOptOutsToFinish clears the list of operator addresses that have opted out and will be
// finished at the end of the provided epoch.
func (k Keeper) ClearOptOutsToFinish(ctx sdk.Context, epoch int64) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.OptOutsToFinishKey(epoch)
	store.Delete(key)
}

// SetOperatorOptOutFinishEpoch sets the epoch at which an operator's opt out will be finished.
func (k Keeper) SetOperatorOptOutFinishEpoch(
	ctx sdk.Context, operatorAddr sdk.AccAddress, epoch int64,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.OperatorOptOutFinishEpochKey(operatorAddr)
	uepoch, _ := types.SafeInt64ToUint64(epoch)
	bz := sdk.Uint64ToBigEndian(uepoch)
	store.Set(key, bz)
}

// GetOperatorOptOutFinishEpoch returns the epoch at which an operator's opt out will be
// finished.
func (k Keeper) GetOperatorOptOutFinishEpoch(
	ctx sdk.Context, operatorAddr sdk.AccAddress,
) int64 {
	store := ctx.KVStore(k.storeKey)
	key := types.OperatorOptOutFinishEpochKey(operatorAddr)
	bz := store.Get(key)
	if bz == nil {
		return -1
	}
	// max int64 is 9 quintillion, and max uint64 is double of that.
	// it is too far in the future to be a concern.
	return int64(sdk.BigEndianToUint64(bz)) // #nosec G701 // see above.
}

// DeleteOperatorOptOutFinishEpoch deletes the epoch at which an operator's opt out will be
// finished.
func (k Keeper) DeleteOperatorOptOutFinishEpoch(
	ctx sdk.Context, operatorAddr sdk.AccAddress,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.OperatorOptOutFinishEpochKey(operatorAddr)
	store.Delete(key)
}

// AppendConsensusAddrToPrune appends a consensus address to the list of consensus addresses to
// prune at the end of the epoch.
func (k Keeper) AppendConsensusAddrToPrune(
	ctx sdk.Context, epoch int64, operatorAddr sdk.ConsAddress,
) {
	prev := k.GetConsensusAddrsToPrune(ctx, epoch)
	next := types.ConsensusAddresses{List: append(prev, operatorAddr)}
	k.setConsensusAddrsToPrune(ctx, epoch, next)
}

// GetConsensusAddrsToPrune returns the list of consensus addresses to prune at the end of the
// epoch.
func (k Keeper) GetConsensusAddrsToPrune(
	ctx sdk.Context, epoch int64,
) [][]byte {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.ConsensusAddrsToPruneKey(epoch)
	bz := store.Get(key)
	if bz == nil {
		return [][]byte{}
	}
	var res types.ConsensusAddresses
	if err := res.Unmarshal(bz); err != nil {
		panic(err)
	}
	return res.GetList()
}

// ClearConsensusAddrsToPrune clears the list of consensus addresses to prune at the end of the
// epoch.
func (k Keeper) ClearConsensusAddrsToPrune(ctx sdk.Context, epoch int64) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.ConsensusAddrsToPruneKey(epoch)
	store.Delete(key)
}

// setConsensusAddrsToPrune sets the list of consensus addresses to prune at the end of the
// epoch.
func (k Keeper) setConsensusAddrsToPrune(
	ctx sdk.Context, epoch int64, addrs types.ConsensusAddresses,
) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.ConsensusAddrsToPruneKey(epoch)
	bz, err := addrs.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}
