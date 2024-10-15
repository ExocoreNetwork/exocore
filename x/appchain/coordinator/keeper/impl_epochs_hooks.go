package keeper

import (
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EpochsHooksWrapper is the wrapper structure that implements the epochs hooks for the dogfood
// keeper.
type EpochsHooksWrapper struct {
	keeper *Keeper
}

// Interface guard
var _ epochstypes.EpochHooks = EpochsHooksWrapper{}

// EpochsHooks returns the epochs hooks wrapper. It follows the "accept interfaces, return
// concretes" pattern.
func (k *Keeper) EpochsHooks() EpochsHooksWrapper {
	return EpochsHooksWrapper{k}
}

// AfterEpochEnd is called after an epoch ends. It is called during the BeginBlock function.
func (wrapper EpochsHooksWrapper) AfterEpochEnd(
	ctx sdk.Context, identifier string, epoch int64,
) {
	// start any chains that are due to start, by creating their genesis state.
	wrapper.keeper.ActivateScheduledChains(ctx, identifier, epoch)

	// slashing is applied during the epoch, so we don't have to do anything about that here.
	// note that slashing should flow through to this keeper via a hook and the impact
	// should be applied to the validator set. first, it should freeze the oracle round,
	// then, it should calculate the USD power, then, it should find the new x/dogfood
	// validator set and lastly, it should find the new appchain validator set for the
	// impacted chains.

	// next, we remove any chains that didn't respond in time: either to the validator
	// set update or to the initialization protocol. the removal is undertaken before
	// generating the validator set update to save resources.
	wrapper.keeper.RemoveTimedoutSubscribers(ctx, identifier, epoch)

	// last, we iterate over the active list and queue the validator set update for them.
	// interchain-security does this in EndBlock, but we can do it now because our validator
	// set is independent of the coordinator chain's.
	wrapper.keeper.QueueValidatorUpdatesForEpochID(ctx, identifier, epoch)
	// send the queued validator updates. the `epoch` is used for scheduling the VSC timeouts
	// and nothing else. it has no bearing on the actual validator set.
	wrapper.keeper.SendQueuedValidatorUpdates(ctx, epoch)
}

// BeforeEpochStart is called before an epoch starts.
func (wrapper EpochsHooksWrapper) BeforeEpochStart(
	sdk.Context, string, int64,
) {
	// no-op
}
