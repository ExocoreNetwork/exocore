package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/exocore/x/reward/types"
)

var _ types.QueryServiceServer = Querier{}

// Querier implements the grpc queries for the nexus module
type Querier struct {
	keeper Keeper
	minter types.Minter
	nexus  types.Exocore
}

// NewGRPCQuerier creates a new nexus Querier
func NewGRPCQuerier(k Keeper, m types.Minter, n types.Exocore) Querier {
	return Querier{
		keeper: k,
		minter: m,
		nexus:  n,
	}
}

// Params returns the reward module params
func (q Querier) Params(c context.Context, req *types.ParamsRequest) (*types.ParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	params := q.keeper.GetParams(ctx)

	return &types.ParamsResponse{
		Params: &params,
	}, nil
}
