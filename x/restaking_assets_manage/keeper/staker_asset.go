package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetStakerAssetInfos(ctx sdk.Context, stakerId string) (assetsInfo map[string]math.Uint, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) GetStakerSpecifiedAssetAmount(ctx sdk.Context, stakerId string, assetId string) (amount math.Uint, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) IncreaseStakerAssetsAmount(ctx sdk.Context, stakerId string, assetsAddAmount map[string]math.Uint) (err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) DecreaseStakerAssetsAmount(ctx sdk.Context, stakerId string, assetsSubAmount map[string]math.Uint) (err error) {
	//TODO implement me
	panic("implement me")
}
