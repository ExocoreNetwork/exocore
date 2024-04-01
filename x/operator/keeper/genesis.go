package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, state types.GenesisState) []abci.ValidatorUpdate {
	// operators.go
	for _, infoCopy := range state.Operators {
		info := infoCopy // prevent implicit memory aliasing
		if err := k.SetOperatorInfo(ctx, info.EarningsAddr, &info); err != nil {
			panic(err)
		}
		operatorAddress := info.EarningsAddr
		// #nosec G703 // already validated
		operatorAccAddress, _ := sdk.AccAddressFromBech32(operatorAddress)
		if err := k.OptIn(ctx, operatorAccAddress, ctx.ChainID()); err != nil {
			panic(err)
		}
	}
	// consensus_keys.go
	for _, record := range state.OperatorRecords {
		operatorAddress := record.OperatorAddress
		// #nosec G703 // already validated
		operatorAccAddress, _ := sdk.AccAddressFromBech32(operatorAddress)
		for _, subRecord := range record.Chains {
			consKeyBytes32 := subRecord.ConsensusKey
			// #nosec G703 // already validated
			consKey, _ := types.HexStringToPubKey(consKeyBytes32)
			if err := k.SetOperatorConsKeyForChainID(
				ctx, operatorAccAddress, subRecord.ChainID, consKey,
			); err != nil {
				panic(err)
			}
		}
	}
	// state_update.go
	for _, level1 := range state.StakerRecords {
		stakerID := level1.StakerID
		for _, level2 := range level1.StakerDetails {
			assetID := level2.AssetID
			for _, level3 := range level2.Details {
				operatorAddress := level3.OperatorAddress
				amount := level3.Amount
				if err := k.UpdateOptedInAssetsState(
					ctx, stakerID, assetID, operatorAddress, amount,
				); err != nil {
					panic(err)
				}

			}
		}
	}
	return []abci.ValidatorUpdate{}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return types.DefaultGenesis()
}
