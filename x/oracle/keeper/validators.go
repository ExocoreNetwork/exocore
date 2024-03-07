package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetValidators set validators in the store
func (k Keeper) SetValidators(ctx sdk.Context, validators types.Validators) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKey))
	b := k.cdc.MustMarshal(&validators)
	store.Set([]byte{0}, b)
}

// GetValidators returns validators
func (k Keeper) GetValidators(ctx sdk.Context) (val types.Validators, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveValidators removes validators from the store
func (k Keeper) RemoveValidators(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKey))
	store.Delete([]byte{0})
}
