package keeper

import (
	"context"

    "github.com/exocore/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) RewardDetail(goCtx context.Context,  msg *types.MsgRewardDetail) (*types.MsgRewardDetailResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Handling the message
    _ = ctx

	return &types.MsgRewardDetailResponse{}, nil
}
