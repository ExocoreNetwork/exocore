package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

func (k Keeper) SetStakingAssetInfo(ctx sdk.Context, info *types2.StakingAssetInfo) (exoCoreAssetIndex uint64, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetStakingAssetInfo(ctx sdk.Context, assetId string) (info *types2.StakingAssetInfo, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetAllStakingAssetsInfo(ctx sdk.Context) (allAssets map[string]*types2.StakingAssetInfo, err error) {
	//TODO implement me
	panic("implement me")
}
