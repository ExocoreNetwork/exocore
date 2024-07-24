package keeper

import (
	"log"

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
	feePool := k.FeePool
	if totalPreviousPower == 0 {
		k.FeePool = types.FeePool{
			CommunityPool: feePool.CommunityPool.Add(feesCollected...),
		}
	}
	// calculate fraction allocated to exocore validators
	remaining := feesCollected
	communityTax, err := k.GetCommunityTax(ctx)
	if err != nil {
		return err
	}
	feeMultiplier := feesCollected.MulDecTruncate(math.LegacyOneDec().Sub(communityTax))
	// allocate tokens proportionally to voting power of different validators
	validatorUpdates := k.StakingKeeper.GetValidatorUpdates(ctx)
	for _, vu := range validatorUpdates {
		powerFraction := math.LegacyNewDec(vu.Power).QuoTruncate(math.LegacyNewDec(totalPreviousPower))
		reward := feeMultiplier.MulDecTruncate(powerFraction)
		pubKey, _ := cryptocodec.FromTmProtoPublicKey(vu.PubKey)
		consAddr := sdk.ConsAddress(pubKey.Address().String())
		validator := k.StakingKeeper.ValidatorByConsAddr(ctx, consAddr)
		if err = k.AllocateTokensToValidator(ctx, validator, reward); err != nil {
			return err
		}
		remaining = remaining.Sub(reward)
	}
	// send to community pool and set remainder in fee pool
	amt, re := remaining.TruncateDecimal()
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, types.ProtocolPoolModuleName, amt); err != nil {
		return err
	}

	// set ToDistribute in protocolpool to keep track of continuous funds distribution
	if err := k.poolKeeper.SetToDistribute(ctx, amt, k.GetAuthority()); err != nil { // TODO: this should be distribution module account
		return err
	}
	//	k.FeePool = types.FeePool{DecimalPool: k.FeePool.DecimalPool.Add(re...)}
	k.FeePool.DecimalPool = k.FeePool.DecimalPool.Add(re...)
	return nil
}

// AllocateTokensToValidator allocate tokens to a particular validator,
// splitting according to commission.
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) error {
	// split tokens between validator and delegators according to commission
	rate := val.GetCommission()
	commission := tokens.MulDec(rate)
	shared := tokens.Sub(commission)
	valBz := val.GetOperator().String()

	// update current commission
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCommission,
		sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		sdk.NewAttribute(types.EventTypeCommission, val.GetOperator().String()),
	))
	if currentCommission, ok := k.ValidatorsAccumulatedCommission[valBz]; ok {
		currentCommission.Commission = currentCommission.Commission.Add(commission...)
		k.ValidatorsAccumulatedCommission[valBz] = currentCommission
	} else {
		log.Printf("currentCommission %s didn't exist", currentCommission)
		// No need to return here
	}

	// update current rewards
	// if the rewards do not exist it's fine, we will just add to zero.
	if currentRewards, ok := k.ValidatorCurrentRewards[valBz]; ok {
		currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
		k.ValidatorCurrentRewards[valBz] = currentRewards
	}

	// update outstanding rewards
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRewards,
		sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
	))

	if outstanding, ok := k.ValidatorOutstandingRewards[valBz]; ok {
		outstanding.Rewards = outstanding.Rewards.Add(tokens...)
		k.ValidatorOutstandingRewards[valBz] = outstanding
	} else {
		log.Printf("ValidatorOutstandingRewards for %s didn't exist", valBz)
	}
	return nil
}
