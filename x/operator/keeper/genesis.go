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
			// checks for cosmos chain as AVS registration.
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
			// there is no need to check whether the asset is registered or not,
			// since an unregistered asset being in the map will fail in the comparison
			// against the delegation module (which checks that assets are registered).
			if _, ok := stakerAssetOperatorMap[stakerID][assetID]; !ok {
				stakerAssetOperatorMap[stakerID][assetID] = make(map[string]math.Int)
			}
			for _, level3 := range level2.Details {
				operatorAddress := level3.OperatorAddress
				// check that the operator is registered. this is necessary to do so here, since
				// the delegation module has not checked it. if this check did not exist, it
				// would be possible to creeate a delegation (according to both modules) for
				// an operator that is not registered.
				// #nosec G703 // already validated
				operatorAddressAcc, _ := sdk.AccAddressFromBech32(operatorAddress)
				if !k.IsOperator(ctx, operatorAddressAcc) {
					panic("operator not found")
				}
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
	// D for delegation
	stakerAssetOperatorMapD := make(map[string]map[string]map[string]math.Int)
	// validate the information in the delegation keeper,
	// which has validated it in the assets keeper.
	// the validation in the delegation keeper is that
	// the delegated amount is less than the deposited amount (before the delegation).
	// so, in summary, our checks are that
	// deposit[staker][asset] >= sum_over_operators(delegated[staker][asset][operator]), and
	// x_delegation[staker][asset][operator] == x_operator[staker][asset][operator].
	checkFunc := func(
		stakerID, assetID, operatorAddress string, state *delegationtypes.DelegationAmounts,
	) error {
		if _, ok := stakerAssetOperatorMapD[stakerID]; !ok {
			stakerAssetOperatorMapD[stakerID] = make(map[string]map[string]math.Int)
		}
		if _, ok := stakerAssetOperatorMapD[stakerID][assetID]; !ok {
			stakerAssetOperatorMapD[stakerID][assetID] = make(map[string]math.Int)
		}
		if _, ok := stakerAssetOperatorMapD[stakerID][assetID][operatorAddress]; !ok {
			stakerAssetOperatorMapD[stakerID][assetID][operatorAddress] = math.ZeroInt()
		}
		stakerAssetOperatorMapD[stakerID][assetID][operatorAddress].Add(state.UndelegatableAmount)
		return nil
	}
	if err := k.delegationKeeper.IterateDelegationState(ctx, checkFunc); err != nil {
		// should never happen
		panic(err)
	}
	if !isEqual(stakerAssetOperatorMap, stakerAssetOperatorMapD) {
		panic("delegation and operator module values are inconsistent")
	}
	// we have checked that the results in operator, delegation (implied assets) are consistent.
	// there is no need to check again for assets again here.
	return []abci.ValidatorUpdate{}
}

func (Keeper) ExportGenesis(sdk.Context) *types.GenesisState {
	return types.DefaultGenesis()
}

// isEqual compares two nested maps and returns true if they are equal, false otherwise.
func isEqual(map1, map2 map[string]map[string]map[string]math.Int) bool {
	// Check if map1 keys exist in map2 and values are equal
	for stakerID, assetMap := range map1 {
		if assetMap2, ok := map2[stakerID]; ok {
			for assetID, operatorMap := range assetMap {
				if operatorMap2, ok := assetMap2[assetID]; ok {
					for operatorAddress, amount := range operatorMap {
						if amount2, ok := operatorMap2[operatorAddress]; ok {
							if !amount.Equal(amount2) {
								return false // Amounts differ
							}
						} else {
							return false // OperatorAddress not found in map2
						}
					}
				} else {
					return false // AssetID not found in map2
				}
			}
		} else {
			return false // StakerID not found in map2
		}
	}

	// Check if map2 has extra keys not present in map1
	for stakerID, assetMap := range map2 {
		if _, ok := map1[stakerID]; !ok {
			return false // Extra StakerID in map2
		}
		for assetID := range assetMap {
			if _, ok := map1[stakerID][assetID]; !ok {
				return false // Extra AssetID in map2
			}
		}
	}

	return true
}
