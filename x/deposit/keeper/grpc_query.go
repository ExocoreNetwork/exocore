package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/deposit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	params, err := k.GetParams(c)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}
