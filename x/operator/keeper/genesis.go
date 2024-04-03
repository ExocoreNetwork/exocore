package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	for _, op := range state.Operators {
		op := op
		if err := k.SetOperatorInfo(ctx, op.EarningsAddr, &op); err != nil {
			panic(err)
		}
	}
	for _, record := range state.OperatorRecords {
		addr := record.OperatorAddress
		// #nosec G703 // already validated
		operatorAddr, _ := sdk.AccAddressFromBech32(addr)
		bootstrapping := false
		for _, detail := range record.Chains {
			// opt into the specified chain (TODO: avs address format)
			if err := k.OptIn(ctx, operatorAddr, detail.ChainID); err != nil {
				panic(err)
			}
			// #nosec G703 // already validated
			key, _ := types.HexStringToPubKey(detail.ConsensusKey)
			// then set pub key
			if err := k.SetOperatorConsKeyForChainID(
				ctx, operatorAddr, detail.ChainID, key,
			); err != nil {
				panic(err)
			}
			bootstrapping = bootstrapping || ctx.ChainID() == detail.ChainID
		}
		if !bootstrapping {
			// TODO: consider removing this check
			panic("registered an operator but they aren't bootstrapping the current chain")
		}
	}
	return []abci.ValidatorUpdate{}
}

func (Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	// TODO
	return types.DefaultGenesis()
}
