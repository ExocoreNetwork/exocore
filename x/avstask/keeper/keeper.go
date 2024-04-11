package keeper

import (
	"context"
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/avs/keeper"
	"github.com/ExocoreNetwork/exocore/x/avstask/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc       codec.BinaryCodec
		storeKey  storetypes.StoreKey
		avsKeeper keeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	avsKeeper keeper.Keeper,
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

type TaskKeeper interface {
	RegisterAVSTask(ctx context.Context, req *types.RegisterAVSTaskReq) (*types.RegisterAVSTaskResponse, error)
}
