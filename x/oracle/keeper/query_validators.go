package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Validators(goCtx context.Context, req *types.QueryGetValidatorsRequest) (*types.QueryGetValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetValidators(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetValidatorsResponse{Validators: val}, nil
}
