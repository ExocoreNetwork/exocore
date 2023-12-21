package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueClientChainInfoByIndex query client chain info by clientChainLzId
func (k Keeper) QueClientChainInfoByIndex(ctx context.Context, info *types2.QueryClientChainInfo) (*types2.ClientChainInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetClientChainInfoByIndex(c, info.ChainIndex)
}

// QueAllClientChainInfo query all client chain info that have been registered in exoCore
// the key of returned map is clientChainLzId, the value is the client chain info.
func (k Keeper) QueAllClientChainInfo(ctx context.Context, info *types2.QueryAllClientChainInfo) (*types2.QueryAllClientChainInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllClientChainInfo(c)
	if err != nil {
		return nil, err
	}
	return &types2.QueryAllClientChainInfoResponse{AllClientChainInfos: allInfo}, nil
}

// QueStakingAssetInfo query the specified client chain asset info by inputting assetId
func (k Keeper) QueStakingAssetInfo(ctx context.Context, info *types2.QueryStakingAssetInfo) (*types2.StakingAssetInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakingAssetInfo(c, info.AssetId)
}

// QueAllStakingAssetsInfo query the info about all client chain assets that have been registered
func (k Keeper) QueAllStakingAssetsInfo(ctx context.Context, info *types2.QueryAllStakingAssetsInfo) (*types2.QueryAllStakingAssetsInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllStakingAssetsInfo(c)
	if err != nil {
		return nil, err
	}
	return &types2.QueryAllStakingAssetsInfoResponse{AllStakingAssetsInfo: allInfo}, nil
}

// QueStakerAssetInfos query th state of all assets for a staker specified by stakerId
func (k Keeper) QueStakerAssetInfos(ctx context.Context, info *types2.QueryStakerAssetInfo) (*types2.QueryAssetInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	assetInfos, err := k.GetStakerAssetInfos(c, info.StakerId)
	if err != nil {
		return nil, err
	}
	return &types2.QueryAssetInfoResponse{AssetInfos: assetInfos}, nil
}

// QueStakerSpecifiedAssetAmount query the specified asset state of a staker, using stakerId and assetId as query parameters
func (k Keeper) QueStakerSpecifiedAssetAmount(ctx context.Context, req *types2.QuerySpecifiedAssetAmountReq) (*types2.StakerSingleAssetOrChangeInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakerSpecifiedAssetInfo(c, req.StakerId, req.AssetId)
}

// QueOperatorAssetInfos query th state of all assets for an operator specified by operator address
func (k Keeper) QueOperatorAssetInfos(ctx context.Context, infos *types2.QueryOperatorAssetInfos) (*types2.QueryOperatorAssetInfosResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(infos.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	assetInfos, err := k.GetOperatorAssetInfos(c, addr)
	if err != nil {
		return nil, err
	}
	return &types2.QueryOperatorAssetInfosResponse{AssetInfos: assetInfos}, nil
}

// QueOperatorSpecifiedAssetAmount query the specified asset state of an operator, using operator address and assetId as query parameters
func (k Keeper) QueOperatorSpecifiedAssetAmount(ctx context.Context, req *types2.QueryOperatorSpecifiedAssetAmountReq) (*types2.OperatorSingleAssetOrChangeInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(req.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return k.GetOperatorSpecifiedAssetInfo(c, addr, req.AssetId)
}

// QueStakerExoCoreAddr outdated,will be deprecated
func (k Keeper) QueStakerExoCoreAddr(ctx context.Context, req *types2.QueryStakerExCoreAddr) (*types2.QueryStakerExCoreAddrResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	exoCoreAddr, err := k.GetStakerExoCoreAddr(c, req.StakerId)
	if err != nil {
		return nil, err
	}
	return &types2.QueryStakerExCoreAddrResponse{ExCoreAddr: exoCoreAddr}, nil
}
