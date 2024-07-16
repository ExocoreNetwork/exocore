package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"errors"
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Based on the epoch, AllocateTokens performs reward and fee distribution to all validators based
// on the F1 fee distribution specification.
func (k Keeper) AllocateTokens(ctx sdk.Context, totalPreviousPower int64) error {
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)
	// transfer collected fees to the distribution module account
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt); err != nil {
		return err
	}
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return err
	}
	if totalPreviousPower == 0 {
		if err := k.FeePool.Set(ctx, types.FeePool{CommunityPool: feePool.CommunityPool.Add(feesCollected...)}); err != nil {
			return err
		}
	}
	// calculate fraction allocated to validators
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

	if err := k.FeePool.Set(ctx, types.FeePool{DecimalPool: feePool.DecimalPool.Add(re...)}); err != nil {
		return err
	}

	return nil

}

// AllocateTokensToValidator allocate tokens to a particular validator,
// splitting according to commission.
func (k Keeper) AllocateTokensToValidator(ctx sdk.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) error {
	// split tokens between validator and delegators according to commission
	rate := val.GetCommission()
	commission := tokens.MulDec(rate)
	shared := tokens.Sub(commission)
	valBz := val.GetOperator().Bytes()
	//valBz, err := k.StakingKeeper.Validator() Valida GetExocoreValidator().StringToBytes(val.GetOperator())
	//if err != nil {
	//	return err
	//}

	// update current commission
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCommission,
		sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		sdk.NewAttribute(types.EventTypeCommission, val.GetOperator().String()),
	))
	currentCommission, err := k.ValidatorsAccumulatedCommission.Get(ctx, valBz)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return err
	}
	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	err = k.ValidatorsAccumulatedCommission.Set(ctx, valBz, currentCommission)
	if err != nil {
		return err
	}

	// update current rewards
	currentRewards, err := k.ValidatorCurrentRewards.Get(ctx, valBz)
	// if the rewards do not exist it's fine, we will just add to zero.
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return err
	}

	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	err = k.ValidatorCurrentRewards.Set(ctx, valBz, currentRewards)
	if err != nil {
		return err
	}

	// update outstanding rewards
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeRewards,
		sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
	))

	outstanding, err := k.ValidatorOutstandingRewards.Get(ctx, valBz)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return err
	}

	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	return k.ValidatorOutstandingRewards.Set(ctx, valBz, outstanding)
}
