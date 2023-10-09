package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueClientChainInfoByIndex(ctx context.Context, info *types2.QueryClientChainInfo) (*types2.ClientChainInfo, error) {
	return k.GetClientChainInfoByIndex(info.ChainIndex)
}

func (k Keeper) QueAllClientChainInfo(ctx context.Context, info *types2.QueryAllClientChainInfo) (*types2.QueryAllClientChainInfoResponse, error) {
	allInfo, err := k.GetAllClientChainInfo()
	if err != nil {
		return nil, err
	}
	return &types2.QueryAllClientChainInfoResponse{AllClientChainInfos: allInfo}, nil
}

func (k Keeper) QueStakingAssetInfo(ctx context.Context, info *types2.QueryStakingAssetInfo) (*types2.StakingAssetInfo, error) {
	return k.GetStakingAssetInfo(info.AssetId)
}

func (k Keeper) QueAllStakingAssetsInfo(ctx context.Context, info *types2.QueryAllStakingAssetsInfo) (*types2.QueryAllStakingAssetsInfoResponse, error) {
	allInfo, err := k.GetAllStakingAssetsInfo()
	if err != nil {
		return nil, err
	}
	return &types2.QueryAllStakingAssetsInfoResponse{AllStakingAssetsInfo: allInfo}, nil
}

func (k Keeper) QueStakerAssetInfos(ctx context.Context, info *types2.QueryStakerAssetInfo) (*types2.QueryAssetInfoResponse, error) {
	assetInfo, err := k.GetStakerAssetInfos(info.StakerId)
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
	amount, err := k.GetStakerSpecifiedAssetAmount(req.StakerId, req.AssetId)
	if err != nil {
		return nil, err
	}
	return &types2.QuerySpecifiedAssetAmountReqResponse{
		Amount: amount,
	}, nil
}

func (k Keeper) QueOperatorAssetInfos(ctx context.Context, infos *types2.QueryOperatorAssetInfos) (*types2.QueryAssetInfoResponse, error) {
	addr, err := sdk.AccAddressFromBech32(infos.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	assetInfo, err := k.GetOperatorAssetInfos(addr)
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
	addr, err := sdk.AccAddressFromBech32(req.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	amount, err := k.GetOperatorSpecifiedAssetAmount(addr, req.AssetId)
	if err != nil {
		return nil, err
	}
	return &types2.QuerySpecifiedAssetAmountReqResponse{
		Amount: amount,
	}, nil
}

func (k Keeper) QueStakerExoCoreAddr(ctx context.Context, req *types2.QueryStakerExCoreAddr) (*types2.QueryStakerExCoreAddrResponse, error) {
	exoCoreAddr, err := k.GetStakerExoCoreAddr(ctx, req.StakerId)
	if err != nil {
		return nil, err
	}
	return &types2.QueryStakerExCoreAddrResponse{ExCoreAddr: exoCoreAddr}, nil
}
