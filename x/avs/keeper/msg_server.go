package keeper

import (
	"context"
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) RegisterAVS(ctx context.Context, req *types.RegisterAVSReq) (*types.RegisterAVSResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	fromAddress := req.FromAddress
	operatorAddress := req.Info.OperatorAddress
	for _, opAddr := range operatorAddress {
		if fromAddress == opAddr {
			// Set purely for AVS itself information.
			if err := k.SetAVSInfo(c, req.Info); err != nil {
				return nil, err
			}
		}
	}
	return nil, fmt.Errorf("The fromAddress %s is different from operatorAddress %s ", fromAddress, operatorAddress)
}

func (k Keeper) DeRegisterAVS(ctx context.Context, req *types.DeRegisterAVSReq) (*types.DeRegisterAVSResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	if err := k.DeleteAVSInfo(c, req.Info); err != nil {
		return nil, err
	}

	return nil, nil
}
