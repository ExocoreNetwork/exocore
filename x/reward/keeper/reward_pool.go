package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type rewardPool struct {
	ctx         sdk.Context
	k           Keeper
	banker      bankkeeper.Keeper
	distributor types.Distributor
	types.Pool
}

func newRewardPool(ctx sdk.Context, k Keeper, banker bankkeeper.Keeper, distributor types.Distributor, p types.Pool) *rewardPool {
	return &rewardPool{
		ctx:         ctx,
		k:           k,
		banker:      banker,
		distributor: distributor,
	}
}

func (p rewardPool) getRewards(earningAddress string) (sdk.Coins, bool) {
	for _, reward := range p.Pool.Rewards {
		if reward.EarningsAddr == earningAddress {
			return reward.Coins, true
		}
	}
	return sdk.Coins{}, false
}

func (p *rewardPool) AddReward(earningAddress string, coin sdk.Coin) {
	defer func() {
		p.k.Logger(p.ctx).Debug("adding rewards in pool", "pool", p.Name, "earningAddress", earningAddress, "coin", coin.String(), "amount", coin.Amount)
		p.k.setPool(p.ctx, p.Pool)
	}()

	if coin.Amount.IsZero() {
		return
	}

	for i, reward := range p.Rewards {
		if reward.EarningsAddr == earningAddress {
			p.Rewards[i].Coins = reward.Coins.Add(coin)
			return
		}
	}

	p.Rewards = append(p.Rewards, types.Pool_Reward{
		EarningsAddr: earningAddress,
		Coins:        sdk.NewCoins(coin),
	})
}
