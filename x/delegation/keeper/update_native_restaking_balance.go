package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) UpdateNativeRestakingBalance(
	ctx sdk.Context, stakerID, assetID string, amount sdkmath.Int,
) error {
	// todo: check if the assetID is native retaking token
	if amount.IsPositive() {
		// If the balance increases due to the client chain PoS staking reward, the increased
		// amount can be considered a virtual deposit event. However, the increased amount needs
		// to be manually delegated by the staker if they want it to contribute to voting power.
		// Of course, we can also treat it as both a virtual deposit and delegation event if we
		// think this approach is better. In that case, we would proportionally delegate the
		// increased amount to all operators to whom the related staker has already delegated
		// this native token.
		err := k.assetsKeeper.UpdateStakerAssetState(ctx, stakerID, assetID, types.DeltaStakerSingleAsset{
			TotalDepositAmount: amount,
			WithdrawableAmount: amount,
		})
		if err != nil {
			return err
		}
	} else if amount.IsNegative() {
		// If the balance decreases due to the client chain PoS slashing, the decreased amount
		// will be slashed from the withdrawable amount first, the pending undelegation second,
		// and the delegated share last if there is still a remaining amount that needs to be slashed.
		// When slash from the delegated share, we can proportionally decrease from all operators
		// to whom the related staker has already delegated.

		// slash from the withdrawable amount
		assetInfo, err := k.assetsKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetID)
		if err != nil {
			return err
		}
		slashFromWithdrawable := amount.Neg()
		remainAmount := slashFromWithdrawable.Sub(assetInfo.WithdrawableAmount)
		if remainAmount.IsPositive() {
			slashFromWithdrawable = assetInfo.WithdrawableAmount
		}
		err = k.assetsKeeper.UpdateStakerAssetState(ctx, stakerID, assetID, types.DeltaStakerSingleAsset{
			TotalDepositAmount: slashFromWithdrawable.Neg(),
			WithdrawableAmount: slashFromWithdrawable.Neg(),
		})
		if err != nil {
			return err
		}
		ctx.Logger().Info("UpdateNativeRestakingBalance slash from withdrawable amount", "stakerID", stakerID, "assetID", assetID, "slashFromWithdrawable", slashFromWithdrawable, "remainAmount", remainAmount)

		// slash from pending undelegations
		if remainAmount.IsPositive() {
			opFunc := func(undelegationKey string, undelegation *delegationtypes.UndelegationRecord) (bool, error) {
				// slash from the single undelegation
				slashAmount := remainAmount
				remainAmount = slashAmount.Sub(undelegation.ActualCompletedAmount)
				if remainAmount.IsPositive() {
					slashAmount = undelegation.ActualCompletedAmount
				}
				undelegation.ActualCompletedAmount = undelegation.ActualCompletedAmount.Sub(slashAmount)
				if !remainAmount.IsPositive() {
					// return ture to break the iteration if there isn't remaining amount to be slashed
					return true, nil
				}
				ctx.Logger().Info("UpdateNativeRestakingBalance slash from undelegation", "stakerID", stakerID, "assetID", assetID, "operator", undelegation.OperatorAddr, "undelegationKey", undelegationKey, "slashAmount", slashAmount, "remainAmount", remainAmount)
				return false, nil
			}
			err = k.IterateUndelegationsByStakerAndAsset(ctx, stakerID, assetID, true, opFunc)
			if err != nil {
				return err
			}
		}

		// slash from the delegated share
		// the delegated share will be proportionally decreased from all operators to
		// whom the related staker has already delegated
		if remainAmount.IsPositive() {
			// calculate the slash proportion
			totalDelegatedAmount, err := k.TotalDelegatedAmountForStakerAsset(ctx, stakerID, assetID)
			if err != nil {
				return err
			}
			slashProportion := sdkmath.LegacyNewDecFromBigInt(remainAmount.BigInt()).Quo(sdkmath.LegacyNewDecFromBigInt(totalDelegatedAmount.BigInt()))
			if slashProportion.GT(sdkmath.LegacyNewDec(1)) {
				slashProportion = sdkmath.LegacyNewDec(1)
			}
			opFunc := func(keys *delegationtypes.SingleDelegationInfoReq, delegationAmount *delegationtypes.DelegationAmounts) (bool, error) {
				opAccAddr, err := sdk.AccAddressFromBech32(keys.OperatorAddr)
				if err != nil {
					return true, err
				}
				slashShare := delegationAmount.UndelegatableShare.Mul(slashProportion)
				actualSlashAmount, err := k.RemoveShare(ctx, false, opAccAddr, stakerID, assetID, slashShare)
				if err != nil {
					return true, err
				}
				ctx.Logger().Info("UpdateNativeRestakingBalance slash from delegated share", "stakerID", stakerID, "assetID", assetID, "operator", keys.OperatorAddr, "slashProportion", "slashShare", slashShare, "actualSlashAmount", actualSlashAmount)
				return false, nil
			}
			err = k.IterateDelegationsForStakerAndAsset(ctx, stakerID, assetID, opFunc)
			if err != nil {
				return err
			}
			remainAmount = sdkmath.LegacyNewDec(1).Sub(slashProportion).MulInt(remainAmount).TruncateInt()
		}
		// In this case, we only print a log as a reminder. This situation will only occur when the total slashing amount
		// from the client chain and Exocore chain is greater than the total staking amount.
		if remainAmount.IsPositive() {
			ctx.Logger().Info("UpdateNativeRestakingBalance all staking funds has been slashed, the remaining amount is:", "stakerID", stakerID, "assetID", assetID, "remainAmount", remainAmount)
		}
	}
	return nil
}
