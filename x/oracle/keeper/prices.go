package keeper

import (
	"encoding/binary"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetPrices set a specific prices in the store from its index
func (k Keeper) SetPrices(ctx sdk.Context, prices types.Prices) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	store = prefix.NewStore(store, types.PricesKey(prices.TokenId))
	for _, v := range prices.PriceList {
		b := k.cdc.MustMarshal(v)
		store.Set(types.PricesRoundKey(v.RoundId), b)
	}
	store.Set(types.PricesNextRountIdKey, types.Uint64Bytes(uint64(prices.NextRountId)))
}

// GetPrices returns a prices from its index
func (k Keeper) GetPrices(
	ctx sdk.Context,
	tokenId int32,

) (val types.Prices, found bool) {
	//	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	//	store = prefix.NewStore(store, types.PricesKey(tokenId))
	store := k.getPriceTRStore(ctx, tokenId)
	nextRoundIdB := store.Get(types.PricesNextRountIdKey)
	if nextRoundIdB == nil {
		return val, false
	}

	nextRoundId := binary.BigEndian.Uint64(nextRoundIdB)

	val.TokenId = tokenId
	val.NextRountId = nextRoundId
	val.PriceList = make([]*types.PriceWithTimeAndRound, nextRoundId)
	//0 roundId is reserved
	val.PriceList[0] = &types.PriceWithTimeAndRound{}
	for i := uint64(1); i < nextRoundId; i++ {
		b := store.Get(types.PricesRoundKey(i))
		val.PriceList[i] = &types.PriceWithTimeAndRound{}
		if b != nil {
			//should alwyas be true since we don't delete prices from history round
			k.cdc.MustUnmarshal(b, val.PriceList[i])
		}
	}

	return val, true
}

// RemovePrices removes a prices from the store
func (k Keeper) RemovePrices(
	ctx sdk.Context,
	tokenId int32,

) {
	//	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	//	store = prefix.NewStore(store, types.PricesKey(tokenId))
	store := k.getPriceTRStore(ctx, tokenId)
	//	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
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

func (k Keeper) AppendPriceTR(ctx sdk.Context, tokenId int32, priceTR types.PriceWithTimeAndRound) {
	nextRoundId := k.GetNextRoundId(ctx, tokenId)
	if nextRoundId != priceTR.RoundId {
		return
	}
	store := k.getPriceTRStore(ctx, tokenId)
	b := k.cdc.MustMarshal(&priceTR)
	store.Set(types.PricesRoundKey(nextRoundId), b)
}

//func(k Keeper) SetPriceTR(ctx sdk.Context, tokenId int32, priceTR){}

func (k Keeper) GetPriceTRRoundId(ctx sdk.Context, tokenId int32, roundId uint64) (price types.PriceWithTimeAndRound, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	store = prefix.NewStore(store, types.PricesKey(tokenId))

	b := store.Get(types.PricesRoundKey(roundId))
	if b == nil {
		return
	}

	k.cdc.Unmarshal(b, &price)
	found = true
	return
}

func (k Keeper) GetPriceTRLatest(ctx sdk.Context, tokenId int32) (price types.PriceWithTimeAndRound, found bool) {
	//	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	//	store = prefix.NewStore(store, types.PricesKey(tokenId))
	store := k.getPriceTRStore(ctx, tokenId)
	nextRoundIdB := store.Get(types.PricesNextRountIdKey)
	if nextRoundIdB == nil {
		return
	}
	nextRoundId := binary.BigEndian.Uint64(nextRoundIdB)
	b := store.Get(types.PricesRoundKey(nextRoundId - 1))
	if b != nil {
		//should always be true
		k.cdc.Unmarshal(b, &price)
		found = true
	}
	return
}

func (k Keeper) GetNextRoundId(ctx sdk.Context, tokenId int32) (nextRoundId uint64) {
	nextRoundId = 1
	//store := getPriceTRStore(ctx, k.storeKey, tokenId)
	store := k.getPriceTRStore(ctx, tokenId)
	nextRoundIdB := store.Get(types.PricesNextRountIdKey)
	if nextRoundIdB != nil {
		nextRoundId = binary.BigEndian.Uint64(nextRoundIdB)
	}
	return
}

//func getPriceTRStore(ctx sdk.Context, storeKey storetypes.StoreKey, tokenId int32) prefix.Store {
//	store := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.PricesKeyPrefix))
//	return prefix.NewStore(store, types.PricesKey(tokenId))
//}

func (k Keeper) getPriceTRStore(ctx sdk.Context, tokenId int32) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	return prefix.NewStore(store, types.PricesKey(tokenId))
}
