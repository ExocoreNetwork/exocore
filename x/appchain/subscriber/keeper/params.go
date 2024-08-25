package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetParams sets the appchain coordinator parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey(), bz)
}

// GetParams gets the appchain coordinator parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey())
	if bz == nil {
		return types.DefaultParams()
	}
	var params types.Params
	k.cdc.MustUnmarshal(bz, &params)
	return params
}
