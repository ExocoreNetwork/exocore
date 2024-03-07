package oracle

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the prices
	for _, elem := range genState.PricesList {
		k.SetPrices(ctx, elem)
	}
	// Set if defined
	if genState.Validators != nil {
		k.SetValidators(ctx, *genState.Validators)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PricesList = k.GetAllPrices(ctx)
	// Get all validators
	validators, found := k.GetValidators(ctx)
	if found {
		genesis.Validators = &validators
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
