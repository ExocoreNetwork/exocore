package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/reward/exported"
	"github.com/exocore/x/reward/types"
)

type rewardPool struct {
	types.Pool
	ctx    sdk.Context
	k      Keeper
	banker types.BankKeeper
}

func newPool(ctx sdk.Context, k Keeper, banker types.BankKeeper, p types.Pool) *rewardPool {
	return &rewardPool{
		ctx:    ctx,
		k:      k,
		banker: banker,
		Pool:   p,
	}
}

var _ exported.RewardPool = &rewardPool{}

// To be decided here!
// func (k Keeper) sendRewards(ctx sdk.Context, rewards []*types.MsgRewardDetail, addr string, rewardProgramDenom string) (sdk.Coin, error) {
// 	amount := sdk.NewInt64Coin(rewardProgramDenom, 0)

// 	if len(rewards) == 0 {
// 		return amount, nil
// 	}

// 	for _, reward := range rewards {
// 		amount.Denom = reward.GetClaimPeriodReward().Denom
// 		amount = amount.Add(reward.GetClaimPeriodReward())
// 	}

// 	return k.sendCoinsToAccount(ctx, amount, addr)
// }

// // sendCoinsToAccount is mainly for `SendCoinsFromModuleToAccount`
// func (k Keeper) sendCoinsToAccount(ctx sdk.Context, amount sdk.Coin, addr string) (sdk.Coin, error) {
// 	if amount.IsZero() {
// 		return sdk.NewInt64Coin(amount.GetDenom(), 0), nil
// 	}

// 	acc, err := sdk.AccAddressFromBech32(addr)
// 	if err != nil {
// 		return sdk.NewInt64Coin(amount.Denom, 0), err
// 	}

// 	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, acc, sdk.NewCoins(amount))
// 	if err != nil {
// 		return sdk.NewInt64Coin(amount.Denom, 0), err
// 	}

// 	return amount, nil
// }

// Get reward value of the corresponding address in the rewards pool
func (p rewardPool) getRewards(address sdk.ValAddress) (sdk.Coins, bool) {
	for _, reward := range p.Rewards {
		if reward.Validator.Equals(address) {
			return reward.Coins, true
		}
	}

	return sdk.Coins{}, false
}

// Add and record for the corresponding reward in the rewards pool
func (p *rewardPool) AddReward(address sdk.ValAddress, coin sdk.Coin) {
	defer func() {
		p.k.Logger(p.ctx).Debug("adding rewards in pool", "pool", p.Name, "validator", address.String(), "coin", coin.String())

		p.k.setPool(p.ctx, p.Pool)
	}()

	if coin.Amount.IsZero() {
		return
	}

	for i, reward := range p.Rewards {
		if reward.Validator.Equals(address) {
			p.Rewards[i].Coins = reward.Coins.Add(coin)
			return
		}
	}

	p.Rewards = append(p.Rewards, types.Pool_Reward{
		Validator: address,
		Coins:     sdk.NewCoins(coin),
	})
}

// Clear rewards of the specific address
func (p *rewardPool) ClearRewards(address sdk.ValAddress) {
	for i, reward := range p.Rewards {
		if reward.Validator.Equals(address) {
			p.k.Logger(p.ctx).Info("clearing rewards in pool", "pool", p.Name, "validator", address.String())

			p.Rewards = append(p.Rewards[:i], p.Rewards[i+1:]...)
			p.k.setPool(p.ctx, p.Pool)
			return
		}
	}
}
