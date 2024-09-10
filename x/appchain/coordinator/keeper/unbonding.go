package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AppendConsAddrToPrune appends a consensus address to the list of consensus addresses
// that will be pruned when the validator set update containing the given vscID is matured
// by the chainID.
func (k Keeper) AppendConsAddrToPrune(
	ctx sdk.Context, chainID string, vscID uint64, consKey sdk.ConsAddress,
) {
	prev := k.GetConsAddrsToPrune(ctx, chainID, vscID)
	prev.List = append(prev.List, consKey)
	k.SetConsAddrsToPrune(ctx, chainID, vscID, prev)
}

// GetConsAddrsToPrune returns the list of consensus addresses that will be pruned when the
// validator set update containing the given vscID is matured by the chainID.
func (k Keeper) GetConsAddrsToPrune(
	ctx sdk.Context, chainID string, vscID uint64,
) (res types.ConsensusAddresses) {
	store := ctx.KVStore(k.storeKey)
	key := types.ConsAddrsToPruneKey(chainID, vscID)
	k.cdc.MustUnmarshal(store.Get(key), &res)
	return res
}

// SetConsAddrsToPrune sets the list of consensus addresses that will be pruned when the
// validator set update containing the given vscID is matured by the chainID.
func (k Keeper) SetConsAddrsToPrune(
	ctx sdk.Context, chainID string, vscID uint64, addrs types.ConsensusAddresses,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.ConsAddrsToPruneKey(chainID, vscID)
	if len(addrs.List) == 0 {
		store.Delete(key)
		return
	}
	store.Set(key, k.cdc.MustMarshal(&addrs))
}

// SetMaturityVscIDForChainIDConsAddr sets the vscID for the given chainID and consensus address.
// When the vscID matures on the chainID, the consensus address will be pruned.
func (k Keeper) SetMaturityVscIDForChainIDConsAddr(
	ctx sdk.Context, chainID string, consAddr sdk.ConsAddress, vscID uint64,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.MaturityVscIDForChainIDConsAddrKey(chainID, consAddr)
	store.Set(key, sdk.Uint64ToBigEndian(vscID))
}

// GetMaturityVscIDForChainIDConsAddr returns the vscID for the given chainID and consensus address.
// The vscID is used to prune the consensus address when the vscID matures on the chainID.
func (k Keeper) GetMaturityVscIDForChainIDConsAddr(
	ctx sdk.Context, chainID string, consAddr sdk.ConsAddress,
) uint64 {
	store := ctx.KVStore(k.storeKey)
	key := types.MaturityVscIDForChainIDConsAddrKey(chainID, consAddr)
	bz := store.Get(key)
	return sdk.BigEndianToUint64(bz)
}

// DeleteMaturityVscIDForChainIDConsAddr deletes the vscID for the given chainID and consensus address.
// The vscID is used to prune the consensus address when the vscID matures on the chainID.
func (k Keeper) DeleteMaturityVscIDForChainIDConsAddr(
	ctx sdk.Context, chainID string, consAddr sdk.ConsAddress,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.MaturityVscIDForChainIDConsAddrKey(chainID, consAddr)
	store.Delete(key)
}

// AppendUndelegationToRelease appends an undelegation record to the list of undelegations to release
// when the validator set update containing the given vscID is matured by the chainID.
func (k Keeper) AppendUndelegationToRelease(
	ctx sdk.Context, chainID string, vscID uint64, recordKey []byte,
) {
	prev := k.GetUndelegationsToRelease(ctx, chainID, vscID)
	prev.List = append(prev.List, recordKey)
	k.SetUndelegationsToRelease(ctx, chainID, vscID, prev)
}

// GetUndelegationsToRelease returns the list of undelegations to release when the validator set update
// containing the given vscID is matured by the chainID.
func (k Keeper) GetUndelegationsToRelease(
	ctx sdk.Context, chainID string, vscID uint64,
) (res types.UndelegationRecordKeys) {
	store := ctx.KVStore(k.storeKey)
	key := types.UndelegationsToReleaseKey(chainID, vscID)
	k.cdc.MustUnmarshal(store.Get(key), &res)
	return res
}

// SetUndelegationsToRelease sets the list of undelegations to release when the validator set update
// containing the given vscID is matured by the chainID.
func (k Keeper) SetUndelegationsToRelease(
	ctx sdk.Context, chainID string, vscID uint64, keys types.UndelegationRecordKeys,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.UndelegationsToReleaseKey(chainID, vscID)
	if len(keys.List) == 0 {
		store.Delete(key)
		return
	}
	store.Set(key, k.cdc.MustMarshal(&keys))
}
