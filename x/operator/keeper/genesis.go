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
	for _, record := range state.OperatorRecords {
		addr := record.OperatorAddress
		// #nosec G703 // already validated
		operatorAddr, _ := sdk.AccAddressFromBech32(addr)
		for _, detail := range record.Chains {
			chainID := detail.ChainID
			// validate that the chain exists
			// TODO: move this check to the avs keeper when implemented.
			if chainID != ctx.ChainID() {
				panic("unknown chain id")
			}
			// opt into the specified chain (TODO: avs address format)
			if err := k.OptIn(ctx, operatorAddr, chainID); err != nil {
				panic(err)
			}
			// #nosec G703 // already validated
			key, _ := types.HexStringToPubKey(detail.ConsensusKey)
			// then set pub key
			if err := k.setOperatorConsKeyForChainID(
				ctx, operatorAddr, chainID, key, true,
			); err != nil {
				panic(err)
			}
		}
	}
	return []abci.ValidatorUpdate{}
}

func (Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	// TODO
	return types.DefaultGenesis()
}
