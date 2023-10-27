package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (p rewardPool) getRewards(validator sdk.ValAddress) (sdk.Coins, bool) {
	for _, reward := range p.Rewards {
		if reward.Validator.Equals(validator) {
			return reward.Coins, true
		}
	}

	return sdk.Coins{}, false
}

func (p *rewardPool) AddReward(validator sdk.ValAddress, coin sdk.Coin) {
	// basic check
	if coin.Amount.IsZero() {
		return
	}


	


}
