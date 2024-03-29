package keeper

import (
	"context"
	"fmt"
	"github.com/ExocoreNetwork/exocore/x/taskmanageravs/types"
	tasktype "github.com/ExocoreNetwork/exocore/x/taskmanageravs/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc       codec.BinaryCodec
		storeKey  storetypes.StoreKey
		avsKeeper tasktype.AvsKeeper
	}
)

func (k Keeper) RegisterAVSTask(ctx context.Context, req *types.RegisterAVSTaskReq) (*types.RegisterAVSTaskResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	avsKeeper tasktype.AvsKeeper,
) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  storeKey,
		avsKeeper: avsKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

type ITask interface {
	RegisterAVSTask(ctx context.Context, req *types.RegisterAVSTaskReq) (*types.RegisterAVSTaskResponse, error)
}
