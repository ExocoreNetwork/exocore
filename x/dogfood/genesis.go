package dogfood

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/keeper"
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	genState types.GenesisState,
) []abci.ValidatorUpdate {
	k.SetParams(ctx, genState.Params)
	return k.ApplyValidatorChanges(ctx, genState.ValSet)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetDogfoodParams(ctx)

	return genesis
}
