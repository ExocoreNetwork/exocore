package keeper

import (
	"fmt"

	exocoretypes "github.com/ExocoreNetwork/exocore/types/keys"
	types "github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetValsetUpdateIDForHeight sets the valset update ID for a given height
func (k Keeper) SetValsetUpdateIDForHeight(
	ctx sdk.Context, height int64, valsetUpdateID uint64,
) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ValsetUpdateIDKey(height), sdk.Uint64ToBigEndian(valsetUpdateID))
}

// GetValsetUpdateIDForHeight gets the valset update ID for a given height
func (k Keeper) GetValsetUpdateIDForHeight(
	ctx sdk.Context, height int64,
) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.ValsetUpdateIDKey(height)
	if !store.Has(key) {
		return 0, false
	}
	bz := store.Get(key)
	return sdk.BigEndianToUint64(bz), true
}

// ApplyValidatorChanges is a wrapper function that returns the provided validator set
// update. The wrapping allows to save the validator set information in the store.
// The caller should (but _not_ must) provide `changes` that are different from the
// ones already with Tendermint.
func (k Keeper) ApplyValidatorChanges(
	ctx sdk.Context,
	// in dogfood, we use the wrappedkeywithpower because the operator module provides
	// keys in that format. since the subscriber chain doesn't need the operator module
	// we can use the tm validator update type.
	changes []abci.ValidatorUpdate,
) []abci.ValidatorUpdate {
	ret := make([]abci.ValidatorUpdate, 0, len(changes))
	logger := k.Logger(ctx)
	for i := range changes {
		change := changes[i] // avoid implicit memory aliasing
		wrappedKey := exocoretypes.NewWrappedConsKeyFromTmProtoKey(&change.PubKey)
		if wrappedKey == nil {
			// an error in deserializing the key would indicate that the coordinator
			// has provided invalid data. this is a critical error and should be
			// investigated.
			panic(fmt.Sprintf("invalid pubkey %s", change.PubKey))
		}
		consAddress := wrappedKey.ToConsAddr()
		val, found := k.GetSubscriberChainValidator(ctx, consAddress)
		switch found {
		case true:
			if change.Power < 1 {
				logger.Info("deleting validator", "consAddress", consAddress)
				k.DeleteSubscriberChainValidator(ctx, consAddress)
			} else {
				logger.Info("updating validator", "consAddress", consAddress)
				val.Power = change.Power
				k.SetSubscriberChainValidator(ctx, val)
			}
		case false:
			if change.Power > 0 {
				ocVal, err := types.NewSubscriberChainValidator(
					consAddress, change.Power, wrappedKey.ToSdkKey(),
				)
				if err != nil {
					// cannot happen, but just in case add this check.
					// simply skip the validator if it does.
					continue
				}
				logger.Info("adding validator", "consAddress", consAddress)
				k.SetSubscriberChainValidator(ctx, ocVal)
				ret = append(ret, change)
			} else {
				// edge case: we received an update for 0 power
				// but the validator is already deleted. Do not forward
				// to tendermint.
				logger.Info(
					"received update for non-existent validator",
					"cons address", consAddress,
				)
				continue
			}
		}
		ret = append(ret, change)
	}
	return ret
}

// SetSubscriberChainValidator stores a validator based on the pub key derived address.
func (k Keeper) SetSubscriberChainValidator(
	ctx sdk.Context, validator types.SubscriberChainValidator,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validator)

	store.Set(types.SubscriberChainValidatorKey(validator.ConsAddress), bz)
}

// GetSubscriberChainValidator gets a validator based on the pub key derived (consensus) address.
func (k Keeper) GetSubscriberChainValidator(
	ctx sdk.Context, addr sdk.ConsAddress,
) (validator types.SubscriberChainValidator, found bool) {
	store := ctx.KVStore(k.storeKey)
	v := store.Get(types.SubscriberChainValidatorKey(addr))
	if v == nil {
		return
	}
	k.cdc.MustUnmarshal(v, &validator)
	found = true

	return
}

// DeleteSubscriberChainValidator deletes a validator based on the pub key derived address.
func (k Keeper) DeleteSubscriberChainValidator(ctx sdk.Context, addr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.SubscriberChainValidatorKey(addr))
}

// GetAllSubscriberChainValidators returns all validators in the store.
func (k Keeper) GetAllSubscriberChainValidators(
	ctx sdk.Context,
) (validators []types.SubscriberChainValidator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.SubscriberChainValidatorBytePrefix})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		val := types.SubscriberChainValidator{}
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		validators = append(validators, val)
	}

	return validators
}
