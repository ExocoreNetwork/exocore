package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetPrices set a specific prices in the store from its index
func (k Keeper) SetPrices(ctx sdk.Context, prices types.Prices) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	b := k.cdc.MustMarshal(&prices)
	store.Set(types.PricesKey(
		prices.TokenId,
	), b)
}

// GetPrices returns a prices from its index
func (k Keeper) GetPrices(
	ctx sdk.Context,
	tokenId int32,

) (val types.Prices, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))

	b := store.Get(types.PricesKey(
		tokenId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePrices removes a prices from the store
func (k Keeper) RemovePrices(
	ctx sdk.Context,
	tokenId int32,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	store.Delete(types.PricesKey(
		tokenId,
	))
}

// GetAllPrices returns all prices
func (k Keeper) GetAllPrices(ctx sdk.Context) (list []types.Prices) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Prices
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
