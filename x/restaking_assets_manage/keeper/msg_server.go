package keeper

import (
	"context"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

var _ types2.MsgServer = &Keeper{}

func (k Keeper) SetStakerExoCoreAddr(ctx context.Context, addrInfo *types2.MsgSetExoCoreAddr) (*types2.MsgSetExoCoreAddrResponse, error) {
	//TODO implement me
	panic("implement me")
}
