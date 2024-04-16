package keeper

import (
	"encoding/binary"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetPrices set a specific prices in the store from its index
func (k Keeper) SetPrices(ctx sdk.Context, prices types.Prices) {
	store := k.getPriceTRStore(ctx, prices.TokenId)
	for _, v := range prices.PriceList {
		b := k.cdc.MustMarshal(v)
		store.Set(types.PricesRoundKey(v.RoundId), b)
	}
	store.Set(types.PricesNextRoundIdKey, types.Uint64Bytes(prices.NextRoundId))
}

// GetPrices returns a prices from its index
func (k Keeper) GetPrices(
	ctx sdk.Context,
	tokenId uint64,
) (val types.Prices, found bool) {
	store := k.getPriceTRStore(ctx, tokenId)
	nextRoundId := k.GetNextRoundId(ctx, tokenId)

	val.TokenId = tokenId
	val.NextRoundId = nextRoundId
	val.PriceList = make([]*types.PriceWithTimeAndRound, nextRoundId)
	// 0 roundId is reserved, expect the roundid corresponds to the slice index
	val.PriceList[0] = &types.PriceWithTimeAndRound{}
	for i := uint64(1); i < nextRoundId; i++ {
		b := store.Get(types.PricesRoundKey(i))
		val.PriceList[i] = &types.PriceWithTimeAndRound{}
		if b != nil {
			// should alwyas be true since we don't delete prices from history round
			k.cdc.MustUnmarshal(b, val.PriceList[i])
			found = true
		}
	}

	return
}

// RemovePrices removes a prices from the store
func (k Keeper) RemovePrices(
	ctx sdk.Context,
	tokenId uint64,
) {
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

	//	prevTokenId := uint32(0)
	//	var val types.PriceWithTimeAndRound
	var price types.Prices
	prevTokenId := uint64(0)
	for ; iterator.Valid(); iterator.Next() {
		tokenId, _, nextRoundId := parseKey(iterator.Key())
		if prevTokenId == 0 {
			prevTokenId = tokenId
			price.TokenId = tokenId
		} else if prevTokenId != tokenId && price.TokenId > 0 {
			list = append(list, price)
			prevTokenId = tokenId
			price = types.Prices{TokenId: tokenId}
		}
		if nextRoundId {
			price.NextRoundId = binary.BigEndian.Uint64(iterator.Value())
		} else {
			var val types.PriceWithTimeAndRound
			k.cdc.MustUnmarshal(iterator.Value(), &val)
			price.PriceList = append(price.PriceList, &val)
		}
	}
	if price.TokenId > 0 {
		list = append(list, price)
	}
	return list
}

func (k Keeper) AppendPriceTR(ctx sdk.Context, tokenId uint64, priceTR types.PriceWithTimeAndRound) {
	nextRoundId := k.GetNextRoundId(ctx, tokenId)
	if nextRoundId != priceTR.RoundId {
		return
	}
	store := k.getPriceTRStore(ctx, tokenId)
	b := k.cdc.MustMarshal(&priceTR)
	store.Set(types.PricesRoundKey(nextRoundId), b)
	k.IncreaseNextRoundId(ctx, tokenId)
}

// func(k Keeper) SetPriceTR(ctx sdk.Context, tokenId int32, priceTR){}
func (k Keeper) GetPriceTRRoundId(ctx sdk.Context, tokenId uint64, roundId uint64) (price types.PriceWithTimeAndRound, found bool) {
	store := k.getPriceTRStore(ctx, tokenId)

	b := store.Get(types.PricesRoundKey(roundId))
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &price)
	found = true
	return
}

func (k Keeper) GetPriceTRLatest(ctx sdk.Context, tokenId uint64) (price types.PriceWithTimeAndRound, found bool) {
	//	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	//	store = prefix.NewStore(store, types.PricesKey(tokenId))
	store := k.getPriceTRStore(ctx, tokenId)
	nextRoundIdB := store.Get(types.PricesNextRoundIdKey)
	if nextRoundIdB == nil {
		return
	}
	nextRoundId := binary.BigEndian.Uint64(nextRoundIdB)
	b := store.Get(types.PricesRoundKey(nextRoundId - 1))
	if b != nil {
		// should always be true
		k.cdc.MustUnmarshal(b, &price)
		found = true
	}
	return
}

func (k Keeper) GetNextRoundId(ctx sdk.Context, tokenId uint64) (nextRoundId uint64) {
	nextRoundId = 1
	// store := getPriceTRStore(ctx, k.storeKey, tokenId)
	store := k.getPriceTRStore(ctx, tokenId)
	nextRoundIdB := store.Get(types.PricesNextRoundIdKey)
	if nextRoundIdB != nil {
		if nextRoundId = binary.BigEndian.Uint64(nextRoundIdB); nextRoundId == 0 {
			nextRoundId = 1
		}
	}
	return
}

func (k Keeper) IncreaseNextRoundId(ctx sdk.Context, tokenId uint64) {
	// store := getPriceTRStore(ctx, k.storeKey, tokenId)
	store := k.getPriceTRStore(ctx, tokenId)
	nextRoundId := k.GetNextRoundId(ctx, tokenId)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextRoundId+1)
	store.Set(types.PricesNextRoundIdKey, b)
}

//func getPriceTRStore(ctx sdk.Context, storeKey storetypes.StoreKey, tokenId int32) prefix.Store {
//	store := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.PricesKeyPrefix))
//	return prefix.NewStore(store, types.PricesKey(tokenId))
//}

func (k Keeper) getPriceTRStore(ctx sdk.Context, tokenId uint64) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	return prefix.NewStore(store, types.PricesKey(tokenId))
}

func parseKey(key []byte) (tokenId uint64, roundId uint64, nextRoundId bool) {
	tokenId = binary.BigEndian.Uint64(key[:8])
	if len(key) == 21 {
		nextRoundId = true
		return
	}
	roundId = binary.BigEndian.Uint64(key[9:17])
	return
}
