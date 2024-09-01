package keeper

import (
	types "github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetCoordinatorClientID sets the clientID of the coordinator chain
func (k Keeper) SetCoordinatorClientID(ctx sdk.Context, clientID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CoordinatorClientIDKey(), []byte(clientID))
}

// GetCoordinatorClientID gets the clientID of the coordinator chain
func (k Keeper) GetCoordinatorClientID(ctx sdk.Context) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.CoordinatorClientIDKey()
	if !store.Has(key) {
		return "", false
	}
	bz := store.Get(key)
	return string(bz), true
}
