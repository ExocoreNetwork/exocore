package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeforeEpochStart: noop, We don't need to do anything here
func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {
}

// AfterEpochEnd avs Handling reward and slash  at the end of each epoch end
func (k Keeper) AfterEpochEnd(
	ctx sdk.Context, epochIdentifier string, epochNumber int64,
) {
	// get all the avs address bypass the epoch end
	epochEndAVS, err := k.GetEpochEndAVSs(ctx, epochIdentifier, epochNumber)
	if err != nil {
		k.Logger(ctx).Error(
			"epochEndAVS got error",
			"epochEndAVS", epochEndAVS,
		)
		return
	}

	//TODO:Handling reward and slash
	avsInfo, err := k.GetAVSInfo(ctx, "")
	assetId := avsInfo.Info.AssetId
	operatorAddress := avsInfo.Info.OperatorAddress
	if err != nil {
		k.Logger(ctx).Error(
			"get avsInfo error",
			"avsInfo err", err,
		)
		return
	}
	for _, operator := range operatorAddress {
		_, err := sdk.AccAddressFromBech32(operator)
		if err != nil {
			k.Logger(ctx).Error(
				"get operatorInfo error",
				"operatorInfo err", err,
			)
			return
		}
		for _, asset := range assetId {
			_, err := k.assetsKeeper.GetStakingAssetInfo(ctx, asset)
			if err != nil {
				k.Logger(ctx).Error(
					"get assetInfo error",
					"assetInfo err", err,
				)
				return
			}

		}
	}
}
