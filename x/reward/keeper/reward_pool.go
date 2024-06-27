package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/reward/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type rewardPool struct {
	ctx         sdk.Context
	k           Keeper
	distributor types.Distributor
	banker      bankkeeper.Keeper
}

func newRewardPool(ctx sdk.Context, k Keeper, banker bankkeeper.Keeper, distributor types.Distributor, p types.Pool) *rewardPool {
	return &rewardPool{
		ctx:         ctx,
		k:           k,
		banker:      banker,
		distributor: distributor,
		staker:      staker,
		Pool:        p,
	}
}
