package keeper

import (
	"context"
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgServerImpl struct {
	keeper Keeper
}

func NewMsgServerImpl(keeper Keeper) *MsgServerImpl {
	return &MsgServerImpl{keeper: keeper}
}

var _ types.MsgServer = &MsgServerImpl{}

func (m MsgServerImpl) SubmitTaskResult(goCtx context.Context, req *types.SubmitTaskResultReq) (*types.SubmitTaskResultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := m.keeper.SetTaskResultInfo(ctx, req.FromAddress, req.Info); err != nil {
		return nil, err
	}
	return &types.SubmitTaskResultResponse{}, nil
}

func (m MsgServerImpl) RegisterAVS(ctx context.Context, req *types.RegisterAVSReq) (*types.RegisterAVSResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m MsgServerImpl) DeRegisterAVS(ctx context.Context, req *types.DeRegisterAVSReq) (*types.DeRegisterAVSResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m MsgServerImpl) RegisterAVSTask(ctx context.Context, req *types.RegisterAVSTaskReq) (*types.RegisterAVSTaskResponse, error) {
	// TODO implement me
	panic("implement me")
}
