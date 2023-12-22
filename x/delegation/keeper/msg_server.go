package keeper

import (
	context "context"

	errorsmod "cosmossdk.io/errors"
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

// DelegateAssetToOperator todo: Delegation and Undelegation from exoCore chain directly will be implemented in future.At the moment,they are executed from client chain
func (k Keeper) DelegateAssetToOperator(ctx context.Context, delegation *types.MsgDelegation) (*types.DelegationResponse, error) {
	return nil, errorsmod.Wrap(types.ErrNotSupportYet, "func:DelegateAssetToOperator")
}

func (k Keeper) UndelegateAssetFromOperator(ctx context.Context, delegation *types.MsgUndelegation) (*types.UndelegationResponse, error) {
	return nil, errorsmod.Wrap(types.ErrNotSupportYet, "func:UndelegateAssetFromOperator")
}
