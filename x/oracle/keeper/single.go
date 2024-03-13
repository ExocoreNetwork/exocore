package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/aggregator"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/cache"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var cs *cache.Cache

var agc *aggregator.AggregatorContext

func GetCaches() *cache.Cache {
	if cs != nil {
		return cs
	}
	cs = cache.NewCache()
	return cs
}

func GetAggregatorContext(ctx sdk.Context, k Keeper) *aggregator.AggregatorContext {
	if agc != nil {
		return agc
	}

	c := GetCaches()
	c.ResetCaches()
	agc = aggregator.NewAggregatorContext()
	if ok := c.RecacheAggregatorContext(ctx, agc, k); !ok {
		//this is the very first time oracle has been started, fill relalted info as initialization
		c.InitAggregatorContext(ctx, agc, k)
	}
	return agc
}
