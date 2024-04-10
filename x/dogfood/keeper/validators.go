// This file is a duplicate of the subscriber module's validators file with minor changes.
// The function ApplyValidatorChanges can likely be carved out into a shared package with
// the appchain module.

package keeper

import (
	"sort"
	"time"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
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
	ctx sdk.Context, changes []abci.ValidatorUpdate, valSetID uint64,
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

	// store the updated validator set against the provided validator set id
	lastVals := types.Validators{}
	for _, v := range k.GetAllExocoreValidators(ctx) {
		// we stored the validators above, so this will never fail.
		pubkey, _ := v.ConsPubKey() // #nosec G703
		// #nosec G703 // This calls NewAnyWithValue internally, which we have already done.
		val, _ := stakingtypes.NewValidator(
			nil, pubkey, stakingtypes.Description{},
		)
		// Set validator to bonded status
		val.Status = stakingtypes.Bonded
		// Compute tokens from voting power
		val.Tokens = sdk.TokensFromConsensusPower(v.Power, sdk.DefaultPowerReduction)
		lastVals.List = append(lastVals.List, val)
	}
	k.setValidatorSet(ctx, valSetID, &lastVals)
	// this validator set is effective as of the next block, so use height + 1.
	// this statement is true for genesis as well, since ctx.BlockHeight() is
	// reported as 0 during InitGenesis.
	k.setValidatorSetID(ctx, ctx.BlockHeight()+1, valSetID)
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

// GetHistoricalInfo gets the historical info at a given height. It is part of the
// implementation of the staking keeper expected by IBC.
func (k Keeper) GetHistoricalInfo(
	ctx sdk.Context, height int64,
) (stakingtypes.HistoricalInfo, bool) {
	headerSubset, found := k.getBlockHeader(ctx, height)
	if !found {
		// only panic in the case of an unmarshal error
		return stakingtypes.HistoricalInfo{}, false
	}
	valSetID, found := k.getValidatorSetID(ctx, height)
	if !found {
		// only panic in the case of an unmarshal error
		return stakingtypes.HistoricalInfo{}, false
	}
	valSet, found := k.getValidatorSet(ctx, valSetID)
	if !found {
		// only panic in the case of an unmarshal error
		return stakingtypes.HistoricalInfo{}, false
	}
	header := tmproto.Header{
		Time:               headerSubset.Time,
		NextValidatorsHash: headerSubset.NextValidatorsHash,
		AppHash:            headerSubset.AppHash,
	}
	return stakingtypes.NewHistoricalInfo(
		header, stakingtypes.Validators(valSet.GetList()), sdk.DefaultPowerReduction,
	), true
}

// setValidatorSet sets the validator set at a given id. This is
// (intentionally) not exported in the genesis state. It can therefore
// be a private function.
func (k Keeper) setValidatorSet(
	ctx sdk.Context, id uint64, vs *types.Validators,
) {
	store := ctx.KVStore(k.storeKey)
	key := types.ValidatorSetKey(id)
	value := k.cdc.MustMarshal(vs)
	store.Set(key, value)
}

// getValidatorSet gets the validator set at a given id.
func (k Keeper) getValidatorSet(
	ctx sdk.Context, id uint64,
) (*types.Validators, bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.ValidatorSetKey(id)
	if !store.Has(key) {
		return nil, false
	}
	value := store.Get(key)
	var hi types.Validators
	k.cdc.MustUnmarshal(value, &hi)
	return &hi, true
}

// deleteValidatorSet deletes the validator set at a given id.
func (k Keeper) deleteValidatorSet(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.ValidatorSetKey(id)
	store.Delete(key)
}

// TrackHistoricalInfo saves the latest historical info and deletes the ones eligible to be
// pruned. The historical info is stored in two parts: one is the header and the other is the
// validator set. Within an epoch, the validator set will only change if there is a slashing
// event. Otherwise, it is constant. The header, however, will change at every block. Since
// the Cosmos SDK does not allow for the retrieval of a past block header, we store the header
// ourselves in this function. The validator set is stored when it changes at the end of an
// epoch or at a slashing event in the corresponding functions. The function is called within
// the EndBlock of the module, so it is kept public.
func (k Keeper) TrackHistoricalInfo(ctx sdk.Context) {
	// Get the number of historical entries to persist, as the number of block heights.
	// #nosec G701 // uint32 fits into int64 always.
	numHistoricalEntries := int64(
		k.GetHistoricalEntries(ctx),
	)

	// we are deleting headers, say, from, 0 to 999 at block 1999
	// for these headers, we must find the corresponding validator set ids to delete.
	// they must be only deleted if no other block is using them.
	lastDeletedID := uint64(0) // contract: starts from 1.
	for i := ctx.BlockHeight() - numHistoricalEntries; i >= 0; i-- {
		_, found := k.getBlockHeader(ctx, i)
		if found {
			// because they are deleted together, and saved one after the other,
			// since the block header exists, so must the validator set id.
			lastDeletedID, _ = k.getValidatorSetID(ctx, i+1)
			// clear both the header and the mapping
			k.deleteBlockHeader(ctx, i)
			k.deleteValidatorSetID(ctx, i)
		} else {
			break
		}
	}
	// even if numHistoricalEntries is 0, this will work because it is called after the
	// validatorSetID for height + 1 is stored.
	// on the opposite side of things, if numHistoricalEntries is too large, currentID
	// will be 0, and the loop will not run.
	currentID, _ := k.getValidatorSetID(ctx,
		ctx.BlockHeight()-numHistoricalEntries+1,
	)
	// lastDeletedID will be the lowest deleted id since we are working backwards
	// from the latest height to the oldest height. this, and upto but not including
	// currentID, are the ids to delete.
	for i := lastDeletedID; i < currentID; i++ {
		k.deleteValidatorSet(ctx, i)
	}

	// if there is no need to persist historicalInfo, return.
	if numHistoricalEntries == 0 {
		return
	}

	// store the header
	k.storeBlockHeader(ctx)

	// we have stored:
	// outside of TrackHistoricalInfo: ValidatorSetID for height, and the validator set.
	// within TrackHistoricalInfo: the header.
	// this is enough information to answer the GetHistoricalInfo query.
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

// getValidatorSetID returns the identifier of the validator set at a given height.
// It is used to "share" the validator set entries across multiple heights within an epoch.
// Typically, the validator set should change only at the end of an epoch. However, in the
// case of a slashing occurrence, the validator set may change within an epoch. This info
// is not exported in the genesis, so it is private.
func (k Keeper) getValidatorSetID(ctx sdk.Context, height int64) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	key, ok := types.ValidatorSetIDKey(height)
	if !ok {
		return 0, false
	}
	value := store.Get(key)
	if value == nil {
		return 0, false
	}
	return sdk.BigEndianToUint64(value), true
}

// setValidatorSetID sets the identifier of the validator set at a given height.
func (k Keeper) setValidatorSetID(ctx sdk.Context, height int64, id uint64) {
	store := ctx.KVStore(k.storeKey)
	// SetValidatorSetID is called with ctx.BlockHeight which Cosmos
	// guarantees to be positive.
	key, _ := types.ValidatorSetIDKey(height)
	value := sdk.Uint64ToBigEndian(id)
	store.Set(key, value)
}

// deleteValidatorSetID deletes the identifier of the validator set at a given height.
func (k Keeper) deleteValidatorSetID(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	// i don't care if we delete a key that is meaningless
	key, _ := types.ValidatorSetIDKey(height)
	store.Delete(key)
}

// getBlockHeader returns the block header at a given height. It is called during IBC's
// GetHistoricalInfo call. Since it is an external call, we will validate that height
// should be non-negative. This info is not exported in the genesis, so it is private.
func (k Keeper) getBlockHeader(ctx sdk.Context, height int64) (types.HeaderSubset, bool) {
	store := ctx.KVStore(k.storeKey)
	key, ok := types.HeaderKey(height)
	if !ok {
		return types.HeaderSubset{}, false
	}
	value := store.Get(key)
	if value == nil {
		return types.HeaderSubset{}, false
	}
	var header types.HeaderSubset
	k.cdc.MustUnmarshal(value, &header)
	return header, true
}

// deleteBlockHeader deletes the block header at a given height.
func (k Keeper) deleteBlockHeader(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	// i don't care if we delete a key that is meaningless
	key, _ := types.HeaderKey(height)
	store.Delete(key)
}

// storeBlockHeader stores the block header subset as of the current height.
func (k Keeper) storeBlockHeader(ctx sdk.Context) {
	// ctx.BlockHeight() is positive so we don't need to validate
	key, _ := types.HeaderKey(ctx.BlockHeight())
	sdkHeader := ctx.BlockHeader()
	header := types.HeaderSubset{
		Time:               sdkHeader.Time,
		NextValidatorsHash: sdkHeader.NextValidatorsHash,
		AppHash:            sdkHeader.GetAppHash(),
	}
	store := ctx.KVStore(k.storeKey)
	value := k.cdc.MustMarshal(&header)
	store.Set(key, value)
}
