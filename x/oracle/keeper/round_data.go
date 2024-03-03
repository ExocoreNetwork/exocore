package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetRoundData set a specific roundData in the store from its index
func (k Keeper) SetRoundData(ctx sdk.Context, roundData types.RoundData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoundDataKeyPrefix))
	b := k.cdc.MustMarshal(&roundData)
	store.Set(types.RoundDataKey(
		roundData.TokenId,
	), b)
}

// GetRoundData returns a roundData from its index
func (k Keeper) GetRoundData(
	ctx sdk.Context,
	tokenId int32,

) (val types.RoundData, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoundDataKeyPrefix))

	b := store.Get(types.RoundDataKey(
		tokenId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveRoundData removes a roundData from the store
func (k Keeper) RemoveRoundData(
	ctx sdk.Context,
	tokenId int32,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoundDataKeyPrefix))
	store.Delete(types.RoundDataKey(
		tokenId,
	))
}

// GetAllRoundData returns all roundData
func (k Keeper) GetAllRoundData(ctx sdk.Context) (list []types.RoundData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RoundDataKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RoundData
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
