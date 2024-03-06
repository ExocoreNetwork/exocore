package keeper

import (
	context "context"

	errorsmod "cosmossdk.io/errors"
	"github.com/ExocoreNetwork/exocore/x/delegation/types"
)

var _ types.MsgServer = &Keeper{}

// DelegateAssetToOperator todo: Delegation and Undelegation from exoCore chain directly will be implemented in future.At the moment,they are executed from client chain
func (k *Keeper) DelegateAssetToOperator(_ context.Context, _ *types.MsgDelegation) (*types.DelegationResponse, error) {
	return nil, errorsmod.Wrap(types.ErrNotSupportYet, "func:DelegateAssetToOperator")
}

func (k *Keeper) UndelegateAssetFromOperator(_ context.Context, _ *types.MsgUndelegation) (*types.UndelegationResponse, error) {
	return nil, errorsmod.Wrap(types.ErrNotSupportYet, "func:UndelegateAssetFromOperator")
}
