package keeper

import (
	"context"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Params(ctx context.Context, _ *assetstype.QueryParamsRequest) (*assetstype.QueryParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	params, err := k.GetParams(c)
	if err != nil {
		return nil, err
	}
	return &assetstype.QueryParamsResponse{
		Params: params,
	}, nil
}

// QueClientChainInfoByIndex query client chain info by clientChainLzID
func (k Keeper) QueClientChainInfoByIndex(ctx context.Context, info *assetstype.QueryClientChainInfo) (*assetstype.ClientChainInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetClientChainInfoByIndex(c, info.ChainIndex)
}

// QueAllClientChainInfo query all client chain info that have been registered in exoCore
// the key of returned map is clientChainLzID, the value is the client chain info.
func (k Keeper) QueAllClientChainInfo(ctx context.Context, _ *assetstype.QueryAllClientChainInfo) (*assetstype.QueryAllClientChainInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllClientChainInfo(c)
	if err != nil {
		return nil, err
	}
	return &assetstype.QueryAllClientChainInfoResponse{AllClientChainInfos: allInfo}, nil
}

// QueStakingAssetInfo query the specified client chain asset info by inputting assetID
func (k Keeper) QueStakingAssetInfo(ctx context.Context, info *assetstype.QueryStakingAssetInfo) (*assetstype.StakingAssetInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakingAssetInfo(c, info.AssetID)
}

// QueAllStakingAssetsInfo query the info about all client chain assets that have been registered
func (k Keeper) QueAllStakingAssetsInfo(ctx context.Context, _ *assetstype.QueryAllStakingAssetsInfo) (*assetstype.QueryAllStakingAssetsInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllStakingAssetsInfo(c)
	if err != nil {
		return nil, err
	}
	return &assetstype.QueryAllStakingAssetsInfoResponse{AllStakingAssetsInfo: allInfo}, nil
}

// QueStakerAssetInfos query th state of all assets for a staker specified by stakerID
func (k Keeper) QueStakerAssetInfos(ctx context.Context, info *assetstype.QueryStakerAssetInfo) (*assetstype.QueryAssetInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	assetInfos, err := k.GetStakerAssetInfos(c, info.StakerID)
	if err != nil {
		return nil, err
	}
	return &assetstype.QueryAssetInfoResponse{AssetInfos: assetInfos}, nil
}

// QueStakerSpecifiedAssetAmount query the specified asset state of a staker, using stakerID and assetID as query parameters
func (k Keeper) QueStakerSpecifiedAssetAmount(ctx context.Context, req *assetstype.QuerySpecifiedAssetAmountReq) (*assetstype.StakerAssetInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakerSpecifiedAssetInfo(c, req.StakerID, req.AssetID)
}

// QueOperatorAssetInfos query th state of all assets for an operator specified by operator address
func (k Keeper) QueOperatorAssetInfos(ctx context.Context, infos *assetstype.QueryOperatorAssetInfos) (*assetstype.QueryOperatorAssetInfosResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(infos.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	assetInfos, err := k.GetOperatorAssetInfos(c, addr, nil)
	if err != nil {
		return nil, err
	}
	return &assetstype.QueryOperatorAssetInfosResponse{AssetInfos: assetInfos}, nil
}

// QueOperatorSpecifiedAssetAmount query the specified asset state of an operator, using operator address and assetID as query parameters
func (k Keeper) QueOperatorSpecifiedAssetAmount(ctx context.Context, req *assetstype.QueryOperatorSpecifiedAssetAmountReq) (*assetstype.OperatorAssetInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(req.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return k.GetOperatorSpecifiedAssetInfo(c, addr, req.AssetID)
}

// QueStakerExoCoreAddr outdated,will be deprecated
func (k Keeper) QueStakerExoCoreAddr(ctx context.Context, req *assetstype.QueryStakerExCoreAddr) (*assetstype.QueryStakerExCoreAddrResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	exoCoreAddr, err := k.GetStakerExoCoreAddr(c, req.Staker)
	if err != nil {
		return nil, err
	}
	return &assetstype.QueryStakerExCoreAddrResponse{ExoCoreAddr: exoCoreAddr}, nil
}
