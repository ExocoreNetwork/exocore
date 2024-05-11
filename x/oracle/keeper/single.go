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
		// this is the very first time oracle has been started, fill relalted info as initialization
		initAggregatorContext(ctx, agc, k, c)
	} else {
		// this is when a node restart and use the persistent state to refill cache, we don't need to commit these data again
		c.SkipCommit()
	}
	return agc
}

func recacheAggregatorContext(ctx sdk.Context, agc *aggregator.AggregatorContext, k Keeper, c *cache.Cache) bool {
	from := ctx.BlockHeight() - common.MaxNonce
	to := ctx.BlockHeight() - 1

	h, ok := k.GetValidatorUpdateBlock(ctx)
	recentParamsMap := k.GetAllRecentParamsAsMap(ctx)
	if !ok || len(recentParamsMap) == 0 {
		// no cache, this is the very first running, so go to initial process instead
		return false
	}

	if int64(h.Block) > from {
		from = int64(h.Block)
	}

	totalPower := big.NewInt(0)
	validatorPowers := make(map[string]*big.Int)
	k.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(sdk.DefaultPowerReduction))
		addr := string(validator.GetOperator())
		validatorPowers[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		return false
	})
	agc.SetValidatorPowers(validatorPowers)
	// TODO: test only
	if k.GetLastTotalPower(ctx).BigInt().Cmp(totalPower) != 0 {
		ctx.Logger().Error("something wrong when get validatorsPower from staking module")
	}

	// reset validators
	c.AddCache(cache.ItemV(validatorPowers))

	recentMsgs := k.GetAllRecentMsgAsMap(ctx)
	var pTmp common.Params
	for ; from < to; from++ {
		// fill params
		prev := int64(0)
		for b, recentParams := range recentParamsMap {
			if b <= from && b > prev {
				pTmp = common.Params(*recentParams)
				agc.SetParams(&pTmp)
				prev = b
			}
		}

		agc.PrepareRound(ctx, uint64(from))

		if msgs := recentMsgs[from+1]; msgs != nil {
			for _, msg := range msgs {
				// these messages are retreived for recache, just skip the validation check and fill the memory cache
				//nolint
				agc.FillPrice(&types.MsgCreatePrice{
					Creator:  msg.Validator,
					FeederID: msg.FeederID,
					Prices:   msg.PSources,
				})
			}
		}
		agc.SealRound(ctx, false)
	}

	// fill params cache
	c.AddCache(cache.ItemP(&pTmp))

	agc.PrepareRound(ctx, uint64(to))

	return true
}

func initAggregatorContext(ctx sdk.Context, agc *aggregator.AggregatorContext, k common.KeeperOracle, c *cache.Cache) {
	// set params
	p := k.GetParams(ctx)
	pTmp := common.Params(p)
	agc.SetParams(&pTmp)
	// set params cache
	c.AddCache(cache.ItemP(&pTmp))

	totalPower := big.NewInt(0)
	validatorPowers := make(map[string]*big.Int)
	k.IterateBondedValidatorsByPower(ctx, func(_ int64, validator stakingtypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(sdk.DefaultPowerReduction))
		addr := validator.GetOperator().String()
		validatorPowers[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		return false
	})
	agc.SetValidatorPowers(validatorPowers)

	// set validatorPower cache
	c.AddCache(cache.ItemV(validatorPowers))

	agc.PrepareRound(ctx, uint64(ctx.BlockHeight()-1))
}

func ResetAggregatorContext() {
	agc = nil
}

func ResetCache() {
	cs = nil
}
