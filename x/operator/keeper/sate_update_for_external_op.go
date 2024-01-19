package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) IncreasedOptedInAssets(ctx sdk.Context, stakerId, assetId, operatorAddr string, opAmount sdkmath.Int) error {
	//get the AVS opted-in by the operator
	avsList, err := k.GetOptedInAVSForOperator(ctx, operatorAddr)
	if err != nil {
		return err
	}

	//get price and priceDecimal from oracle
	price, priceDecimal, err := k.oracleKeeper.GetSpecifiedAssetsPrice(ctx, assetId)
	if err != nil {
		return err
	}

	opUsdValue := opAmount.Mul()

}

func (k Keeper) DecreaseOptedInAssets(ctx sdk.Context, stakerId, assetId, operatorAddr string, opAmount sdkmath.Int) error {

}

// OptIn call this function to opt in AVS
func (k Keeper) OptIn(ctx sdk.Context, OperatorAddress sdk.AccAddress, AVSAddr string) error {

	return nil
}

// OptOut call this function to opt out of AVS
func (k Keeper) OptOut(ctx sdk.Context, OperatorAddress sdk.AccAddress, AVSAddr string) error {

	return nil
}
