package keeper

import (
	"context"

    "github.com/exocore/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) ClaimRewardResponse(goCtx context.Context,  msg *types.MsgClaimRewardResponse) (*types.MsgClaimRewardResponseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

	return &types.MsgClaimRewardResponseResponse{}, nil
}
