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
	assetInfo, err := k.GetStakerAssetInfos(c, info.StakerId)
	if err != nil {
		return nil, err
	}
	response := &types2.QueryAssetInfoResponse{AssetAmounts: make(map[string]*types2.QueryAssetInfoResponse_ValueField, 0)}
	for k, v := range assetInfo {
		response.AssetAmounts[k] = &types2.QueryAssetInfoResponse_ValueField{
			Amount: v,
		}
	}
	return response, nil
}

func (k Keeper) QueStakerSpecifiedAssetAmount(ctx context.Context, req *types2.QuerySpecifiedAssetAmountReq) (*types2.QuerySpecifiedAssetAmountReqResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	amount, err := k.GetStakerSpecifiedAssetAmount(c, req.StakerId, req.AssetId)
	if err != nil {
		return nil, err
	}
	return &types2.QuerySpecifiedAssetAmountReqResponse{
		Amount: amount,
	}, nil
}

func (k Keeper) QueOperatorAssetInfos(ctx context.Context, infos *types2.QueryOperatorAssetInfos) (*types2.QueryAssetInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(infos.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	assetInfo, err := k.GetOperatorAssetInfos(c, addr)
	if err != nil {
		return nil, err
	}
	response := &types2.QueryAssetInfoResponse{AssetAmounts: make(map[string]*types2.QueryAssetInfoResponse_ValueField, 0)}
	for k, v := range assetInfo {
		response.AssetAmounts[k] = &types2.QueryAssetInfoResponse_ValueField{
			Amount: v,
		}
	}
	return response, nil
}

func (k Keeper) QueOperatorSpecifiedAssetAmount(ctx context.Context, req *types2.QueryOperatorSpecifiedAssetAmountReq) (*types2.QuerySpecifiedAssetAmountReqResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(req.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	amount, err := k.GetOperatorSpecifiedAssetAmount(c, addr, req.AssetId)
	if err != nil {
		return nil, err
	}
	return &types2.QuerySpecifiedAssetAmountReqResponse{
		Amount: amount,
	}, nil
}

func (k Keeper) QueStakerExoCoreAddr(ctx context.Context, req *types2.QueryStakerExCoreAddr) (*types2.QueryStakerExCoreAddrResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	exoCoreAddr, err := k.GetStakerExoCoreAddr(c, req.StakerId)
	if err != nil {
		return nil, err
	}
	return &types2.QueryStakerExCoreAddrResponse{ExCoreAddr: exoCoreAddr}, nil
}
