package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	epochID := genState.Params.EpochIdentifier
	_, found := k.epochsKeeper.GetEpochInfo(ctx, epochID)
	if !found {
		// the panic is suitable here because it is being done at genesis, when the node
		// is not running. it means that the genesis file is malformed.
		panic("not found the epoch info")
		return
	}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	return genesis
}
