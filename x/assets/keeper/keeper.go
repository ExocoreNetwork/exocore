package keeper

import (
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
	assetstype.OracleKeeper
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	oracleKeeper assetstype.OracleKeeper,
) Keeper {
	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		OracleKeeper: oracleKeeper,
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

// IAssets interface will be implemented by assets keeper
type IAssets interface {
	SetClientChainInfo(ctx sdk.Context, info *assetstype.ClientChainInfo) (err error)
	GetClientChainInfoByIndex(ctx sdk.Context, index uint64) (info *assetstype.ClientChainInfo, err error)
	GetAllClientChainInfo(ctx sdk.Context) (infos map[uint64]*assetstype.ClientChainInfo, err error)

	SetStakingAssetInfo(ctx sdk.Context, info *assetstype.StakingAssetInfo) (err error)
	GetStakingAssetInfo(ctx sdk.Context, assetID string) (info *assetstype.StakingAssetInfo, err error)
	GetAllStakingAssetsInfo(ctx sdk.Context) (allAssets map[string]*assetstype.StakingAssetInfo, err error)

	GetStakerAssetInfos(ctx sdk.Context, stakerID string) (assetsInfo map[string]*assetstype.StakerAssetInfo, err error)
	GetStakerSpecifiedAssetInfo(ctx sdk.Context, stakerID string, assetID string) (info *assetstype.StakerAssetInfo, err error)
	UpdateStakerAssetState(ctx sdk.Context, stakerID string, assetID string, changeAmount assetstype.DeltaStakerSingleAsset) (err error)

	GetOperatorAssetInfos(ctx sdk.Context, operatorAddr sdk.Address, assetsFilter map[string]interface{}) (assetsInfo map[string]*assetstype.OperatorAssetInfo, err error)
	GetOperatorSpecifiedAssetInfo(ctx sdk.Context, operatorAddr sdk.Address, assetID string) (info *assetstype.OperatorAssetInfo, err error)
	UpdateOperatorAssetState(ctx sdk.Context, operatorAddr sdk.Address, assetID string, changeAmount assetstype.DeltaOperatorSingleAsset) (err error)
}
