package keeper

import (
	"cosmossdk.io/math"
	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RegisterClientChain(ctx sdk.Context, clientChain *restakingtype.ClientChainInfo) error {
	// check if the client chain has been registered
	isExist := k.ClientChainInfoIsExist(ctx, clientChain.LayerZeroChainId)
	if isExist {
		return restakingtype.ErrClientChainIsExist
	}

	return k.SetClientChainInfo(ctx, clientChain)
}

func (k Keeper) DeregisterClientChain(ctx sdk.Context, clientChainId uint64) error {
	// todo: Check whether all state related to the input clientChainId hasn't been used, for example all assets from the
	// client chain have been deregistered.

	// delete the client chain info
	isExist := k.ClientChainInfoIsExist(ctx, clientChainId)
	if !isExist {
		return restakingtype.ErrNoClientChainKey
	}
	return k.DeleteClientChainInfo(ctx, clientChainId)
}

func (k Keeper) RegisterAsset(ctx sdk.Context, asset *restakingtype.ClientChainTokenInfo) error {
	_, assetID := restakingtype.GetStakeIDAndAssetIdFromStr(asset.LayerZeroChainId, "", asset.Address)
	// check if the asset has been registered
	isExist := k.IsStakingAsset(ctx, assetID)
	if isExist {
		return restakingtype.ErrAssetIsExist
	}

	return k.SetStakingAssetInfo(ctx, &restakingtype.StakingAssetInfo{
		AssetBasicInfo:     asset,
		StakingTotalAmount: math.NewInt(0),
	})
}

func (k Keeper) DeregisterAsset(ctx sdk.Context, assetID string) error {
	// todo: Check whether all state related to the asset hasn't been used
	// Only when all stakers have withdrawn this asset, could it be deregistered.
	isExist := k.IsStakingAsset(ctx, assetID)
	if !isExist {
		return restakingtype.ErrNoClientChainAssetKey
	}
	return k.DeleteStakingAssetInfo(ctx, assetID)
}
