package keeper

import (
	context "context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/delegation/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) RegisterOperator(ctx context.Context, req *types.RegisterOperatorReq) (*types.RegisterOperatorResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetOperatorInfo(c, req.FromAddress, req.Info)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// DelegateAssetToOperator todo: Delegation and unDelegation from exoCore chain directly will be implemented in future.At 他the moment,they are executed from client chain
func (k Keeper) DelegateAssetToOperator(ctx context.Context, delegation *types.MsgDelegation) (*types.DelegationResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) UnDelegateAssetFromOperator(ctx context.Context, delegation *types.MsgUnDelegation) (*types.UnDelegationResponse, error) {
	//TODO implement me
	panic("implement me")
}
