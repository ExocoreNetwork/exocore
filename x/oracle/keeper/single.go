package keeper

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/aggregator"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/cache"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	logger := k.Logger(ctx)
	from := ctx.BlockHeight() - int64(common.MaxNonce) + 1
	to := ctx.BlockHeight()

	h, ok := k.GetValidatorUpdateBlock(ctx)
	recentParamsMap := k.GetAllRecentParamsAsMap(ctx)
	if !ok || len(recentParamsMap) == 0 {
		logger.Info("no validatorUpdateBlock found, go to initial process", "height", ctx.BlockHeight())
		// no cache, this is the very first running, so go to initial process instead
		return false
	}

	if int64(h.Block) >= from {
		from = int64(h.Block) + 1
	}

	logger.Info("recacheAggregatorContext", "from", from, "to", to, "height", ctx.BlockHeight())
	totalPower := big.NewInt(0)
	validatorPowers := make(map[string]*big.Int)
	validatorSet := k.GetAllExocoreValidators(ctx)
	for _, v := range validatorSet {
		validatorPowers[sdk.AccAddress(v.Address).String()] = big.NewInt(v.Power)
		totalPower = new(big.Int).Add(totalPower, big.NewInt(v.Power))
	}
	agc.SetValidatorPowers(validatorPowers)
	// TODO: test only
	if k.GetLastTotalPower(ctx).BigInt().Cmp(totalPower) != 0 {
		ctx.Logger().Error("something wrong when get validatorsPower from dogfood module")
	}

	// reset validators
	c.AddCache(cache.ItemV(validatorPowers))

	recentMsgs := k.GetAllRecentMsgAsMap(ctx)

	var pTmp common.Params
	if from >= to {
		// backwards compatible for that the validatorUpdateBlock updated every block
		prev := int64(0)
		for b, p := range recentParamsMap {
			if b > prev {
				// pTmp be set at least once, since len(recentParamsMap)>0
				pTmp = common.Params(*p)
				prev = b
			}
		}
		agc.SetParams(&pTmp)
		setCommonParams(types.Params(pTmp))
	} else {
		for ; from < to; from++ {
			// fill params
			prev := int64(0)
			for b, recentParams := range recentParamsMap {
				if b <= from && b > prev {
					pTmp = common.Params(*recentParams)
					agc.SetParams(&pTmp)
					prev = b
					setCommonParams(*recentParams)
				}
			}

			agc.PrepareRoundBeginBlock(ctx, uint64(from))

			if msgs := recentMsgs[from]; msgs != nil {
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
			ctxReplay := ctx.WithBlockHeight(from)
			agc.SealRound(ctxReplay, false)
		}
	}

	var pRet common.Params
	if updated := c.GetCache(cache.ItemP(&pRet)); !updated {
		c.AddCache(cache.ItemP(&pTmp))
	}

	return true
}

func initAggregatorContext(ctx sdk.Context, agc *aggregator.AggregatorContext, k common.KeeperOracle, c *cache.Cache) {
	ctx.Logger().Info("initAggregatorContext", "height", ctx.BlockHeight())
	// set params
	p := k.GetParams(ctx)
	pTmp := common.Params(p)
	agc.SetParams(&pTmp)
	// set params cache
	c.AddCache(cache.ItemP(&pTmp))
	setCommonParams(p)
	totalPower := big.NewInt(0)
	validatorPowers := make(map[string]*big.Int)
	validatorSet := k.GetAllExocoreValidators(ctx)
	for _, v := range validatorSet {
		validatorPowers[sdk.AccAddress(v.Address).String()] = big.NewInt(v.Power)
		totalPower = new(big.Int).Add(totalPower, big.NewInt(v.Power))
	}

	agc.SetValidatorPowers(validatorPowers)
	// TODO: test only
	if k.GetLastTotalPower(ctx).BigInt().Cmp(totalPower) != 0 {
		ctx.Logger().Error("something wrong when get validatorsPower from dogfood module")
	}
	// set validatorPower cache
	c.AddCache(cache.ItemV(validatorPowers))

	agc.PrepareRoundBeginBlock(ctx, uint64(ctx.BlockHeight()))
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

func setCommonParams(p types.Params) {
	common.MaxNonce = p.MaxNonce
	common.ThresholdA = p.ThresholdA
	common.ThresholdB = p.ThresholdB
	common.MaxDetID = p.MaxDetId
	common.Mode = p.Mode
}
