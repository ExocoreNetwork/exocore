package keeper

import (
	"context"
	"errors"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) StakerInfos(goCtx context.Context, req *types.QueryStakerInfosRequest) (*types.QueryStakerInfosResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request for stakerInfos")
	}
	if !assetstypes.IsNativeToken(req.AssetID) {
		return nil, errors.New("assetID doesn't reprensents any supported nativeRestakingToken")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	stakerInfos := k.GetStakerInfos(ctx, req.AssetID)
	return &types.QueryStakerInfosResponse{StakerInfos: stakerInfos}, nil
}

func (k Keeper) StakerList(goCtx context.Context, req *types.QueryStakerListRequest) (*types.QueryStakerListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request for stakerList")
	}
	if !assetstypes.IsNativeToken(req.AssetID) {
		return nil, errors.New("assetID doesn't reprensents any supported nativeRestakingToken")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	stakerList := k.GetStakerList(ctx, req.AssetID)
	return &types.QueryStakerListResponse{StakerList: &stakerList}, nil
}
