package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeforeEpochStart: noop, We don't need to do anything here
func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {
}

// AfterEpochEnd distribute the reward at the end of each epoch end
func (k Keeper) AfterEpochEnd(
	ctx sdk.Context, epochIdentifier string, epochNumber int64,
) {
	expEpochIdentifier := k.GetEpochIdentifier(ctx)
	if epochIdentifier != expEpochIdentifier {
		k.Logger(ctx).Error(
			"epochIdentifier didn't equal to expEpochIdentifier",
			"epochIdentifier", epochIdentifier,
		)
		return
	}

	// get all the avs address bypass the epoch end
	epochEndAVS, err := k.avsKeeper.GetEpochEndAVSs(ctx)
	if err != nil {
		k.Logger(ctx).Error(
			"epochEndAVS got error",
			"epochEndAVS", epochEndAVS,
		)
		return
	}
	pool := k.getPool(ctx, types.ModuleName)
	// distribute the reward to the avs accordingly
	ForEach(epochEndAVS, func(p string) {
		if err := pool.ReleaseRewards(p); err != nil {
			k.Logger(ctx).Error(
				"release reward error",
				"error message", err,
			)
			return
		}
	})
}

// ForEach apply the function on every element within the slice
func ForEach[T any](source []T, f func(T)) {
	for i := range source {
		f(source[i])
	}
}
