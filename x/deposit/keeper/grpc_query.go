package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/deposit/types"
)

func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	params, err := k.GetParams(c)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}
