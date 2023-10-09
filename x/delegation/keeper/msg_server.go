package keeper

import (
	context "context"
	"github.com/exocore/x/delegation/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) RegisterOperator(ctx context.Context, info *types.OperatorInfo) (*types.RegisterOperatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) DelegateAssetToOperator(ctx context.Context, delegation *types.MsgDelegation) (*types.DelegationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) UnDelegateAssetFromOperator(ctx context.Context, delegation *types.MsgUnDelegation) (*types.UnDelegationResponse, error) {
	//TODO implement me
	panic("implement me")
}
