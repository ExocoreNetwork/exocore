package keeper

import "cosmossdk.io/math"

func (k Keeper) GetStakerAssetInfos(reStakerId string) (assetsInfo map[string]math.Uint, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetStakerSpecifiedAssetAmount(stakerId string, assetId string) (amount math.Uint, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) IncreaseStakerAssetsAmount(stakerId string, assetsAddAmount map[string]math.Uint) (err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) DecreaseStakerAssetsAmount(stakerId string, assetsSubAmount map[string]math.Uint) (err error) {
	//TODO implement me
	panic("implement me")
}
