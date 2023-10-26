package keeper

import (
	"context"

    "github.com/exocore/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) ClaimRewardRequest(goCtx context.Context,  msg *types.MsgClaimRewardRequest) (*types.MsgClaimRewardRequestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

	return &types.MsgClaimRewardRequestResponse{}, nil
}
