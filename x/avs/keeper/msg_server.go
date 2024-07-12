package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) RegisterAVS(_ context.Context, _ *types.RegisterAVSReq) (*types.RegisterAVSResponse, error) {
	// Disable cosmos transaction temporarily
	// c := sdk.UnwrapSDKContext(ctx)
	// fromAddress := req.FromAddress
	// operatorAddress := req.Info.OperatorAddress
	// for _, opAddr := range operatorAddress {
	// 	if fromAddress == opAddr {
	// 		// Set purely for AVS itself information.
	// 		if err := k.SetAVSInfo(c, req.Info); err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }
	return nil, nil
}

func (k Keeper) DeRegisterAVS(_ context.Context, _ *types.DeRegisterAVSReq) (*types.DeRegisterAVSResponse, error) {
	// Disable cosmos transaction temporarily
	// c := sdk.UnwrapSDKContext(ctx)
	// if err := k.DeleteAVSInfo(c, req.Info); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

func (k Keeper) RegisterAVSTask(ctx context.Context, req *types.RegisterAVSTaskReq) (*types.RegisterAVSTaskResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	avs, err := k.GetAVSInfo(c, req.FromAddress)
	if err != nil {
		return nil, err
	}
	if avs.GetInfo() == nil {
		return nil, errorsmod.Wrap(types.ErrNotYetRegistered, fmt.Sprintf("RegisterAVSTask: avs address is %s", req.FromAddress))
	}
	err = k.SetAVSTaskInfo(c, req)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
