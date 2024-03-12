package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) IndexRecentParams(goCtx context.Context, req *types.QueryGetIndexRecentParamsRequest) (*types.QueryGetIndexRecentParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetIndexRecentParams(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetIndexRecentParamsResponse{IndexRecentParams: val}, nil
}
