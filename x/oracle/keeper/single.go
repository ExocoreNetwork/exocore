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

var agc, agcCheckTx *aggregator.AggregatorContext

func GetCaches() *cache.Cache {
	if cs != nil {
		return cs
	}
	cs = cache.NewCache()
	return cs
}

// GetAggregatorContext returns singleton aggregatorContext used to calculate final price for each round of each tokenFeeder
func GetAggregatorContext(ctx sdk.Context, k Keeper) *aggregator.AggregatorContext {
	if ctx.IsCheckTx() {
		if agcCheckTx != nil {
			return agcCheckTx
		}
		if agc == nil {
			c := GetCaches()
			c.ResetCaches()
			agcCheckTx = aggregator.NewAggregatorContext()
			if ok := recacheAggregatorContext(ctx, agcCheckTx, k, c); !ok {
				// this is the very first time oracle has been started, fill relalted info as initialization
				initAggregatorContext(ctx, agcCheckTx, k, c)
			}
			return agcCheckTx
		}
		agcCheckTx = agc.Copy4CheckTx()
		return agcCheckTx
	}

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
	from := ctx.BlockHeight() - int64(common.MaxNonce)
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
		addr := validator.GetOperator().String()
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
	for ; from < to; from++ {
		// fill params
		prev := int64(0)
		for b, recentParams := range recentParamsMap {
			if b <= from && b > prev {
				agc.SetParams(recentParams)
				prev = b
				setCommonParams(recentParams)
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
	var p *types.Params
	var b int64
	if from >= to {
		// backwards compatible for that the validatorUpdateBlock updated every block
		prev := int64(0)
		for b, p = range recentParamsMap {
			if b > prev {
				// pTmp be set at least once, since len(recentParamsMap)>0
				prev = b
			}
		}
		agc.SetParams(p)
		setCommonParams(p)
	}

	var pRet cache.ItemP
	if updated := c.GetCache(&pRet); !updated {
		c.AddCache(cache.ItemP(*p))
	}
	// fill params cache
	agc.PrepareRound(ctx, uint64(to))

	return true
}

func initAggregatorContext(ctx sdk.Context, agc *aggregator.AggregatorContext, k common.KeeperOracle, c *cache.Cache) {
	// set params
	p := k.GetParams(ctx)
	agc.SetParams(&p)
	// set params cache
	c.AddCache(cache.ItemP(p))
	setCommonParams(&p)

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

func ResetAggregatorContextCheckTx() {
	agcCheckTx = nil
}

func setCommonParams(p *types.Params) {
	common.MaxNonce = p.MaxNonce
	common.ThresholdA = p.ThresholdA
	common.ThresholdB = p.ThresholdB
	common.MaxDetID = p.MaxDetId
	common.Mode = p.Mode
}
