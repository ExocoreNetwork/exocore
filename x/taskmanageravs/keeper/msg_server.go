package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"github.com/ExocoreNetwork/exocore/x/taskmanageravs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

func (m msgServer) RegisterAVSTask(ctx context.Context, req *types.RegisterAVSTaskReq) (*types.RegisterAVSTaskResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	avs := m.avsKeeper.IsAVS(c, sdk.AccAddress(req.AVSAddress))
	if !avs {
		return nil, errorsmod.Wrap(types.ErrNotYetRegistered, fmt.Sprintf("RegisterAVSTask: avs address is %s", req.GetAVSAddress()))

	}
	_, err := m.Keeper.SetAvsTaskInfo(c, req)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
