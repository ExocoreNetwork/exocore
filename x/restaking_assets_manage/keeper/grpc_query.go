package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueClientChainInfoByIndex(ctx context.Context, info *types2.QueryClientChainInfo) (*types2.ClientChainInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetClientChainInfoByIndex(c, info.ChainIndex)
}

func (k Keeper) QueAllClientChainInfo(ctx context.Context, info *types2.QueryAllClientChainInfo) (*types2.QueryAllClientChainInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllClientChainInfo(c)
	if err != nil {
		return nil, err
	}
	return &types2.QueryAllClientChainInfoResponse{AllClientChainInfos: allInfo}, nil
}

func (k Keeper) QueStakingAssetInfo(ctx context.Context, info *types2.QueryStakingAssetInfo) (*types2.StakingAssetInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakingAssetInfo(c, info.AssetId)
}

func (k Keeper) QueAllStakingAssetsInfo(ctx context.Context, info *types2.QueryAllStakingAssetsInfo) (*types2.QueryAllStakingAssetsInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllStakingAssetsInfo(c)
	if err != nil {
		return nil, err
	}
	return &types2.QueryAllStakingAssetsInfoResponse{AllStakingAssetsInfo: allInfo}, nil
}

func (k Keeper) QueStakerAssetInfos(ctx context.Context, info *types2.QueryStakerAssetInfo) (*types2.QueryAssetInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	assetInfos, err := k.GetStakerAssetInfos(c, info.StakerId)
	if err != nil {
		return nil, err
	}
	return &types2.QueryAssetInfoResponse{AssetInfos: assetInfos}, nil
}

func (k Keeper) QueStakerSpecifiedAssetAmount(ctx context.Context, req *types2.QuerySpecifiedAssetAmountReq) (*types2.StakerSingleAssetOrChangeInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakerSpecifiedAssetInfo(c, req.StakerId, req.AssetId)
}

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

func (k Keeper) QueOperatorSpecifiedAssetAmount(ctx context.Context, req *types2.QueryOperatorSpecifiedAssetAmountReq) (*types2.OperatorSingleAssetOrChangeInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(req.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return k.GetOperatorSpecifiedAssetInfo(c, addr, req.AssetId)
}

func (k Keeper) QueStakerExoCoreAddr(ctx context.Context, req *types2.QueryStakerExCoreAddr) (*types2.QueryStakerExCoreAddrResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	exoCoreAddr, err := k.GetStakerExoCoreAddr(c, req.StakerId)
	if err != nil {
		return nil, err
	}
	return &types2.QueryStakerExCoreAddrResponse{ExCoreAddr: exoCoreAddr}, nil
}
