package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis loads the initial state from a genesis file.

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) []abci.ValidatorUpdate {

	//code := []byte("chain-id-code")        // declare as a constant in x/avs/types, here is just for illustration
	//codeHash := crypto.Keccak256Hash(code) // same as above
	//k.evmKeeper.SetCode(ctx, codeHash.Bytes(), code)
	return []abci.ValidatorUpdate{}
}

func (Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	// TODO
	return types.DefaultGenesis()
}
