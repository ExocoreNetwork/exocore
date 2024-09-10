package keeper

import (
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetSubscriberValidatorForChain sets the subscriber validator for a chain.
// Storing this historical information allows us to minimize the number/size of
// validator set updates sent to the subscriber by skipping the keys for which
// there is no change in vote power.
func (k Keeper) SetSubscriberValidatorForChain(
	ctx sdk.Context, chainID string, validator commontypes.SubscriberValidator,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.SubscriberValidatorKey(chainID, validator.ConsAddress)
	bz := k.cdc.MustMarshal(&validator)
	store.Set(key, bz)
}

// GetSubscriberValidatorForChain gets the subscriber validator for a chain.
func (k Keeper) GetSubscriberValidatorForChain(
	ctx sdk.Context, chainID string, consAddress []byte,
) (validator commontypes.SubscriberValidator, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.SubscriberValidatorKey(chainID, consAddress)
	if !store.Has(key) {
		return validator, false
	}
	bz := store.Get(key)
	k.cdc.MustUnmarshal(bz, &validator)
	return validator, true
}

// GetAllSubscriberValidatorsForChain gets all subscriber validators for a chain, ordered
// by the consensus address bytes.
func (k Keeper) GetAllSubscriberValidatorsForChain(
	ctx sdk.Context, chainID string,
) (validators []commontypes.SubscriberValidator) {
	store := ctx.KVStore(k.storeKey)
	partialKey := types.SubscriberValidatorKey(chainID, nil)
	iterator := sdk.KVStorePrefixIterator(store, partialKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var validator commontypes.SubscriberValidator
		k.cdc.MustUnmarshal(iterator.Value(), &validator)
		validators = append(validators, validator)
	}

	return validators
}

// DeleteSubscriberValidatorForChain deletes the subscriber validator for a chain, given
// the consensus address.
func (k Keeper) DeleteSubscriberValidatorForChain(
	ctx sdk.Context, chainID string, consAddress []byte,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.SubscriberValidatorKey(chainID, consAddress)
	store.Delete(key)
}

// SetMaxValidatorsForChain sets the maximum number of validators for a chain.
func (k Keeper) SetMaxValidatorsForChain(
	ctx sdk.Context, chainID string, maxValidators uint32,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.MaxValidatorsKey(chainID)
	store.Set(key, sdk.Uint64ToBigEndian(uint64(maxValidators)))
}

// GetMaxValidatorsForChain gets the maximum number of validators for a chain.
func (k Keeper) GetMaxValidatorsForChain(
	ctx sdk.Context, chainID string,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	key := types.MaxValidatorsKey(chainID)
	bz := store.Get(key)
	// #nosec G115 // we stored it, we trust it
	return uint32(sdk.BigEndianToUint64(bz))
}
