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
		// When slash from the delegated share, we can slash sequentially from all delegated operators
		// according to the order of the operator addresses in the KV store. We can also proportionally
		// slash from all operators to whom the related staker has already delegated, if we think
		// this approach is better

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

		// slash from pending undelegations
		if remainAmount.IsPositive() {
			opFunc := func(undelegation *delegationtypes.UndelegationRecord) (bool, error) {
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
				return false, nil
			}
			err = k.IterateUndelegationsByStakerAndAsset(ctx, stakerID, assetID, true, opFunc)
			if err != nil {
				return err
			}
		}

		// slash from the delegated share
		if remainAmount.IsPositive() {
			opFunc := func(keys *delegationtypes.SingleDelegationInfoReq, _ *delegationtypes.DelegationAmounts) (bool, error) {
				opAccAddr, err := sdk.AccAddressFromBech32(keys.OperatorAddr)
				if err != nil {
					return true, err
				}
				slashShare, err := k.CalculateSlashShare(ctx, opAccAddr, stakerID, assetID, remainAmount)
				if err != nil {
					return true, err
				}
				actualSlashAmount, err := k.RemoveShare(ctx, false, opAccAddr, stakerID, assetID, slashShare)
				if err != nil {
					return true, err
				}
				remainAmount = remainAmount.Sub(actualSlashAmount)
				if !remainAmount.IsPositive() {
					return true, nil
				}
				return false, nil
			}
			err = k.IterateDelegationsForStakerAndAsset(ctx, stakerID, assetID, opFunc)
			if err != nil {
				return err
			}
		}
		// In this case, we only print a log as a reminder. This situation will only occur when the total slashing amount
		// from the client chain and Exocore chain is greater than the total staking amount.
		if remainAmount.IsPositive() {
			ctx.Logger().Info("all staking funds has been slashed, the remaining amount is:", "remainAmount", remainAmount)
		}
	}
	return nil
}
