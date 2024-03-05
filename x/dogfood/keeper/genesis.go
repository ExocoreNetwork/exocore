package keeper

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(
	ctx sdk.Context,
	genState types.GenesisState,
) []abci.ValidatorUpdate {
	k.SetParams(ctx, genState.Params)
	// the `params` validator is not super useful to validate state level information
	// so, it must be done here. by extension, the `InitGenesis` of the epochs module
	// should be called before that of this module.
	_, found := k.epochsKeeper.GetEpochInfo(ctx, genState.Params.EpochIdentifier)
	if !found {
		// the panic is suitable here because it is being done at genesis, when the node
		// is not running. it means that the genesis file is malformed.
		panic("epoch info not found")
	}
	// apply the same logic to the staking assets.
	for _, assetID := range genState.Params.AssetIDs {
		if !k.restakingKeeper.IsStakingAsset(ctx, assetID) {
			panic(fmt.Errorf("staking asset %s not found", assetID))
		}
	}
	// genState must not be malformed.
	if len(genState.ValSet) > int(k.GetMaxValidators(ctx)) { // #nosec G701 // ok if 64-bit.
		panic(fmt.Errorf(
			"cannot have more than %d validators in the genesis state",
			k.GetMaxValidators(ctx),
		))
	}
	return k.ApplyValidatorChanges(
		ctx, genState.ValSet, types.InitialValidatorSetID,
	)
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetDogfoodParams(ctx)
	// TODO(mm)
	return genesis
}
