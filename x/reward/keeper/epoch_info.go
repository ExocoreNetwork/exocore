package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetEpochIdentifier gets the epoch identifier
func (k Keeper) GetEpochIdentifier(_ sdk.Context) string {
	// TODO: compatible with evmos v16
	// store := ctx.KVStore(k.storeKey)
	// bz := store.Get(types.KeyPrefixEpochIdentifier)
	// if len(bz) == 0 {
	// 	return ""
	// }
	//
	// return string(bz)
	return ""
}

// SetEpochsPerPeriod stores the epoch identifier
func (k Keeper) SetEpochIdentifier(_ sdk.Context, _ string) {
	// TODO: compatible with evmos v16
	// store := ctx.KVStore(k.storeKey)
	// store.Set(types.KeyPrefixEpochIdentifier, []byte(epochIdentifier))
}
