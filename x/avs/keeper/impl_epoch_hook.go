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

	taskResponseEpochEndAVSepochEndAVS, err := wrapper.keeper.GetTaskChallengeEpochEndAVSs(ctx, epochIdentifier, epochNumber)
	if err != nil {
		wrapper.keeper.Logger(ctx).Error(
			"epochEndAVS got error",
			"epochEndAVS", taskResponseEpochEndAVSepochEndAVS,
		)
		return
	}

	// TODO:Handling reward and slash
	avsInfo, err := wrapper.keeper.GetAVSInfo(ctx, "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr")
	if avsInfo.Info == nil {
		return
	}
	assetID := avsInfo.Info.AssetId
	if err != nil {
		wrapper.keeper.Logger(ctx).Error(
			"get avsInfo error",
			"avsInfo err", err,
		)
		return
	}

	for _, asset := range assetID {
		_, err := wrapper.keeper.assetsKeeper.GetStakingAssetInfo(ctx, asset)
		if err != nil {
			wrapper.keeper.Logger(ctx).Error(
				"get assetInfo error",
				"assetInfo err", err,
			)
			return
		}

	}
}

// BeforeEpochStart is called before an epoch starts.
func (wrapper EpochsHooksWrapper) BeforeEpochStart(
	sdk.Context, string, int64,
) {
	// no-op
}
