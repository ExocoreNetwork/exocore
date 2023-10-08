package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetOperatorAssetInfos(operatorAddr sdk.Address) (assetsInfo map[string]math.Uint, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetOperatorSpecifiedAssetAmount(operatorAddr sdk.Address, assetId string) (amount math.Uint, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) IncreaseOperatorAssetsAmount(operatorAddr sdk.Address, assetsAddAmount map[string]math.Uint) (err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) DecreaseOperatorAssetsAmount(operatorAddr sdk.Address, assetsSubAmount map[string]math.Uint) (err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error) {
	//TODO implement me
	panic("implement me")
}
