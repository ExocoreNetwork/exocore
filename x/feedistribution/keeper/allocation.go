package keeper

import (
	"cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
	//
	// TODO: Consider parallelizing later
	validatorUpdates := k.StakingKeeper.GetValidatorUpdates(ctx)
	for _, vu := range validatorUpdates {
		powerFraction := math.LegacyNewDec(vu.Power).QuoTruncate(math.LegacyNewDec(totalPreviousPower))
		reward := feeMultiplier.MulDecTruncate(powerFraction)
		pubKey, _ := cryptocodec.FromTmProtoPublicKey(vu.PubKey)
		consAddr := sdk.ConsAddress(pubKey.Address().String())
		validator := k.StakingKeeper.ValidatorByConsAddr(ctx, consAddr)
		k.AllocateTokensToValidator(ctx, validator, reward)
		remaining = remaining.Sub(reward)
	}

	// allocate community funding
	feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)
	k.SetFeePool(ctx, feePool)
	return nil
}

// AllocateTokensToValidator allocate tokens to a particular validator,
// splitting according to commission.
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) {
	// split tokens between validator and delegators according to commission
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
	// update current rewards
	// if the rewards do not exist it's fine, we will just add to zero.
	currentRewards := k.GetValidatorCurrentRewards(ctx, valBz)
	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	k.SetValidatorCurrentRewards(ctx, valBz, currentRewards)

	// update outstanding rewards
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRewards,
		sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
	))

	outstanding := k.GetValidatorOutstandingRewards(ctx, valBz)
	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	k.SetValidatorOutstandingRewards(ctx, valBz, outstanding)
}
