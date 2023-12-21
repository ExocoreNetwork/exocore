package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	depositkeeper "github.com/exocore/x/deposit/keeper"
	restakingkeeper "github.com/exocore/x/restaking_assets_manage/keeper"
	"github.com/exocore/x/withdraw/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey

		// restaking keepers for asset status update
		restakingStateKeeper restakingkeeper.Keeper
		depositKeeper        depositkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	restakingStateKeeper restakingkeeper.Keeper,
	depositKeeper depositkeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:                  cdc,
		storeKey:             storeKey,
		restakingStateKeeper: restakingStateKeeper,
		depositKeeper:        depositKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
