package keeper

import (
	paramstypes "github.com/ExocoreNetwork/exocore/x/deposit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (*paramstypes.Params, error) {
	// store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixParams)
	// ifExist := store.Has(types.ParamsKey)
	// if !ifExist {
	// 	return nil, types.ErrNoParamsKey
	// }

	// value := store.Get(types.ParamsKey)

	// ret := &types.Params{}
	// k.cdc.MustUnmarshal(value, ret)
	// return ret, nil
	// Uify the way to obtain Params from deposit keeper
	return k.depositKeeper.GetParams(ctx)
}
