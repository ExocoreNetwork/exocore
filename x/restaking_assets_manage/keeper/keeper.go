package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) Keeper {
	return Keeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

func (k Keeper) GetAllOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address) (optedInInfos map[string][]sdk.Address, err error) {
	// TODO implement me
	panic("implement me")
}

func (k Keeper) SetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, setInfo map[string]sdk.Address) (middleWares []sdk.Address, err error) {
	// TODO implement me
	panic("implement me")
}

// IRestakingAssetsManage interface will be implemented by restaking_assets_manage keeper
type IRestakingAssetsManage interface {
	SetClientChainInfo(ctx sdk.Context, info *restakingtype.ClientChainInfo) (err error)
	GetClientChainInfoByIndex(ctx sdk.Context, index uint64) (info *restakingtype.ClientChainInfo, err error)
	GetAllClientChainInfo(ctx sdk.Context) (infos map[uint64]*restakingtype.ClientChainInfo, err error)

	SetStakingAssetInfo(ctx sdk.Context, info *restakingtype.StakingAssetInfo) (err error)
	GetStakingAssetInfo(ctx sdk.Context, assetId string) (info *restakingtype.StakingAssetInfo, err error)
	GetAllStakingAssetsInfo(ctx sdk.Context) (allAssets map[string]*restakingtype.StakingAssetInfo, err error)

	GetStakerAssetInfos(ctx sdk.Context, stakerId string) (assetsInfo map[string]*restakingtype.StakerSingleAssetOrChangeInfo, err error)
	GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerId string, assetId string) (info *restakingtype.StakerSingleAssetOrChangeInfo, err error)
	UpdateStakerAssetState(ctx sdk.Context, stakerId string, assetId string, changeAmount restakingtype.StakerSingleAssetOrChangeInfo) (err error)

	GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address) (assetsInfo map[string]*restakingtype.OperatorSingleAssetOrChangeInfo, err error)
	GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetId string) (info *restakingtype.OperatorSingleAssetOrChangeInfo, err error)
	UpdateOperatorAssetState(ctx sdk.Context, operatorAddr sdk.Address, assetId string, changeAmount restakingtype.OperatorSingleAssetOrChangeInfo) (err error)

	// SetStakerExoCoreAddr handle the SetStakerExoCoreAddr txs from msg service
	SetStakerExoCoreAddr(ctx context.Context, addr *restakingtype.MsgSetExoCoreAddr) (*restakingtype.MsgSetExoCoreAddrResponse, error)
	GetStakerExoCoreAddr(ctx sdk.Context, stakerId string) (string, error)

	// GetOperatorAssetOptedInMiddleWare :the following three interfaces can be implemented in operator optedIn module
	GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetId string) (middleWares []sdk.Address, err error)
	GetAllOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address) (optedInInfos map[string][]sdk.Address, err error)
	SetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, setInfo map[string]sdk.Address) (middleWares []sdk.Address, err error)
}
