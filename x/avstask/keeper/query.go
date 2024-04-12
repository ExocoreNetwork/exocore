package keeper

import (
	"context"

	avstasktypes "github.com/ExocoreNetwork/exocore/x/avstask/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ avstasktypes.QueryServer = &Keeper{}

func (k *Keeper) QueryAVSTaskInfo(ctx context.Context, req *avstasktypes.GetAVSTaskInfoReq) (*avstasktypes.TaskContractInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetAVSTaskInfo(c, req.TaskAddr)
}
