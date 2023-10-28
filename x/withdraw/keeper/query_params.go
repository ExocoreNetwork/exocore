package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/withdraw/types"
)

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	c := sdk.UnwrapSDKContext(goCtx)
	params, err := k.GetParams(c)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{
		Params: params,
	}, nil
}
