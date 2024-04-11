package keeper

import (
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
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

// CalculateShare calculates the S_j
// S_j = S * T_j / T, `S` and `T` is the current asset share and amount of operator,
// and the T_j represents the change in staker's asset amount when some external
// operations occur, such as: delegation, undelegation and non-instantaneous slashing.
// A special case is the initial delegation, when T = 0 and S = 0, so T_j / T is undefined.
// For the initial delegation, delegator j who delegates T_j tokens receive S_j = T_j shares.
func (k Keeper) CalculateShare(ctx sdk.Context, operator sdk.AccAddress, assetID string, amount sdkmath.Int) (share sdkmath.LegacyDec, err error) {
	// get the total share of operator
	info, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, operator, assetID)
	if err != nil {
		return share, err
	}

	if info.TotalShare.IsZero() {
		// the first delegation to a validator sets the exchange rate to one
		share = sdk.NewDecFromInt(amount)
	} else {
		share, err = SharesFromTokens(info.TotalShare, amount, info.TotalAmount)
		if err != nil {
			return share, err
		}
	}
	return share, nil
}

// ValidateUndeleagtionAmount validates that a given undelegation amount is
// valid based on upon the converted shares. If the amount is valid, the total
// amount of respective shares is returned, otherwise an error is returned.
func (k Keeper) ValidateUndeleagtionAmount(
	ctx sdk.Context, operator sdk.AccAddress, stakerID, assetID string, amount sdkmath.Int,
) (share sdkmath.LegacyDec, err error) {
	delegationInfo, err := k.GetSingleDelegationInfo(ctx, stakerID, assetID, operator.String())
	if err != nil {
		return share, err
	}

	info, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, operator, assetID)
	if err != nil {
		return share, err
	}

	share, err = SharesFromTokens(info.TotalShare, amount, info.TotalAmount)
	if err != nil {
		return share, err
	}

	if share.GT(delegationInfo.UndelegatableShare) {
		return share, delegationtypes.ErrInsufficientShares
	}

	// Depending on the share, amount can be smaller than unit amount(1stake).
	// If the remain amount after unbonding is smaller than the minimum share,
	// it's completely unbonded to avoid leaving dust shares.
	tolerance, err := SharesFromTokens(info.TotalShare, sdkmath.OneInt(), info.TotalAmount)
	if err != nil {
		return share, err
	}

	if delegationInfo.UndelegatableShare.Sub(share).LT(tolerance) {
		share = delegationInfo.UndelegatableShare
	}

	return share, nil
}

func (k Keeper) RemoveShareFromOperator(
	ctx sdk.Context, operator sdk.AccAddress, assetID string, share sdkmath.LegacyDec,
) (token sdkmath.Int, err error) {
	operatorAssetState, err := k.assetsKeeper.GetOperatorSpecifiedAssetInfo(ctx, operator, assetID)
	if err != nil {
		return token, err
	}
	if share.LT(operatorAssetState.TotalShare) {
		return token, delegationtypes.ErrInsufficientShares
	}

	var removedToken sdkmath.Int
	if operatorAssetState.TotalShare.Equal(share) {
		// last delegation share gets any trimmings
		removedToken = operatorAssetState.TotalAmount
	} else {
		// leave excess tokens in the validator
		// however fully use all the delegator shares
		removedToken, err = TokensFromShares(share, operatorAssetState.TotalShare, operatorAssetState.TotalAmount)
		if err != nil {
			return token, err
		}
	}

	delta := assetstype.DeltaOperatorSingleAsset{
		TotalAmount:         removedToken.Neg(),
		WaitUnbondingAmount: removedToken,
		TotalShare:          share.Neg(),
	}
	err = k.assetsKeeper.UpdateOperatorAssetState(ctx, operator, assetID, delta)
	if err != nil {
		return token, err
	}
	return token, nil
}
