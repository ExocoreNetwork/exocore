package keeper

import (
	"context"

	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueClientChainInfoByIndex query client chain info by clientChainLzID
func (k Keeper) QueClientChainInfoByIndex(ctx context.Context, info *restakingtype.QueryClientChainInfo) (*restakingtype.ClientChainInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetClientChainInfoByIndex(c, info.ChainIndex)
}

// QueAllClientChainInfo query all client chain info that have been registered in exoCore
// the key of returned map is clientChainLzID, the value is the client chain info.
func (k Keeper) QueAllClientChainInfo(ctx context.Context, _ *restakingtype.QueryAllClientChainInfo) (*restakingtype.QueryAllClientChainInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllClientChainInfo(c)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryAllClientChainInfoResponse{AllClientChainInfos: allInfo}, nil
}

// QueStakingAssetInfo query the specified client chain asset info by inputting assetID
func (k Keeper) QueStakingAssetInfo(ctx context.Context, info *restakingtype.QueryStakingAssetInfo) (*restakingtype.StakingAssetInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakingAssetInfo(c, info.AssetID)
}

// QueAllStakingAssetsInfo query the info about all client chain assets that have been registered
func (k Keeper) QueAllStakingAssetsInfo(ctx context.Context, _ *restakingtype.QueryAllStakingAssetsInfo) (*restakingtype.QueryAllStakingAssetsInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	allInfo, err := k.GetAllStakingAssetsInfo(c)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryAllStakingAssetsInfoResponse{AllStakingAssetsInfo: allInfo}, nil
}

// QueStakerAssetInfos query th state of all assets for a staker specified by stakerID
func (k Keeper) QueStakerAssetInfos(ctx context.Context, info *restakingtype.QueryStakerAssetInfo) (*restakingtype.QueryAssetInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	assetInfos, err := k.GetStakerAssetInfos(c, info.StakerID)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryAssetInfoResponse{AssetInfos: assetInfos}, nil
}

// QueStakerSpecifiedAssetAmount query the specified asset state of a staker, using stakerID and assetID as query parameters
func (k Keeper) QueStakerSpecifiedAssetAmount(ctx context.Context, req *restakingtype.QuerySpecifiedAssetAmountReq) (*restakingtype.StakerSingleAssetOrChangeInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetStakerSpecifiedAssetInfo(c, req.StakerID, req.AssetID)
}

// QueOperatorAssetInfos query th state of all assets for an operator specified by operator address
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

// QueOperatorSpecifiedAssetAmount query the specified asset state of an operator, using operator address and assetID as query parameters
func (k Keeper) QueOperatorSpecifiedAssetAmount(ctx context.Context, req *restakingtype.QueryOperatorSpecifiedAssetAmountReq) (*restakingtype.OperatorSingleAssetOrChangeInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(req.OperatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return k.GetOperatorSpecifiedAssetInfo(c, addr, req.AssetID)
}

// QueStakerExoCoreAddr outdated,will be deprecated
func (k Keeper) QueStakerExoCoreAddr(ctx context.Context, req *restakingtype.QueryStakerExCoreAddr) (*restakingtype.QueryStakerExCoreAddrResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	exoCoreAddr, err := k.GetStakerExoCoreAddr(c, req.StakerID)
	if err != nil {
		return nil, err
	}
	return &restakingtype.QueryStakerExCoreAddrResponse{ExCoreAddr: exoCoreAddr}, nil
}
