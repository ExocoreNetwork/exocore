package keeper

import (
	"errors"

	sdkmath "cosmossdk.io/math"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationkeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	oracletypes "github.com/ExocoreNetwork/exocore/x/oracle/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/xerrors"
)

// CalculateUSDValueForOperator calculates the total and self usd value for the
// operator according to the input assets filter and prices.
// This function will be used in slashing calculations and voting power updates per epoch.
// The inputs/outputs and calculation logic for these two cases are different,
// so an `isForSlash` flag is used to distinguish between them.
// When it's called by the voting power update, the needed outputs are the current total
// staking amount and the self-staking amount of the operator. The current total
// staking amount excludes the pending unbonding amount, so it's used to calculate the voting power.
// The self-staking amount is also needed to check if the operator's self-staking is sufficient.
// At the same time, the prices of all assets have been retrieved in the caller's function, so they
// are inputted as a parameter.
// When it's called by the slash execution, the needed output is the sum of the current total amount and
// the pending unbonding amount, because the undelegation also needs to be slashed. And the prices of
// all assets haven't been prepared by the caller, so the prices should be retrieved in this function.
func (k *Keeper) CalculateUSDValueForOperator(
	ctx sdk.Context,
	isForSlash bool,
	operator string,
	assetsFilter map[string]interface{},
	decimals map[string]uint32,
	prices map[string]oracletypes.Price,
) (operatortypes.OperatorUSDValue, error) {
	var err error
	ret := operatortypes.OperatorUSDValue{
		Staking:                 sdkmath.LegacyNewDec(0),
		SelfStaking:             sdkmath.LegacyNewDec(0),
		StakingAndWaitUnbonding: sdkmath.LegacyNewDec(0),
	}
	// iterate all assets owned by the operator to calculate its voting power
	opFuncToIterateAssets := func(assetID string, state *assetstypes.OperatorAssetInfo) error {
		//		var price operatortypes.Price
		var price oracletypes.Price
		var decimal uint32
		if isForSlash {
			// when calculated the USD value for slashing, the input prices map is null
			// so the price needs to be retrieved here
			price, err = k.oracleKeeper.GetSpecifiedAssetsPrice(ctx, assetID)
			if err != nil {
				// TODO: when assetID is not registered in oracle module, this error will finally lead to panic
				if !errors.Is(err, oracletypes.ErrGetPriceRoundNotFound) {
					return err
				}
				// TODO: for now, we ignore the error when the price round is not found and set the price to 1 to avoid panic
			}
			assetInfo, err := k.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
			if err != nil {
				return err
			}
			decimal = assetInfo.AssetBasicInfo.Decimals
			ret.StakingAndWaitUnbonding = ret.StakingAndWaitUnbonding.Add(CalculateUSDValue(state.TotalAmount.Add(state.WaitUnbondingAmount), price.Value, decimal, price.Decimal))
		} else {
			if prices == nil {
				return xerrors.Errorf("CalculateUSDValueForOperator, the input prices map is nil")
			}
			price = prices[assetID]
			decimal = decimals[assetID]
			ret.Staking = ret.Staking.Add(CalculateUSDValue(state.TotalAmount, price.Value, decimal, price.Decimal))
			// calculate the token amount from the share for the operator
			selfAmount, err := delegationkeeper.TokensFromShares(state.OperatorShare, state.TotalShare, state.TotalAmount)
			if err != nil {
				return err
			}
			ret.SelfStaking = ret.SelfStaking.Add(CalculateUSDValue(selfAmount, price.Value, decimal, price.Decimal))
		}
		return nil
	}
	err = k.assetsKeeper.IterateAssetsForOperator(ctx, false, operator, assetsFilter, opFuncToIterateAssets)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

// UpdateVotingPower update the voting power of the specified AVS and its operators at
// the end of epoch.
func (k *Keeper) UpdateVotingPower(ctx sdk.Context, avsAddr string) error {
	// get assets supported by the AVS
	assets, err := k.avsKeeper.GetAVSSupportedAssets(ctx, avsAddr)
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
	// TODO: for now, we ignore the error when the price round is not found and set the price to 1 to avoid panic
	if err != nil {
		// TODO: when assetID is not registered in oracle module, this error will finally lead to panic
		if !errors.Is(err, oracletypes.ErrGetPriceRoundNotFound) {
			return err
		}
		// TODO: for now, we ignore the error when the price round is not found and set the price to 1 to avoid panic
	}
	// update the voting power of operators and AVS
	avsVotingPower := sdkmath.LegacyNewDec(0)
	// check if self USD value is more than the minimum self delegation.
	minimumSelfDelegation, err := k.avsKeeper.GetAVSMinimumSelfDelegation(ctx, avsAddr)
	if err != nil {
		return err
	}
	opFunc := func(operator string, votingPower *sdkmath.LegacyDec) error {
		// clear the old voting power for the operator
		*votingPower = sdkmath.LegacyNewDec(0)
		usdValues, err := k.CalculateUSDValueForOperator(ctx, false, operator, assets, decimals, prices)
		if err != nil {
			return err
		}
		if usdValues.SelfStaking.GTE(minimumSelfDelegation) {
			*votingPower = votingPower.Add(usdValues.Staking)
			avsVotingPower = avsVotingPower.Add(*votingPower)
		}
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
	return []abci.ValidatorUpdate{}
}
