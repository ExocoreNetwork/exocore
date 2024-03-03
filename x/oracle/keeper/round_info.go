package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRoundInfo set a specific roundInfo in the store from its index
func (k Keeper) SetRoundInfo(ctx sdk.Context, roundInfo types.RoundInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoundInfoKeyPrefix))
	b := k.cdc.MustMarshal(&roundInfo)
	store.Set(types.RoundInfoKey(
		roundInfo.TokenId,
	), b)
}

// GetRoundInfo returns a roundInfo from its index
func (k Keeper) GetRoundInfo(
	ctx sdk.Context,
	tokenId int32,

) (val types.RoundInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoundInfoKeyPrefix))

	b := store.Get(types.RoundInfoKey(
		tokenId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveRoundInfo removes a roundInfo from the store
func (k Keeper) RemoveRoundInfo(
	ctx sdk.Context,
	tokenId int32,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoundInfoKeyPrefix))
	store.Delete(types.RoundInfoKey(
		tokenId,
	))
}

// GetAllRoundInfo returns all roundInfo
func (k Keeper) GetAllRoundInfo(ctx sdk.Context) (list []types.RoundInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoundInfoKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RoundInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
