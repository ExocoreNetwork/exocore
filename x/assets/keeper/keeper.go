package keeper

import (
	"context"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// GetAllOperatorAssetOptedInMiddleWare This function should be implemented in the operator opt-in module
func (k Keeper) GetAllOperatorAssetOptedInMiddleWare(sdk.Address) (optedInInfos map[string][]sdk.Address, err error) {
	// TODO implement me
	panic("implement me")
}

// SetOperatorAssetOptedInMiddleWare This function should be implemented in the operator opt-in module
func (k Keeper) SetOperatorAssetOptedInMiddleWare(sdk.Address, map[string]sdk.Address) (middleWares []sdk.Address, err error) {
	// TODO implement me
	panic("implement me")
}

// IRestakingAssetsManage interface will be implemented by assets keeper
type IRestakingAssetsManage interface {
	SetClientChainInfo(ctx sdk.Context, info *assetstype.ClientChainInfo) (err error)
	GetClientChainInfoByIndex(ctx sdk.Context, index uint64) (info *assetstype.ClientChainInfo, err error)
	GetAllClientChainInfo(ctx sdk.Context) (infos map[uint64]*assetstype.ClientChainInfo, err error)

	SetStakingAssetInfo(ctx sdk.Context, info *assetstype.StakingAssetInfo) (err error)
	GetStakingAssetInfo(ctx sdk.Context, assetID string) (info *assetstype.StakingAssetInfo, err error)
	GetAllStakingAssetsInfo(ctx sdk.Context) (allAssets map[string]*assetstype.StakingAssetInfo, err error)

	GetStakerAssetInfos(ctx sdk.Context, stakerID string) (assetsInfo map[string]*assetstype.StakerSingleAssetInfo, err error)
	GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerID string, assetID string) (info *assetstype.StakerSingleAssetInfo, err error)
	UpdateStakerAssetState(ctx sdk.Context, stakerID string, assetID string, changeAmount assetstype.StakerSingleAssetChangeInfo) (err error)

	GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address) (assetsInfo map[string]*assetstype.OperatorSingleAssetInfo, err error)
	GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetID string) (info *assetstype.OperatorSingleAssetInfo, err error)
	UpdateOperatorAssetState(ctx sdk.Context, operatorAddr sdk.Address, assetID string, changeAmount assetstype.OperatorSingleAssetChangeInfo) (err error)

	// SetStakerExoCoreAddr handle the SetStakerExoCoreAddr txs from msg service
	SetStakerExoCoreAddr(ctx context.Context, addr *assetstype.MsgSetExoCoreAddr) (*assetstype.MsgSetExoCoreAddrResponse, error)
	GetStakerExoCoreAddr(ctx sdk.Context, stakerID string) (string, error)

	// GetOperatorAssetOptedInMiddleWare :the following three interfaces should be implemented in operator opt-in module
	GetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, assetID string) (middleWares []sdk.Address, err error)
	GetAllOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address) (optedInInfos map[string][]sdk.Address, err error)
	SetOperatorAssetOptedInMiddleWare(operatorAddr sdk.Address, setInfo map[string]sdk.Address) (middleWares []sdk.Address, err error)
}
