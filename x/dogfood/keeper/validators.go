// This file is a duplicate of the subscriber module's validators file with minor changes.
// The function ApplyValidatorChanges can likely be carved out into a shared package with
// the appchain module.

package keeper

import (
	"sort"
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
	// no need to check for found, as the epoch info is validated at genesis.
	epoch, _ := k.epochsKeeper.GetEpochInfo(ctx, identifier)
	durationPerEpoch := epoch.Duration
	return time.Duration(count) * durationPerEpoch
}

// ApplyValidatorChanges returns the validator set as is. However, it also
// stores the validators that are added or those that are removed, and updates
// the stored power for the existing validators. It also allows any hooks registered
// on the keeper to be executed. Lastly, it stores the validator set against the
// provided validator set id.
func (k Keeper) ApplyValidatorChanges(
	ctx sdk.Context, changes []abci.ValidatorUpdate,
) []abci.ValidatorUpdate {
	ret := []abci.ValidatorUpdate{}
	for _, change := range changes {
		// convert TM pubkey to SDK pubkey for storage within the validator object.
		pubkey, err := cryptocodec.FromTmProtoPublicKey(change.GetPubKey())
		if err != nil {
			// An error here would indicate that this change is invalid.
			// The change is received either from the genesis file, or from
			// other parts of the module.
			// In no situation it should happen; however, if it does,
			// we do not panic. Simply skip the change.
			continue
		}
		// the address is just derived from the public key and
		// has no correlation with the operator address on Exocore.
		addr := pubkey.Address()
		val, found := k.GetValidator(ctx, addr)
		switch found {
		case true:
			// update or delete an existing validator.
			// assumption: power can not be negative.
			if change.Power < 1 {
				k.DeleteValidator(ctx, addr)
			} else {
				val.Power = change.Power
				k.SetValidator(ctx, val)
			}
		case false:
			if change.Power > 0 {
				// create a new validator.
				ocVal, err := types.NewExocoreValidator(addr, change.Power, pubkey)
				if err != nil {
					continue
				}
				// guard for errors within the AfterValidatorBonded hook.
				cc, writeFunc := ctx.CacheContext()
				k.SetValidator(cc, ocVal)
				err = k.Hooks().AfterValidatorBonded(cc, sdk.ConsAddress(addr), nil)
				if err != nil {
					// If an error is returned, the validator is not added to the `ret` slice.
					continue
				}
				writeFunc()
			} else {
				// edge case: we received an update for 0 power
				// but the validator is already deleted. Do not forward
				// to tendermint.
				continue
			}
		}
		ret = append(ret, change)
	}

	// sort for determinism
	sort.Slice(ret, func(i, j int) bool {
		if ret[i].Power != ret[j].Power {
			return ret[i].Power > ret[j].Power
		}
		return ret[i].PubKey.String() > ret[j].PubKey.String()
	})

	return ret
}

// SetValidator stores a validator based on the pub key derived address. This
// is accessible in the genesis state via `val_set`.
func (k Keeper) SetValidator(ctx sdk.Context, validator types.ExocoreValidator) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validator)

	store.Set(types.ExocoreValidatorKey(validator.Address), bz)
}

// GetValidator gets a validator based on the pub key derived (consensus) address.
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

// GetHistoricalInfo gets the historical info at a given height
func (k Keeper) GetHistoricalInfo(
	ctx sdk.Context, height int64,
) (stakingtypes.HistoricalInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.HistoricalInfoKey(height)

	value := store.Get(key)
	if value == nil {
		return stakingtypes.HistoricalInfo{}, false
	}

	return stakingtypes.MustUnmarshalHistoricalInfo(k.cdc, value), true
}

// SetHistoricalInfo sets the historical info at a given height
func (k Keeper) SetHistoricalInfo(
	ctx sdk.Context, height int64, hi *stakingtypes.HistoricalInfo,
) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.HistoricalInfoKey(height)
	value := k.cdc.MustMarshal(hi)

	store.Set(key, value)
}

// DeleteHistoricalInfo deletes the historical info at a given height
func (k Keeper) DeleteHistoricalInfo(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	key, _ := types.HistoricalInfoKey(height)

	store.Delete(key)
}

// TrackHistoricalInfo saves the latest historical info and deletes the ones eligible to be
// pruned. The function is called within the EndBlock of the module, so it is kept public.
// It is mostly a copy of the function used by interchain-security.
// If the historical info were only used by IBC, this function would store a subset of the
// header for each block, since only those parts were used.
// However, the historical info is used by the EVM keeper as well, which hashes the full header
// to report it via Solidity to the caller. Therefore, the full header must be stored.
func (k Keeper) TrackHistoricalInfo(ctx sdk.Context) {
	// Get the number of historical entries to persist, as the number of block heights.
	// #nosec G701 // uint32 fits into int64 always.
	numHistoricalEntries := int64(
		k.GetHistoricalEntries(ctx),
	)

	// we are deleting headers, say, from, 0 to 999 at block 1999
	for i := ctx.BlockHeight() - numHistoricalEntries; i >= 0; i-- {
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
			// since we stored the validator in the first place, something like this
			// should never happen, but if it does it is an extremely grave error
			// that will result in a block mismatch and hence that node will halt.
			continue
		}
		val, err := stakingtypes.NewValidator(nil, pk, stakingtypes.Description{})
		if err != nil {
			// same as above.
			continue
		}

		// Set validator to bonded status
		val.Status = stakingtypes.Bonded
		// Compute tokens from voting power
		val.Tokens = sdk.TokensFromConsensusPower(v.Power, sdk.DefaultPowerReduction)
		lastVals = append(lastVals, val)
	}

	// Create historical info entry which sorts the validator set by voting power
	historicalEntry := stakingtypes.NewHistoricalInfo(
		ctx.BlockHeader(), lastVals, sdk.DefaultPowerReduction,
	)

	// Set latest HistoricalInfo at current height
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
