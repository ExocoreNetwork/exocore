package keeper

import (
	"context"

    "github.com/exocore/x/withdraw/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) Withdraw(goCtx context.Context,  msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

	return &types.MsgWithdrawResponse{}, nil
}
