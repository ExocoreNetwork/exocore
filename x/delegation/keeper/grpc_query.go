package keeper

import (
	"context"
	types2 "github.com/exocore/x/delegation/types"
)

var _ types2.QueryServer = Keeper{}

func (k Keeper) GetDelegationInfo(ctx context.Context, info *types2.QueryDelegationInfo) (*types2.QueryDelegationInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}
