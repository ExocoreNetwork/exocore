package keeper

import (
	"fmt"

	restakingkeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	depositkeeper "github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	"github.com/ExocoreNetwork/exocore/x/withdraw/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey

		// restaking keepers for asset status update
		assetsKeeper  restakingkeeper.Keeper
		depositKeeper depositkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	assetsKeeper restakingkeeper.Keeper,
	depositKeeper depositkeeper.Keeper,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		assetsKeeper:  assetsKeeper,
		depositKeeper: depositKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
