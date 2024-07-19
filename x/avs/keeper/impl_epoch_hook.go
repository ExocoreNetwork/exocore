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
	ctx sdk.Context, epochIdentifier string, epochNumber int64,
) {
	// get all the avs address bypass the epoch end
	epochEndAVS, err := wrapper.keeper.GetEpochEndAVSs(ctx, epochIdentifier, epochNumber)
	if err != nil {
		wrapper.keeper.Logger(ctx).Error(
			"epochEndAVS got error",
			"epochEndAVS", epochEndAVS,
		)
		return
	}

	taskChallengeEpochEndAVS, err := wrapper.keeper.GetTaskChallengeEpochEndAVSs(ctx, epochIdentifier, epochNumber)
	if err != nil {
		wrapper.keeper.Logger(ctx).Error(
			"epochEndAVS got error",
			"epochEndAVS", taskChallengeEpochEndAVS,
		)
		return
	}
}

// BeforeEpochStart is called before an epoch starts.
func (wrapper EpochsHooksWrapper) BeforeEpochStart(
	sdk.Context, string, int64,
) {
}
