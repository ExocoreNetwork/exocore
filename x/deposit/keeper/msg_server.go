package keeper

import (
	context "context"
	"github.com/exocore/x/deposit/types"
)

var _ types.MsgServer = &Keeper{}

func (k Keeper) SetStakerExoCoreAddr(ctx context.Context, addr *types.MsgSetExoCoreAddr) (*types.MsgSetExoCoreAddrResponse, error) {
	//TODO implement me
	panic("implement me")
}
