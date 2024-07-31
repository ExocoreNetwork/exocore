package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = &Keeper{}

func (k Keeper) QueryAVSInfo(ctx context.Context, req *types.QueryAVSInfoReq) (*types.QueryAVSInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetAVSInfo(c, req.AVSAddress)
}

func (k Keeper) QueryAVSTaskInfo(ctx context.Context, req *types.QueryAVSTaskInfoReq) (*types.TaskInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetTaskInfo(c, req.TaskId, req.TaskAddr)
}

// QueryAVSAddrByChainID is an implementation of the QueryAVSAddrByChainID gRPC method
func (k Keeper) QueryAVSAddrByChainID(ctx context.Context, req *types.QueryAVSAddrByChainIDReq) (*types.QueryAVSAddrByChainIDResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	isChainAvs, avsAddr := k.IsAVSByChainID(c, types.ChainIDWithoutRevision(req.ChainID))
	if !isChainAvs {
		return nil, types.ErrNotYetRegistered
	}
	return &types.QueryAVSAddrByChainIDResponse{AVSAddress: avsAddr.String()}, nil
}
