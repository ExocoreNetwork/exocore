package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	for i := range state.Operators {
		op := state.Operators[i] // avoid implicit memory aliasing
		if err := k.SetOperatorInfo(ctx, op.EarningsAddr, &op); err != nil {
			panic(err)
		}
	}
	return []abci.ValidatorUpdate{}
}

func (Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	// TODO
	return types.DefaultGenesis()
}
