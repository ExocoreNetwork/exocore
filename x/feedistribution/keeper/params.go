package keeper

import (
	"cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefixParams
	bz := store.Get(key)
	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyPrefixParams
	bz := k.cdc.MustMarshal(&params)
	store.Set(key, bz)
}

// GetCommunityTax returns the current distribution community tax.
func (k Keeper) GetCommunityTax(ctx sdk.Context) (math.LegacyDec, error) {
	params := k.GetParams(ctx)

	return params.CommunityTax, nil
}
