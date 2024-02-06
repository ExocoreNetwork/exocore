// This file is a duplicate of the subscriber module's validators file with minor changes.
// The function ApplyValidatorChanges can likely be carved out into a shared package.

package keeper

import (
	"time"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// UnbondingTime returns the time duration of the unbonding period. It is part of the
// implementation of the staking keeper expected by IBC.
// It is calculated as the number of epochs until unbonded multiplied by the duration of an
// epoch. This function is used by IBC's client keeper to validate the self client, and
// nowhere else. As long as it reports a consistent value, it's fine.
func (k Keeper) UnbondingTime(ctx sdk.Context) time.Duration {
	count := k.GetEpochsUntilUnbonded(ctx)
	identifier := k.GetEpochIdentifier(ctx)
	epoch, found := k.epochsKeeper.GetEpochInfo(ctx, identifier)
	if !found {
		panic("epoch info not found")
	}
	durationPerEpoch := epoch.Duration
	return time.Duration(count) * durationPerEpoch
}

// ApplyValidatorChanges returns the validator set as is. However, it also
// stores the validators that are added or those that are removed, and updates
// the power for the existing validators. It also allows any hooks registered
// on the keeper to be executed.
func (k Keeper) ApplyValidatorChanges(
	ctx sdk.Context,
	changes []abci.ValidatorUpdate,
) []abci.ValidatorUpdate {
	ret := []abci.ValidatorUpdate{}
	for _, change := range changes {
		// convert TM pubkey to SDK pubkey
		pubkey, err := cryptocodec.FromTmProtoPublicKey(change.GetPubKey())
		if err != nil {
			// An error here would indicate that the validator updates
			// received from other modules are invalid.
			panic(err)
		}
		addr := pubkey.Address()
		val, found := k.GetValidator(ctx, addr)

		if found {
			// update or delete an existing validator
			if change.Power < 1 {
				k.DeleteValidator(ctx, addr)
			} else {
				val.Power = change.Power
				k.SetValidator(ctx, val)
			}
		} else if change.Power > 0 {
			// create a new validator - the address is just derived from the public key and has
			// no correlation with the operator address on Exocore
			ocVal, err := types.NewExocoreValidator(addr, change.Power, pubkey)
			if err != nil {
				// An error here would indicate that the validator updates
				// received are invalid.
				panic(err)
			}

			k.SetValidator(ctx, ocVal)
			err = k.Hooks().AfterValidatorBonded(ctx, sdk.ConsAddress(addr), nil)
			if err != nil {
				// AfterValidatorBonded is hooked by the Slashing module and should not return
				// an error. If any other module were to hook it, they should also not.
				panic(err)
			}
		} else {
			// edge case: we received an update for 0 power
			// but the validator is already deleted. Do not forward
			// to tendermint.
			continue
		}

		ret = append(ret, change)
	}
	return ret
}

// SetValidator stores a validator based on the pub key derived address. This
// is accessible in the genesis state via `val_set`.
func (k Keeper) SetValidator(ctx sdk.Context, validator types.ExocoreValidator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validator)

	store.Set(types.ExocoreValidatorKey(validator.Address), bz)
}

// GetValidator gets a validator based on the pub key derived address.
func (k Keeper) GetValidator(
	ctx sdk.Context, addr []byte,
) (validator types.ExocoreValidator, found bool) {
	store := ctx.KVStore(k.storeKey)
	v := store.Get(types.ExocoreValidatorKey(addr))
	if v == nil {
		return
	}
	k.cdc.MustUnmarshal(v, &validator)
	found = true

	return
}

// DeleteValidator deletes a validator based on the pub key derived address.
func (k Keeper) DeleteValidator(ctx sdk.Context, addr []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ExocoreValidatorKey(addr))
}

// GetAllExocoreValidators returns all validators in the store.
func (k Keeper) GetAllExocoreValidators(
	ctx sdk.Context,
) (validators []types.ExocoreValidator) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ExocoreValidatorBytePrefix})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		val := types.ExocoreValidator{}
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		validators = append(validators, val)
	}

	return validators
}

// GetHistoricalInfo gets the historical info at a given height. It is part of the
// implementation of the staking keeper expected by IBC.
func (k Keeper) GetHistoricalInfo(
	ctx sdk.Context, height int64,
) (stakingtypes.HistoricalInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.HistoricalInfoKey(height)

	value := store.Get(key)
	if value == nil {
		return stakingtypes.HistoricalInfo{}, false
	}

	return stakingtypes.MustUnmarshalHistoricalInfo(k.cdc, value), true
}

// SetHistoricalInfo sets the historical info at a given height. This is
// (intentionally) not exported in the genesis state.
func (k Keeper) SetHistoricalInfo(
	ctx sdk.Context, height int64, hi *stakingtypes.HistoricalInfo,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.HistoricalInfoKey(height)
	value := k.cdc.MustMarshal(hi)

	store.Set(key, value)
}

// DeleteHistoricalInfo deletes the historical info at a given height.
func (k Keeper) DeleteHistoricalInfo(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	key := types.HistoricalInfoKey(height)

	store.Delete(key)
}

// TrackHistoricalInfo saves the latest historical-info and deletes the oldest
// heights that are below pruning height.
func (k Keeper) TrackHistoricalInfo(ctx sdk.Context) {
	numHistoricalEntries := k.GetHistoricalEntries(ctx)

	// Prune store to ensure we only have parameter-defined historical entries.
	// In most cases, this will involve removing a single historical entry.
	// In the rare scenario when the historical entries gets reduced to a lower value k'
	// from the original value k. k - k' entries must be deleted from the store.
	// Since the entries to be deleted are always in a continuous range, we can iterate
	// over the historical entries starting from the most recent version to be pruned
	// and then return at the first empty entry.
	for i := ctx.BlockHeight() - int64(numHistoricalEntries); i >= 0; i-- {
		_, found := k.GetHistoricalInfo(ctx, i)
		if found {
			k.DeleteHistoricalInfo(ctx, i)
		} else {
			break
		}
	}

	// if there is no need to persist historicalInfo, return.
	if numHistoricalEntries == 0 {
		return
	}

	// Create HistoricalInfo struct
	lastVals := []stakingtypes.Validator{}
	for _, v := range k.GetAllExocoreValidators(ctx) {
		pk, err := v.ConsPubKey()
		if err != nil {
			// This should never happen as the pubkey is assumed
			// to be stored correctly earlier.
			panic(err)
		}
		val, err := stakingtypes.NewValidator(nil, pk, stakingtypes.Description{})
		if err != nil {
			// This should never happen as the pubkey is assumed
			// to be stored correctly earlier.
			panic(err)
		}

		// Set validator to bonded status.
		val.Status = stakingtypes.Bonded
		// Compute tokens from voting power.
		val.Tokens = sdk.TokensFromConsensusPower(
			v.Power,
			// TODO(mm)
			// note that this is not super relevant for the historical info
			// since IBC does not seem to use the tokens field.
			sdk.NewInt(1),
		)
		lastVals = append(lastVals, val)
	}

	// Create historical info entry which sorts the validator set by voting power.
	historicalEntry := stakingtypes.NewHistoricalInfo(
		ctx.BlockHeader(), lastVals,
		// TODO(mm)
		// this should match the power reduction number above
		// and is also thus not relevant.
		sdk.NewInt(1),
	)

	// Set latest HistoricalInfo at current height.
	k.SetHistoricalInfo(ctx, ctx.BlockHeight(), &historicalEntry)
}

// MustGetCurrentValidatorsAsABCIUpdates gets all validators converted
// to the ABCI validator update type. It panics in case of failure.
func (k Keeper) MustGetCurrentValidatorsAsABCIUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	vals := k.GetAllExocoreValidators(ctx)
	valUpdates := make([]abci.ValidatorUpdate, 0, len(vals))
	for _, v := range vals {
		pk, err := v.ConsPubKey()
		if err != nil {
			// This should never happen as the pubkey is assumed
			// to be stored correctly earlier.
			panic(err)
		}
		tmPK, err := cryptocodec.ToTmProtoPublicKey(pk)
		if err != nil {
			// This should never happen as the pubkey is assumed
			// to be stored correctly earlier.
			panic(err)
		}
		valUpdates = append(valUpdates, abci.ValidatorUpdate{PubKey: tmPK, Power: v.Power})
	}
	return valUpdates
}
