package keeper

import (
	sdkmath "cosmossdk.io/math"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationkeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalculateUSDValueForOperator calculates the total and self usd value for the
// operator according to the input assets filter and prices.
func (k *Keeper) CalculateUSDValueForOperator(
	ctx sdk.Context,
	operator string,
	assetsFilter map[string]interface{},
	decimals map[string]uint32,
	prices map[string]operatortypes.Price,
) (sdkmath.LegacyDec, sdkmath.LegacyDec, error) {
	usdValue := sdkmath.LegacyNewDec(0)
	selfUSDValue := sdkmath.LegacyNewDec(0)
	// iterate all assets owned by the operator to calculate its voting power
	opFuncToIterateAssets := func(assetID string, state *assetstypes.OperatorAssetInfo) error {
		price := prices[assetID]
		decimal := decimals[assetID]
		usdValue = usdValue.Add(CalculateUSDValue(state.TotalAmount, price.Value, decimal, price.Decimal))
		// calculate the token amount from the share for the operator
		selfAmount, err := delegationkeeper.TokensFromShares(state.OperatorShare, state.TotalShare, state.TotalAmount)
		if err != nil {
			return err
		}
		selfUSDValue = selfUSDValue.Add(CalculateUSDValue(selfAmount, price.Value, decimal, price.Decimal))
		return nil
	}
	err := k.assetsKeeper.IteratorAssetsForOperator(ctx, operator, assetsFilter, opFuncToIterateAssets)
	if err != nil {
		return sdkmath.LegacyDec{}, sdkmath.LegacyDec{}, err
	}
	return usdValue, selfUSDValue, nil
}

// UpdateVotingPower update the voting power of the specified AVS and its operators at
// the end of epoch.
func (k *Keeper) UpdateVotingPower(ctx sdk.Context, avsAddr string) error {
	// get assets supported by the AVS
	assets, err := k.avsKeeper.GetAvsSupportedAssets(ctx, avsAddr)
	if err != nil {
		return err
	}
	if assets == nil {
		return nil
	}

	// get the prices and decimals of assets
	decimals, err := k.assetsKeeper.GetAssetsDecimal(ctx, assets)
	if err != nil {
		return err
	}
	prices, err := k.oracleKeeper.GetMultipleAssetsPrices(ctx, assets)
	if err != nil {
		return err
	}
	// update the voting power of operators and AVS
	avsVotingPower := sdkmath.LegacyNewDec(0)
	opFunc := func(operator string, votingPower *sdkmath.LegacyDec) error {
		// clear the old voting power for the operator
		*votingPower = sdkmath.LegacyNewDec(0)
		usdValue, selfUSDValue, err := k.CalculateUSDValueForOperator(ctx, operator, assets, decimals, prices)
		if err != nil {
			return err
		}
		// check if self USD value is more than the minimum self delegation.
		minimumSelfDelegation, err := k.avsKeeper.GetAVSMinimumSelfDelegation(ctx, avsAddr)
		if err != nil {
			return err
		}
		if selfUSDValue.GTE(minimumSelfDelegation) {
			*votingPower = votingPower.Add(usdValue)
		}
		avsVotingPower = avsVotingPower.Add(*votingPower)
		return nil
	}
	// iterate all operators of the AVS to update their voting power
	// and calculate the voting power for AVS
	err = k.IterateOperatorsForAVS(ctx, avsAddr, true, opFunc)
	if err != nil {
		return err
	}
	// set the voting power for AVS
	err = k.SetAVSUSDValue(ctx, avsAddr, avsVotingPower)
	if err != nil {
		return err
	}
	return nil
}

// ClearPreConsensusPK clears the previous consensus public key for all operators
func (k *Keeper) ClearPreConsensusPK(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(
		store,
		[]byte{operatortypes.BytePrefixForOperatorAndChainIDToPrevConsKey},
	)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
	return nil
}

// EndBlock : update the assets' share when their prices change
func (k *Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// todo: need to consider the calling order
	avsList, err := k.avsKeeper.GetEpochEndAVSs(ctx)
	if err != nil {
		panic(err)
	}
	for _, avs := range avsList {
		err = k.UpdateVotingPower(ctx, avs)
		if err != nil {
			panic(err)
		}
	}

	err = k.ClearPreConsensusPK(ctx)
	if err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}
