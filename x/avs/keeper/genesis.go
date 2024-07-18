package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// InitGenesis initializes the module's state from a provided genesis state.
// Since this action typically occurs on chain starts, this function is allowed to panic.
func (k Keeper) InitGenesis(
	ctx sdk.Context,
	_ types.GenesisState,
) []abci.ValidatorUpdate {
	code := []byte(types.ChainID)
	codeHash := crypto.Keccak256Hash(code)
	k.evmKeeper.SetCode(ctx, codeHash.Bytes(), code)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the module's exported genesis
func (Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	// TODO
	return types.DefaultGenesis()
}
