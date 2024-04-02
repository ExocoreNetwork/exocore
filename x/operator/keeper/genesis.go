package keeper

import (
	"cosmossdk.io/math"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
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
	stakerAssetOperatorMap := make(map[string]map[string]map[string]math.Int)
	for _, level1 := range state.StakerRecords {
		stakerID := level1.StakerID
		if _, ok := stakerAssetOperatorMap[stakerID]; !ok {
			stakerAssetOperatorMap[stakerID] = make(map[string]map[string]math.Int)
		}
		for _, level2 := range level1.StakerDetails {
			assetID := level2.AssetID
			if _, ok := stakerAssetOperatorMap[stakerID][assetID]; !ok {
				stakerAssetOperatorMap[stakerID][assetID] = make(map[string]math.Int)
			}
			for _, level3 := range level2.Details {
				operatorAddress := level3.OperatorAddress
				amount := level3.Amount
				if err := k.UpdateOptedInAssetsState(
					ctx, stakerID, assetID, operatorAddress, amount,
				); err != nil {
					panic(err)
				}
				if _, ok := stakerAssetOperatorMap[stakerID][assetID][operatorAddress]; !ok {
					stakerAssetOperatorMap[stakerID][assetID][operatorAddress] = math.ZeroInt()
				}
				stakerAssetOperatorMap[stakerID][assetID][operatorAddress].Add(amount)
			}
		}
	}
	// validate the information in the delegation keeper,
	// which has validated it in the assets keeper.
	checkFunc := func(
		stakerID, assetID, operatorAddress string, state *delegationtypes.DelegationAmounts,
	) error {
		valueHere := stakerAssetOperatorMap[stakerID][assetID][operatorAddress]
		if !valueHere.Equal(state.UndelegatableAmount) {
			return types.ErrInvalidGenesisData
		}
		return nil
	}
	// since this module only knows the delegated value (and not the deposit value),
	// it cannot do any further validation with the data in the assets keeper.
	if err := k.delegationKeeper.IterateDelegationState(ctx, checkFunc); err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}

func (Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	return types.DefaultGenesis()
}
