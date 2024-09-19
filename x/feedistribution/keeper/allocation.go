package keeper

import (
	"sort"

	"cosmossdk.io/math"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Based on the epoch, AllocateTokens performs reward and fee distribution to all validators.
func (k Keeper) AllocateTokens(ctx sdk.Context, totalPreviousPower int64) error {
	logger := k.Logger()
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
	logger.Info("Allocate tokens to all validators", "feesCollected amount is ", feesCollected)
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
	logger := k.Logger()
	valBz := val.GetOperator()
	accAddr := sdk.AccAddress(valBz)
	ops, err := k.StakingKeeper.OperatorInfo(ctx, accAddr.String())
	if err != nil {
		ctx.Logger().Error("Failed to get operator info", "error", err)
	}
	commission := tokens.MulDec(ops.GetCommission().Rate)
	shared := tokens.Sub(commission)
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
	operatorAccAddress := sdk.AccAddress(valBz)
	k.AllocateTokensToStakers(ctx, operatorAccAddress, shared, feePool)

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
	logger.Info("Allocate tokens to validator successfully", "allocated amount is", outstanding.Rewards.String())
}

func (k Keeper) AllocateTokensToStakers(ctx sdk.Context, operatorAddress sdk.AccAddress, rewardToAllStakers sdk.DecCoins, feePool *types.FeePool) {
	logger := k.Logger()
	logger.Info("Allocate to stakers of operatorAddress:", operatorAddress.String())
	avsList, err := k.StakingKeeper.GetOptedInAVSForOperator(ctx, operatorAddress.String())
	if err != nil {
		logger.Debug("avs address lists not found; skipping")
		return
	}
	stakersPowerMap, curTotalStakersPowers := make(map[string]math.LegacyDec), math.LegacyNewDec(0)
	globalStakerAddressList := make([]string, 0)
	for _, avsAddress := range avsList {
		avsAssets, err := k.StakingKeeper.GetAVSSupportedAssets(ctx, avsAddress)
		if err != nil {
			logger.Debug("avs address lists not found; skipping")
			continue
		}
		for assetID := range avsAssets {
			stakerList, err := k.StakingKeeper.GetStakersByOperator(ctx, operatorAddress.String(), assetID)
			if err != nil {
				logger.Debug("staker lists not found; skipping")
				continue
			}
			for _, staker := range stakerList.Stakers {
				if curStakerPower, err := k.StakingKeeper.CalculateUSDValueForStaker(ctx, staker, avsAddress, operatorAddress.Bytes()); err != nil {
					logger.Error("curStakerPower error", err)
				} else {
					stakersPowerMap[staker] = curStakerPower
					globalStakerAddressList = append(globalStakerAddressList, staker)
					curTotalStakersPowers = curTotalStakersPowers.Add(curStakerPower)
				}
			}
		}
	}
	sort.Slice(globalStakerAddressList, func(i, j int) bool {
		return stakersPowerMap[globalStakerAddressList[i]].GT(stakersPowerMap[globalStakerAddressList[j]])
	})
	remaining := rewardToAllStakers
	// allocate to stakers in voting power descending order if the curTotalStakersPower is positive
	for _, staker := range globalStakerAddressList {
		stakerPower := stakersPowerMap[staker]
		powerFraction := stakerPower.QuoTruncate(curTotalStakersPowers)
		rewardToSingleStaker := rewardToAllStakers.MulDecTruncate(powerFraction)
		k.AllocateTokensToSingleStaker(ctx, staker, rewardToSingleStaker)
		remaining = remaining.Sub(rewardToSingleStaker)
	}

	feePool.CommunityPool = feePool.CommunityPool.Add(rewardToAllStakers...)
	logger.Info("allocate tokens to stakers successfully", "allocated amount is", rewardToAllStakers.String())
}

func (k Keeper) AllocateTokensToSingleStaker(ctx sdk.Context, stakerAddress string, reward sdk.DecCoins) {
	logger := k.Logger()
	currentStakerRewards := k.GetStakerRewards(ctx, stakerAddress)
	currentStakerRewards.Rewards = currentStakerRewards.Rewards.Add(reward...)
	k.SetStakerRewards(ctx, stakerAddress, currentStakerRewards)
	logger.Info("allocate tokens to single staker successfully", "allocated amount is", currentStakerRewards.Rewards.String())
}
