package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
//func (k Keeper) GetParams(ctx sdk.Context) types.Params {
//	pRes := &types.Params{}
//	k.paramstore.GetParamSet(ctx, pRes)
//	return *pRes
//	// return types.NewParams()
//}
//
//// SetParams set the params
//func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
//	//TODO: update the aggregator's params, call k.UpdateParams
//	k.paramstore.SetParamSet(ctx, &params)
//}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey) // return types.NewParams()
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &params)
	}
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	// TODO: validation check
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
	return nil
}
