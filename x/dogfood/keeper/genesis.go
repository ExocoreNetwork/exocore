package keeper

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
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
	epochID := genState.Params.EpochIdentifier
	_, found := k.epochsKeeper.GetEpochInfo(ctx, epochID)
	if !found {
		// the panic is suitable here because it is being done at genesis, when the node
		// is not running. it means that the genesis file is malformed.
		panic(fmt.Sprintf("epoch info not found %s", epochID))
	}
	// apply the same logic to the staking assets.
	for _, assetID := range genState.Params.AssetIDs {
		if !k.restakingKeeper.IsStakingAsset(ctx, assetID) {
			panic(fmt.Errorf("staking param %s not found in assets module", assetID))
		}
	}
	// at genesis, not chain restarts, each operator may not necessarily be an initial
	// validator. this is because the operator may not have enough minimum self delegation
	// to be considered, or may not be in the top N operators. so checking that count here
	// is meaningless as well.
	totalPower := sdk.NewInt(0)
	out := make([]abci.ValidatorUpdate, len(genState.ValSet))
	for _, val := range genState.ValSet {
		// #nosec G703 // already validated
		consKey, _ := operatortypes.HexStringToPubKey(val.PublicKey)
		// #nosec G703 // this only fails if the key is of a type not already defined.
		consAddr, _ := operatortypes.TMCryptoPublicKeyToConsAddr(consKey)
		// if GetOperatorAddressForChainIDAndConsAddr returns found, it means
		// that the operator is registered and also (TODO) that it has opted into
		// the dogfood AVS.
		found, _ := k.operatorKeeper.GetOperatorAddressForChainIDAndConsAddr(
			ctx, ctx.ChainID(), consAddr,
		)
		if !found {
			panic(fmt.Sprintf("operator not found: %s", consAddr))
		}
		out = append(out, abci.ValidatorUpdate{
			PubKey: *consKey,
			Power:  val.Power,
		})
		totalPower = totalPower.Add(sdk.NewInt(val.Power))
	}
	k.SetLastTotalPower(ctx, totalPower)

	// ApplyValidatorChanges will sort it internally
	return k.ApplyValidatorChanges(
		ctx, out,
	)
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetDogfoodParams(ctx)
	// TODO(mm)
	return genesis
}
