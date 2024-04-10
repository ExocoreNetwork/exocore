package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TokensFromShares calculate the token amount of provided shares, truncated to Int
func (k Keeper) TokensFromShares(stakerShare, totalShare sdkmath.LegacyDec, operatorAmount sdkmath.Int) sdkmath.Int {
	return (stakerShare.MulInt(operatorAmount)).Quo(totalShare).TruncateInt()
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the validator has no tokens.
func (k Keeper) SharesFromTokens(totalShare sdkmath.LegacyDec, stakerAmount, operatorAmount sdkmath.Int) (sdkmath.LegacyDec, error) {
	if v.Tokens.IsZero() {
		return math.LegacyZeroDec(), ErrInsufficientShares
	}

	return v.GetDelegatorShares().MulInt(amt).QuoInt(v.GetTokens()), nil
}

// CalculateShare calculates the S_j
// S_j = S * T_j / T, `S` and `T` is the current asset share and amount of operator,
// and the T_j represents the change in staker's asset amount when some external
// operations occur, such as: delegation, undelegation and non-instantaneous slashing.
// A special case is the initial delegation, when T = 0 and S = 0, so T_j / T is undefined.
// For the initial delegation, delegator j who delegates T_j tokens receive S_j = T_j shares.
func (k Keeper) CalculateShare(ctx sdk.Context, operator sdk.AccAddress, assetID string, amount sdkmath.Int) (*sdkmath.LegacyDec, error) {
	// get the total share of operator
	info, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, operator, assetID)
	if err != nil {
		return nil, err
	}

	var share sdkmath.LegacyDec
	if info.TotalShare.IsZero() {
		// the first delegation to a validator sets the exchange rate to one
		share = sdk.NewDecFromInt(amount)
	} else {
		shares, err := v.SharesFromTokens(amount)
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}

	v.Tokens = v.Tokens.Add(amount)
	v.DelegatorShares = v.DelegatorShares.Add(issuedShares)

	return v, issuedShares
}
