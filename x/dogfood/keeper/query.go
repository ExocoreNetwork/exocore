package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(
	goCtx context.Context,
	req *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetDogfoodParams(ctx)}, nil
}
