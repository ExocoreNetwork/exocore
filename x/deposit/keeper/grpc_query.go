package keeper

import (
	"context"
	"github.com/exocore/x/deposit/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) GetStakerExoCoreAddr(ctx context.Context, addr *types.QueryStakerExCoreAddr) (*types.QueryStakerExCoreAddrResponse, error) {
	//TODO implement me
	panic("implement me")
}
