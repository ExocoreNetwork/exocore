package keeper

import (
	sdkmath "cosmossdk.io/math"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TokensFromShares calculate the token amount of provided shares, truncated to Int
func TokensFromShares(stakerShare, totalShare sdkmath.LegacyDec, operatorAmount sdkmath.Int) (sdkmath.Int, error) {
	if totalShare.IsZero() {
		return sdkmath.NewInt(0), delegationtypes.ErrDivisorIsZero
	}
	return (stakerShare.MulInt(operatorAmount)).Quo(totalShare).TruncateInt(), nil
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the validator has no tokens.
func SharesFromTokens(totalShare sdkmath.LegacyDec, stakerAmount, operatorAmount sdkmath.Int) (sdkmath.LegacyDec, error) {
	if operatorAmount.IsZero() {
		return sdkmath.LegacyZeroDec(), delegationtypes.ErrDivisorIsZero
	}
	return totalShare.MulInt(stakerAmount).QuoInt(operatorAmount), nil
}

// SharesFromTokensTruncated returns the truncated shares of a delegation given
// a bond amount. It returns an error if the validator has no tokens.
func SharesFromTokensTruncated(totalShare sdkmath.LegacyDec, stakerAmount, operatorAmount sdkmath.Int) (sdkmath.LegacyDec, error) {
	if operatorAmount.IsZero() {
		return sdkmath.LegacyZeroDec(), delegationtypes.ErrDivisorIsZero
	}
	return totalShare.MulInt(stakerAmount).QuoTruncate(sdkmath.LegacyNewDecFromInt(operatorAmount)), nil
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
		share, err = SharesFromTokens(info.TotalShare, amount, info.TotalAmount)
		if err != nil {
			return nil, err
		}
	}
	return &share, nil
}

// ValidateUndeleagtionAmount validates that a given undelegation amount is
// valid based on upon the converted shares. If the amount is valid, the total
// amount of respective shares is returned, otherwise an error is returned.
func (k Keeper) ValidateUndeleagtionAmount(
	ctx sdk.Context, operator sdk.AccAddress, stakerID, assetID string, amount sdkmath.Int,
) (*sdkmath.LegacyDec, error) {
	delegationInfo, err := k.GetSingleDelegationInfo(ctx, stakerID, assetID, operator.String())
	if err != nil {
		return nil, err
	}

	info, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, operator, assetID)
	if err != nil {
		return nil, err
	}

	shares, err := SharesFromTokens(info.TotalShare, amount, info.TotalAmount)
	if err != nil {
		return nil, err
	}

	sharesTruncated, err := SharesFromTokensTruncated(info.TotalShare, amount, info.TotalAmount)
	if err != nil {
		return nil, err
	}

	if sharesTruncated.GT(delegationInfo.UndelegatableShare) {
		return nil, delegationtypes.ErrInsufficientShares
	}

	// Depending on the share, amount can be smaller than unit amount(1stake).
	// If the remain amount after unbonding is smaller than the minimum share,
	// it's completely unbonded to avoid leaving dust shares.
	tolerance, err := SharesFromTokens(info.TotalShare, sdkmath.OneInt(), info.TotalAmount)
	if err != nil {
		return nil, err
	}

	if delegationInfo.UndelegatableShare.Sub(shares).LT(tolerance) {
		shares = delegationInfo.UndelegatableShare
	}

	return &shares, nil
}
