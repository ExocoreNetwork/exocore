package keeper

import types2 "github.com/exocore/x/restaking_assets_manage/types"

func (k Keeper) SetStakingAssetInfo(info *types2.StakingAssetInfo) (exoCoreAssetIndex uint64, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetStakingAssetInfo(assetId string) (info *types2.StakingAssetInfo, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetAllStakingAssetsInfo() (allAssets map[string]*types2.StakingAssetInfo, err error) {
	//TODO implement me
	panic("implement me")
}
