package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/reward/types"
)

func (k msgServer) ClaimRewardRequest(goCtx context.Context, msg *types.MsgClaimRewardRequest) (*types.MsgClaimRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgClaimRewardResponse{}, nil
}
