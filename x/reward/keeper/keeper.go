package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ExocoreNetwork/exocore/x/assets/keeper"
	"github.com/ExocoreNetwork/exocore/x/reward/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	// other keepers
	assetsKeeper keeper.Keeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	assetsKeeper keeper.Keeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		assetsKeeper: assetsKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
