package keeper

import (
	"cosmossdk.io/math"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Based on the epoch, AllocateTokens performs reward and fee distribution to all validators.
func (k Keeper) AllocateTokens(ctx sdk.Context, totalPreviousPower int64) error {
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// transfer collected fees to the distribution module account
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt); err != nil {
		return err
	}

	feePool := k.GetFeePool(ctx)
	if totalPreviousPower == 0 {
		feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected...)
		k.SetFeePool(ctx, feePool)
		return nil
	}

	// calculate fraction allocated to exocore validators
	remaining := feesCollected
	communityTax, err := k.GetCommunityTax(ctx)
	if err != nil {
		return err
	}
	feeMultiplier := feesCollected.MulDecTruncate(math.LegacyOneDec().Sub(communityTax))

	// allocate tokens proportionally to voting power of different validators
	// TODO: Consider parallelizing later
	allValidators := k.StakingKeeper.GetAllExocoreValidators(ctx) // GetAllValidators(suite.Ctx)
	for i, val := range allValidators {
		pk, err := val.ConsPubKey()
		if err != nil {
			ctx.Logger().Error("Failed to deserialize public key; skipping", "error", err, "i", i)
			continue
		}
		validatorDetail, found := k.StakingKeeper.ValidatorByConsAddrForChainID(
			ctx, sdk.GetConsAddress(pk), avstypes.ChainIDWithoutRevision(ctx.ChainID()),
		)
		if !found {
			ctx.Logger().Error("Operator address not found; skipping", "consAddress", sdk.GetConsAddress(pk), "i", i)
			continue
		}
		powerFraction := math.LegacyNewDec(val.Power).QuoTruncate(math.LegacyNewDec(totalPreviousPower))
		reward := feeMultiplier.MulDecTruncate(powerFraction)
		k.AllocateTokensToValidator(ctx, validatorDetail, reward, feePool)
		remaining = remaining.Sub(reward)
	}

	// allocate community funding
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.SetFeePool(ctx, feePool)
	return nil
}

// AllocateTokensToValidator allocate tokens to a particular validator,
// splitting according to commission.
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins, feePool *types.FeePool) {
	rate := val.GetCommission()
	commission := tokens.MulDec(rate)
	shared := tokens.Sub(commission)
	valBz := val.GetOperator()

	// update current commission
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCommission,
		sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		sdk.NewAttribute(types.EventTypeCommission, val.GetOperator().String()),
	))
	currentCommission := k.GetValidatorAccumulatedCommission(ctx, valBz)
	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	k.SetValidatorAccumulatedCommission(ctx, valBz, currentCommission)
	// update current rewards, i.e. the rewards to stakers
	// if the rewards do not exist it's fine, we will just add to zero.
	// allocate share tokens to all stakers of this operator.
	k.AllocateTokensToStakers(ctx, val.GetOperator(), shared, feePool)

	// update outstanding rewards
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRewards,
		sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
	))

	// ValidatorOutstandingRewards is the rewards of a validator address.
	outstanding := k.GetValidatorOutstandingRewards(ctx, valBz)
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	k.SetValidatorOutstandingRewards(ctx, valBz, outstanding)
}

func (k Keeper) AllocateTokensToStakers(ctx sdk.Context, operatorAddress sdk.ValAddress, rewardToAllStakers sdk.DecCoins, feePool *types.FeePool) {
	avsList, err := k.StakingKeeper.GetOptedInAVSForOperator(ctx, operatorAddress.String())
	if err != nil {
		ctx.Logger().Error("avs address lists not found; skipping")
	}
	stakersPowerMap, curTotoalStakersPowers := make(map[string]math.LegacyDec), math.LegacyNewDec(1)
	for _, avsAddress := range avsList {
		avsAssets, err := k.StakingKeeper.GetAVSSupportedAssets(ctx, avsAddress)
		if err != nil {
			ctx.Logger().Error("avs address lists not found; skipping")
		}
		for assetID := range avsAssets {
			stakerList, err := k.StakingKeeper.GetStakersByOperator(ctx, operatorAddress.String(), assetID)
			if err != nil {
				ctx.Logger().Error("staker lists not found; skipping")
			}
			for _, staker := range stakerList.Stakers {
				if curStakerPower, err := k.StakingKeeper.CalculateUSDValueForStaker(ctx, staker, avsAddress, operatorAddress.Bytes()); err != nil {
					ctx.Logger().Error("curStakerPower error", err)
				} else {
					stakersPowerMap[staker] = curStakerPower
					curTotoalStakersPowers.Add(curStakerPower)
				}
			}
		}
	}

	for staker, stakerPower := range stakersPowerMap {
		powerFraction := stakerPower.QuoTruncate(curTotoalStakersPowers)
		rewardToSingleStaker := rewardToAllStakers.MulDecTruncate(powerFraction)
		k.AllocateTokensToSingleStaker(ctx, staker, rewardToSingleStaker)
		rewardToAllStakers = rewardToAllStakers.Sub(rewardToSingleStaker)
	}
	feePool.CommunityPool = feePool.CommunityPool.Add(rewardToAllStakers...)
}

func (k Keeper) AllocateTokensToSingleStaker(ctx sdk.Context, stakerAddress string, reward sdk.DecCoins) {
	currentStakerRewards := k.GetStakerRewards(ctx, stakerAddress)
	currentStakerRewards.Rewards = currentStakerRewards.Rewards.Add(reward...)
	k.SetStakerRewards(ctx, stakerAddress, currentStakerRewards)
}
