package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueClientChainInfoByIndex(ctx context.Context, info *restakingtype.QueryClientChainInfo) (*restakingtype.ClientChainInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetClientChainInfoByIndex(c, info.ChainIndex)
}

func (k Keeper) QueAllClientChainInfo(ctx context.Context, info *restakingtype.QueryAllClientChainInfo) (*restakingtype.QueryAllClientChainInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllClientChainInfo(c)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryAllClientChainInfoResponse{AllClientChainInfos: allInfo}, nil
}

func (k Keeper) QueStakingAssetInfo(ctx context.Context, info *restakingtype.QueryStakingAssetInfo) (*restakingtype.StakingAssetInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakingAssetInfo(c, info.AssetId)
}

func (k Keeper) QueAllStakingAssetsInfo(ctx context.Context, info *restakingtype.QueryAllStakingAssetsInfo) (*restakingtype.QueryAllStakingAssetsInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllStakingAssetsInfo(c)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryAllStakingAssetsInfoResponse{AllStakingAssetsInfo: allInfo}, nil
}

func (k Keeper) QueStakerAssetInfos(ctx context.Context, info *restakingtype.QueryStakerAssetInfo) (*restakingtype.QueryAssetInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	assetInfos, err := k.GetStakerAssetInfos(c, info.StakerId)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryAssetInfoResponse{AssetInfos: assetInfos}, nil
}

func (k Keeper) QueStakerSpecifiedAssetAmount(ctx context.Context, req *restakingtype.QuerySpecifiedAssetAmountReq) (*restakingtype.StakerSingleAssetOrChangeInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakerSpecifiedAssetInfo(c, req.StakerId, req.AssetId)
}

func (k Keeper) QueOperatorAssetInfos(ctx context.Context, infos *restakingtype.QueryOperatorAssetInfos) (*restakingtype.QueryOperatorAssetInfosResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(infos.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	assetInfos, err := k.GetOperatorAssetInfos(c, addr)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryOperatorAssetInfosResponse{AssetInfos: assetInfos}, nil
}

func (k Keeper) QueOperatorSpecifiedAssetAmount(ctx context.Context, req *restakingtype.QueryOperatorSpecifiedAssetAmountReq) (*restakingtype.OperatorSingleAssetOrChangeInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(req.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return k.GetOperatorSpecifiedAssetInfo(c, addr, req.AssetId)
}

func (k Keeper) QueStakerExoCoreAddr(ctx context.Context, req *restakingtype.QueryStakerExCoreAddr) (*restakingtype.QueryStakerExCoreAddrResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	exoCoreAddr, err := k.GetStakerExoCoreAddr(c, req.StakerId)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryStakerExCoreAddrResponse{ExCoreAddr: exoCoreAddr}, nil
}
