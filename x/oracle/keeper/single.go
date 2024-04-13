package keeper

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/aggregator"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/cache"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// GetAggregatorContext returns singleton aggregatorContext used to calculate final price for each round of each tokenFeeder
func GetAggregatorContext(ctx sdk.Context, k Keeper) *aggregator.AggregatorContext {
	if agc != nil {
		return agc
	}

	c := GetCaches()
	c.ResetCaches()
	agc = aggregator.NewAggregatorContext()
	if ok := recacheAggregatorContext(ctx, agc, k, c); !ok {
		//this is the very first time oracle has been started, fill relalted info as initialization
		initAggregatorContext(ctx, agc, k, c)
	}
	return agc
}

func recacheAggregatorContext(ctx sdk.Context, agc *aggregator.AggregatorContext, k Keeper, c *cache.Cache) bool {

	from := uint64(ctx.BlockHeight()) - common.MaxNonce
	to := uint64(ctx.BlockHeight()) - 1

	h, ok := k.GetValidatorUpdateBlock(ctx)
	recentParamsMap := k.GetAllRecentParamsAsMap(ctx)
	if !ok || recentParamsMap == nil {
		//no cache, this is the very first running, so go to initial proces instead
		return false
	}

	if h.Block > from {
		from = h.Block
	}

	totalPower := big.NewInt(0)
	validatorPowers := make(map[string]*big.Int)
	k.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(validator.GetBondedTokens()))
		addr := string(validator.GetOperator())
		validatorPowers[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		return false
	})
	agc.SetValidatorPowers(validatorPowers)
	//TODO: test only
	if k.GetLastTotalPower(ctx).BigInt().Cmp(totalPower) != 0 {
		panic("something wrong when get validatorsPower from staking module")
	}

	//reset validators
	c.AddCache(cache.CacheItemV(validatorPowers))

	recentMsgs := k.GetAllRecentMsgAsMap(ctx)
	var pTmp common.Params
	for ; from < to; from++ {
		//fill params
		prev := uint64(0)
		for b, recentParams := range recentParamsMap {
			if b <= from && b > prev {
				pTmp = common.Params(*recentParams)
				agc.SetParams(&pTmp)
				prev = b
			}
		}

		agc.PrepareRound(ctx, from)

		if msgs := recentMsgs[from+1]; msgs != nil {
			for _, msg := range msgs {
				//these messages are retreived for recache, just skip the validation check and fill the memory cache
				agc.FillPrice(&types.MsgCreatePrice{
					Creator:  msg.Validator,
					FeederId: msg.FeederId,
					Prices:   msg.PSources,
				})
			}
		}
		agc.SealRound(ctx, false)
	}

	//fill params cache
	c.AddCache(cache.CacheItemP(&pTmp))

	agc.PrepareRound(ctx, to)

	return true
}

func initAggregatorContext(ctx sdk.Context, agc *aggregator.AggregatorContext, k common.KeeperOracle, c *cache.Cache) error {
	//set params
	p := k.GetParams(ctx)
	m := make(map[uint64]*types.Params)
	m[uint64(ctx.BlockHeight())] = &p
	//	k.setParams4CacheRecover(m) //used to trace tokenFeeder's update during cache recover
	pTmp := common.Params(p)
	agc.SetParams(&pTmp)
	//set params cache
	c.AddCache(cache.CacheItemP(&pTmp))

	totalPower := big.NewInt(0)
	validatorPowers := make(map[string]*big.Int)
	k.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingtypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(validator.GetBondedTokens()))
		//addr := string(validator.GetOperator())
		addr := validator.GetOperator().String()
		//agc.validatorsPower[addr] = power
		validatorPowers[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		return false
	})
	//	agc.SetTotalPower(totalPower)
	agc.SetValidatorPowers(validatorPowers)
	if k.GetLastTotalPower(ctx).BigInt().Cmp(totalPower) != 0 {
		panic("-")
	}

	//set validatorPower cache
	c.AddCache(cache.CacheItemV(validatorPowers))

	agc.PrepareRound(ctx, uint64(ctx.BlockHeight())-1)
	return nil
}

func ResetAggregatorContext() {
	agc = nil
}

func ResetCache() {
	cs = nil
}
