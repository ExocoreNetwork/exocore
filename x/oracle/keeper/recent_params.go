package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRecentParams set a specific recentParams in the store from its index
func (k Keeper) SetRecentParams(ctx sdk.Context, recentParams types.RecentParams) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecentParamsKeyPrefix))
	b := k.cdc.MustMarshal(&recentParams)
	store.Set(types.RecentParamsKey(
		recentParams.Block,
	), b)
}

// GetRecentParams returns a recentParams from its index
func (k Keeper) GetRecentParams(
	ctx sdk.Context,
	block uint64,

) (val types.RecentParams, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecentParamsKeyPrefix))

	b := store.Get(types.RecentParamsKey(
		block,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveRecentParams removes a recentParams from the store
func (k Keeper) RemoveRecentParams(
	ctx sdk.Context,
	block uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecentParamsKeyPrefix))
	store.Delete(types.RecentParamsKey(
		block,
	))
}

// GetAllRecentParams returns all recentParams
func (k Keeper) GetAllRecentParams(ctx sdk.Context) (list []types.RecentParams) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecentParamsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RecentParams
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetAllRecentParamsAsMap(ctx sdk.Context) (result map[uint64]*types.Params) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RecentParamsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RecentParams
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		result[val.Block] = val.Params
	}

	return
}
