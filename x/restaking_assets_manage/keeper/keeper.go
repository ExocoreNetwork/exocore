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

func (k Keeper) GetAllOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address) (optedInInfos map[string][]sdk.Address, err error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) SetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, setInfo map[string]sdk.Address) (middleWares []sdk.Address, err error) {
	//TODO implement me
	panic("implement me")
}

// IReStakingAssetsManage interface will be implemented by restaking_assets_manage keeper
type IReStakingAssetsManage interface {
	SetClientChainInfo(ctx sdk.Context, info *types2.ClientChainInfo) (exoCoreChainIndex uint64, err error)
	GetClientChainInfoByIndex(ctx sdk.Context, index uint64) (info *types2.ClientChainInfo, err error)
	GetAllClientChainInfo(ctx sdk.Context) (infos map[uint64]*types2.ClientChainInfo, err error)

	SetStakingAssetInfo(ctx sdk.Context, info *types2.StakingAssetInfo) (exoCoreAssetIndex uint64, err error)
	GetStakingAssetInfo(ctx sdk.Context, assetId string) (info *types2.StakingAssetInfo, err error)
	GetAllStakingAssetsInfo() (ctx sdk.Context, allAssets map[string]*types2.StakingAssetInfo, err error)

	GetStakerAssetInfos(ctx sdk.Context, stakerId string) (assetsInfo map[string]math.Uint, err error)
	GetStakerSpecifiedAssetAmount(ctx sdk.Context, stakerId string, assetId string) (amount math.Uint, err error)
	IncreaseStakerAssetsAmount(ctx sdk.Context, stakerId string, assetsAddAmount map[string]math.Uint) (err error)
	DecreaseStakerAssetsAmount(ctx sdk.Context, stakerId string, assetsSubAmount map[string]math.Uint) (err error)

	GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address) (assetsInfo map[string]math.Uint, err error)
	GetOperatorSpecifiedAssetAmount(ctx sdk.Context, operatorAddr sdk.Address, assetId string) (amount math.Uint, err error)
	IncreaseOperatorAssetsAmount(ctx sdk.Context, operatorAddr sdk.Address, assetsAddAmount map[string]math.Uint) (err error)
	DecreaseOperatorAssetsAmount(ctx sdk.Context, operatorAddr sdk.Address, assetsSubAmount map[string]math.Uint) (err error)

	// SetStakerExoCoreAddr handle the SetStakerExoCoreAddr txs from msg service
	SetStakerExoCoreAddr(ctx sdk.Context, addr *types2.MsgSetExoCoreAddr) (*types2.MsgSetExoCoreAddrResponse, error)
	GetStakerExoCoreAddr(ctx sdk.Context, stakerId string) (string, error)

	// GetOperatorAssetOptedInMiddleWare :the following three interfaces can be implemented in operator optedIn module
	GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error)
	GetAllOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address) (optedInInfos map[string][]sdk.Address, err error)
	SetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, setInfo map[string]sdk.Address) (middleWares []sdk.Address, err error)
}
