// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
}

// IReStakingAssetsManage interface will be implemented by restaking_assets_manage keeper
type IReStakingAssetsManage interface {
	SetClientChainInfo(info *types2.ClientChainInfo) (exoCoreChainIndex uint64, err error)
	GetClientChainInfoByIndex(exoCoreChainIndex uint64) (info types2.ClientChainInfo, err error)
	GetAllClientChainInfo() (infos map[uint64]types2.ClientChainInfo, err error)

	SetReStakingAssetInfo(info *types2.ReStakingAssetInfo) (exoCoreAssetIndex uint64, err error)
	GetReStakingAssetInfo(assetId string) (info types2.ReStakingAssetInfo, err error)
	GetAllReStakingAssetsInfo() (allAssets map[string]types2.ReStakingAssetInfo, err error)

	GetReStakerAssetInfos(reStakerId string) (assetsInfo map[string]math.Uint, err error)
	GetReStakerSpecifiedAssetAmount(reStakerId string, assetId string) (amount math.Uint, err error)
	IncreaseReStakerAssetsAmount(reStakerId string, assetsAddAmount map[string]math.Uint) (err error)
	DecreaseReStakerAssetsAmount(reStakerId string, assetsSubAmount map[string]math.Uint) (err error)

	GetOperatorAssetInfos(operatorAddr sdk.Address) (assetsInfo map[string]math.Uint, err error)
	GetOperatorSpecifiedAssetAmount(operatorAddr sdk.Address, assetId string) (amount math.Uint, err error)
	IncreaseOperatorAssetsAmount(operatorAddr sdk.Address, assetsAddAmount map[string]math.Uint) (err error)
	DecreaseOperatorAssetsAmount(operatorAddr sdk.Address, assetsSubAmount map[string]math.Uint) (err error)
	GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error)

	// GetAllOperatorAssetOptedInMiddleWare can also be implemented in operator optedIn module
	GetAllOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address) (optedInInfos map[string][]sdk.Address, err error)
	SetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, setInfo map[string]sdk.Address) (middleWares []sdk.Address, err error)
}
