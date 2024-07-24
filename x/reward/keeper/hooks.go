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
	epochEndAVS := k.avsKeeper.GetEpochEndAVSs(ctx, epochIdentifier, epochNumber)

	pool := k.getPool(ctx, types.ModuleName)
	// distribute the reward to the avs accordingly
	ForEach(epochEndAVS, func(p string) {
		avsInfo, err := pool.k.avsKeeper.GetAVSInfo(ctx, p)
		if err != nil {
			k.Logger(ctx).Error(
				"get avsInfo error",
				"avsInfo err", err,
			)
			return
		}
		assetID := avsInfo.Info.AssetIDs
		ownerAddress := avsInfo.Info.AvsOwnerAddress

		for _, operator := range ownerAddress {
			opAddr, err := sdk.AccAddressFromBech32(operator)
			if err != nil {
				k.Logger(ctx).Error(
					"get operatorInfo error",
					"operatorInfo err", err,
				)
				return
			}
			for _, asset := range assetID {
				assetInfo, err := k.assetsKeeper.GetStakingAssetInfo(ctx, asset)
				if err != nil {
					k.Logger(ctx).Error(
						"get assetInfo error",
						"assetInfo err", err,
					)
					return
				}
				if k.assetsKeeper.IsOperatorAssetExist(ctx, opAddr, asset) {
					coin := sdk.Coin{
						Denom:  assetInfo.AssetBasicInfo.Symbol,
						Amount: sdk.NewInt(avsInfo.Info.AssetRewardAmountEpochBasis[asset]),
					}
					pool.AddReward(p, coin)
				}
			}
		}
	})
}

// ForEach apply the function on every element within the slice
func ForEach[T any](source []T, f func(T)) {
	for i := range source {
		f(source[i])
	}
}
