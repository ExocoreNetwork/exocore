package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/reward/types"
)

// sendRewards internal method called with ClaimedRewardPeriodDetail of a single reward program
func (k Keeper) sendRewards(ctx sdk.Context, rewards []*types.ClaimedRewardPeriodDetail, addr string, rewardProgramDenom string) (sdk.Coin, error) {
	amount := sdk.NewInt64Coin(rewardProgramDenom, 0)

	if len(rewards) == 0 {
		return amount, nil
	}

	for _, reward := range rewards {
		amount.Denom = reward.GetClaimPeriodReward().Denom
		amount = amount.Add(reward.GetClaimPeriodReward())
	}

	return k.sendCoinsToAccount(ctx, amount, addr)
}

// sendCoinsToAccount is mainly for `SendCoinsFromModuleToAccount`
func (k Keeper) sendCoinsToAccount(ctx sdk.Context, amount sdk.Coin, addr string) (sdk.Coin, error) {
	if amount.IsZero() {
		return sdk.NewInt64Coin(amount.GetDenom(), 0), nil
	}

	acc, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return sdk.NewInt64Coin(amount.Denom, 0), err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acc, sdk.NewCoins(amount))
	if err != nil {
		return sdk.NewInt64Coin(amount.Denom, 0), err
	}

	return amount, nil
}
