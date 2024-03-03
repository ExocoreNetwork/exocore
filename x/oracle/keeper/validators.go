package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetValidators set a specific validators in the store from its index
func (k Keeper) SetValidators(ctx sdk.Context, validators types.Validators) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKeyPrefix))
	b := k.cdc.MustMarshal(&validators)
	store.Set(types.ValidatorsKey(
		validators.Block,
	), b)
}

// GetValidators returns a validators from its index
func (k Keeper) GetValidators(
	ctx sdk.Context,
	block uint64,

) (val types.Validators, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKeyPrefix))

	b := store.Get(types.ValidatorsKey(
		block,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveValidators removes a validators from the store
func (k Keeper) RemoveValidators(
	ctx sdk.Context,
	block uint64,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKeyPrefix))
	store.Delete(types.ValidatorsKey(
		block,
	))
}

// GetAllValidators returns all validators
func (k Keeper) GetAllValidators(ctx sdk.Context) (list []types.Validators) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ValidatorsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Validators
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
