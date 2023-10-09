package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address) (assetsInfo map[string]math.Uint, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetOperatorSpecifiedAssetAmount(ctx sdk.Context, operatorAddr sdk.Address, assetId string) (amount math.Uint, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) IncreaseOperatorAssetsAmount(ctx sdk.Context, operatorAddr sdk.Address, assetsAddAmount map[string]math.Uint) (err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) DecreaseOperatorAssetsAmount(ctx sdk.Context, operatorAddr sdk.Address, assetsSubAmount map[string]math.Uint) (err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetOperatorAssetOptedInMiddleWare(ctx sdk.Context, operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error) {
	//TODO implement me
	panic("implement me")
}
