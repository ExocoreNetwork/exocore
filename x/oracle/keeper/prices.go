package keeper

import (
	"encoding/binary"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetPrices set a specific prices in the store from its index
func (k Keeper) SetPrices(ctx sdk.Context, prices types.Prices) {
	store := k.getPriceTRStore(ctx, prices.TokenID)
	for _, v := range prices.PriceList {
		b := k.cdc.MustMarshal(v)
		store.Set(types.PricesRoundKey(v.RoundID), b)
	}
	store.Set(types.PricesNextRoundIDKey, types.Uint64Bytes(prices.NextRoundID))
}

// GetPrices returns a prices from its index
func (k Keeper) GetPrices(
	ctx sdk.Context,
	tokenID uint64,
) (val types.Prices, found bool) {
	store := k.getPriceTRStore(ctx, tokenID)
	nextRoundID := k.GetNextRoundID(ctx, tokenID)

	val.TokenID = tokenID
	val.NextRoundID = nextRoundID
	val.PriceList = make([]*types.PriceTimeRound, nextRoundID)
	// 0 roundId is reserved, expect the roundid corresponds to the slice index
	val.PriceList[0] = &types.PriceTimeRound{}
	for i := uint64(1); i < nextRoundID; i++ {
		b := store.Get(types.PricesRoundKey(i))
		val.PriceList[i] = &types.PriceTimeRound{}
		if b != nil {
			// should always be true since we don't delete prices from history round
			k.cdc.MustUnmarshal(b, val.PriceList[i])
			found = true
		}
	}

	return
}

// RemovePrices removes a prices from the store
func (k Keeper) RemovePrices(
	ctx sdk.Context,
	tokenID uint64,
) {
	store := k.getPriceTRStore(ctx, tokenID)
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

	var price types.Prices
	prevTokenID := uint64(0)
	for ; iterator.Valid(); iterator.Next() {
		tokenID, _, nextRoundID := parseKey(iterator.Key())
		if prevTokenID == 0 {
			prevTokenID = tokenID
			price.TokenID = tokenID
		} else if prevTokenID != tokenID && price.TokenID > 0 {
			list = append(list, price)
			prevTokenID = tokenID
			price = types.Prices{TokenID: tokenID}
		}
		if nextRoundID {
			price.NextRoundID = binary.BigEndian.Uint64(iterator.Value())
		} else {
			var val types.PriceTimeRound
			k.cdc.MustUnmarshal(iterator.Value(), &val)
			price.PriceList = append(price.PriceList, &val)
		}
	}
	if price.TokenID > 0 {
		list = append(list, price)
	}
	return list
}

func (k Keeper) AppendPriceTR(ctx sdk.Context, tokenID uint64, priceTR types.PriceTimeRound) {
	nextRoundID := k.GetNextRoundID(ctx, tokenID)
	if nextRoundID != priceTR.RoundID {
		// TODO: return error to tell this round adding fail
		return
	}
	store := k.getPriceTRStore(ctx, tokenID)
	b := k.cdc.MustMarshal(&priceTR)
	store.Set(types.PricesRoundKey(nextRoundID), b)
	k.IncreaseNextRoundID(ctx, tokenID)
}

// func(k Keeper) SetPriceTR(ctx sdk.Context, tokenID int32, priceTR){}
func (k Keeper) GetPriceTRRoundID(ctx sdk.Context, tokenID uint64, roundID uint64) (price types.PriceTimeRound, found bool) {
	store := k.getPriceTRStore(ctx, tokenID)

	b := store.Get(types.PricesRoundKey(roundID))
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &price)
	found = true
	return
}

func (k Keeper) GetPriceTRLatest(ctx sdk.Context, tokenID uint64) (price types.PriceTimeRound, found bool) {
	//	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	//	store = prefix.NewStore(store, types.PricesKey(tokenID))
	store := k.getPriceTRStore(ctx, tokenID)
	nextRoundIDB := store.Get(types.PricesNextRoundIDKey)
	if nextRoundIDB == nil {
		return
	}
	nextRoundID := binary.BigEndian.Uint64(nextRoundIDB)
	b := store.Get(types.PricesRoundKey(nextRoundID - 1))
	if b != nil {
		// should always be true
		k.cdc.MustUnmarshal(b, &price)
		found = true
	}
	return
}

func (k Keeper) GetNextRoundID(ctx sdk.Context, tokenID uint64) (nextRoundID uint64) {
	nextRoundID = 1
	store := k.getPriceTRStore(ctx, tokenID)
	nextRoundIDB := store.Get(types.PricesNextRoundIDKey)
	if nextRoundIDB != nil {
		if nextRoundID = binary.BigEndian.Uint64(nextRoundIDB); nextRoundID == 0 {
			nextRoundID = 1
		}
	}
	return
}

func (k Keeper) IncreaseNextRoundID(ctx sdk.Context, tokenID uint64) {
	store := k.getPriceTRStore(ctx, tokenID)
	nextRoundID := k.GetNextRoundID(ctx, tokenID)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, nextRoundID+1)
	store.Set(types.PricesNextRoundIDKey, b)
}

func (k Keeper) getPriceTRStore(ctx sdk.Context, tokenID uint64) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PricesKeyPrefix))
	return prefix.NewStore(store, types.PricesKey(tokenID))
}

func parseKey(key []byte) (tokenID uint64, roundID uint64, nextRoundID bool) {
	tokenID = binary.BigEndian.Uint64(key[:8])
	if len(key) == 21 {
		nextRoundID = true
		return
	}
	roundID = binary.BigEndian.Uint64(key[9:17])
	return
}
