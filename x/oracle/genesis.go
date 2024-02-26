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
	// Set all the roundInfo
	for _, elem := range genState.RoundInfoList {
		k.SetRoundInfo(ctx, elem)
	}
	// Set all the roundData
	for _, elem := range genState.RoundDataList {
		k.SetRoundData(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.PricesList = k.GetAllPrices(ctx)
	genesis.RoundInfoList = k.GetAllRoundInfo(ctx)
	genesis.RoundDataList = k.GetAllRoundData(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
