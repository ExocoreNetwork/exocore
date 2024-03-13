package cache

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/aggregator"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var zeroBig = big.NewInt(0)

type cacheItemV map[string]*big.Int
type cacheItemP *common.Params
type cacheItemM struct {
	feederId  int32
	pSources  []*types.PriceWithSource
	validator string
}

type Cache struct {
	msg        cacheMsgs
	validators *cacheValidator
	params     *cacheParams
}

type cacheMsgs map[int32][]*cacheItemM

// used to track validator change
type cacheValidator struct {
	validators map[string]*big.Int
	update     bool
}

// used to track params change
type cacheParams struct {
	params *common.Params
	update bool
}

func (c cacheMsgs) add(item *cacheItemM) {
	c[item.feederId] = append(c[item.feederId], item)
}
func (c cacheMsgs) remove(item *cacheItemM) {
	delete(c, item.feederId)
}

func (c cacheMsgs) commit(ctx sdk.Context, k common.KeeperOracle) {
	block := uint64(ctx.BlockHeight())
	recentMsgs := types.RecentMsg{
		Block: block,
		Msgs:  make([]*types.MsgItem, 0),
	}
	for _, msgs4Feeder := range c {
		for _, msg := range msgs4Feeder {
			recentMsgs.Msgs = append(recentMsgs.Msgs, &types.MsgItem{
				FeederId:  msg.feederId,
				PSources:  msg.pSources,
				Validator: msg.validator,
			})
		}
	}
	index, _ := k.GetIndexRecentMsg(ctx)
	for i, b := range index.Index {
		if b >= block-common.MaxNonce {
			index.Index = index.Index[i:]
			break
		}
		k.RemoveRecentMsg(ctx, b)
	}
	k.SetRecentMsg(ctx, recentMsgs)
	index.Index = append(index.Index, block)
	k.SetIndexRecentMsg(ctx, index)
}

func (c *cacheValidator) add(validators map[string]*big.Int) {
	for operator, newPower := range validators {
		if power, ok := c.validators[operator]; ok {
			if newPower.Cmp(zeroBig) == 0 {
				delete(c.validators, operator)
				c.update = true
			} else if power.Cmp(newPower) != 0 {
				c.validators[operator].Set(newPower)
				c.update = true
			}
		} else {
			c.update = true
			np := *newPower
			c.validators[operator] = &np
		}
	}
}

func (c *cacheValidator) commit(ctx sdk.Context, k common.KeeperOracle) {
	block := uint64(ctx.BlockHeight())
	k.SetValidatorUpdateBlock(ctx, types.ValidatorUpdateBlock{Block: block})
}

func (c *cacheParams) add(p *common.Params) {
	//params' update is triggered when params is actually updated, so no need to do comparison here, just udpate and mark the flag
	//TODO: add comparison check, that's something should be done for validation
	c.params = p
	c.update = true
}

func (c *cacheParams) commit(ctx sdk.Context, k common.KeeperOracle) {
	block := uint64(ctx.BlockHeight())
	index, _ := k.GetIndexRecentParams(ctx)
	for i, b := range index.Index {
		if b >= block-common.MaxNonce {
			index.Index = index.Index[i:]
			break
		}
		k.RemoveRecentParams(ctx, b)
	}
	//remove and append for KVStore
	k.SetIndexRecentParams(ctx, index)
	index.Index = append(index.Index, block)
	k.SetIndexRecentParams(ctx, index)
}

// memory cache
func (c *Cache) AddCache(i any, k common.KeeperOracle) {
	switch item := i.(type) {
	case *cacheItemM:
		c.msg.add(item)
		//	case *params:
	case cacheItemP:
		c.params.add(item)
	case cacheItemV:
		c.validators.add(item)
	default:
		panic("no other types are support")
	}
}
func (c *Cache) RemoveCache(i any, k common.KeeperOracle) {
	switch item := i.(type) {
	case *cacheItemM:
		c.msg.remove(item)
	default:
	}
}

func (c *Cache) CommitCache(ctx sdk.Context, reset bool, k common.KeeperOracle) {
	if len(c.msg) > 0 {
		c.msg.commit(ctx, k)
		c.msg = make(map[int32][]*cacheItemM)
	}

	if c.validators.update {
		c.validators.commit(ctx, k)
		c.validators.update = false
	}

	if c.params.update {
		c.params.commit(ctx, k)
		c.params.update = false
	}
	if reset {
		c.ResetCaches()
	}
}

func (c *Cache) ResetCaches() {
	*c = *(NewCache())
}

func NewCache() *Cache {
	return &Cache{
		msg: make(map[int32][]*cacheItemM),
		validators: &cacheValidator{
			validators: make(map[string]*big.Int),
		},
		params: &cacheParams{
			params: &common.Params{},
		},
	}
}

func (c *Cache) RecacheAggregatorContext(ctx sdk.Context, agc *aggregator.AggregatorContext, k common.KeeperOracle) bool {

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
	k.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingTypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(validator.GetBondedTokens()))
		addr := string(validator.GetOperator())
		validatorPowers[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		return false
	})
	agc.SetValidatorPowers(validatorPowers)
	agc.SetTotalPower(totalPower)
	//TODO: test only
	if k.GetLastTotalPower(ctx).Cmp(totalPower) != 0 {
		panic("something wrong when get validatorsPower from staking module")
	}

	//reset validators
	c.AddCache(cacheItemV(validatorPowers), k)

	recentMsgs := k.GetAllRecentMsgAsMap(ctx)
	var pTmp common.Params
	for ; from < to; from++ {
		//fill params
		for b, recentParams := range recentParamsMap {
			prev := uint64(0)
			if b <= from && b > prev {
				pTmp = common.Params(*recentParams)
				agc.SetParams(&pTmp)
				if prev > 0 {
					//TODO: safe delete
					delete(recentParamsMap, prev)
				}
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
		agc.SealRound(ctx)
	}

	//fill params cache
	c.AddCache(cacheItemP(&pTmp), k)

	agc.PrepareRound(ctx, to)

	return true
}

func (c *Cache) InitAggregatorContext(ctx sdk.Context, agc *aggregator.AggregatorContext, k common.KeeperOracle) error {
	//set params
	p := k.GetParams(ctx)
	m := make(map[uint64]*types.Params)
	m[uint64(ctx.BlockHeight())] = &p
	//	k.setParams4CacheRecover(m) //used to trace tokenFeeder's update during cache recover
	pTmp := common.Params(p)
	agc.SetParams(&pTmp)
	//set params cache
	c.AddCache(cacheItemP(&pTmp), k)

	totalPower := big.NewInt(0)
	validatorPowers := make(map[string]*big.Int)
	k.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingTypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(validator.GetBondedTokens()))
		addr := string(validator.GetOperator())
		//agc.validatorsPower[addr] = power
		validatorPowers[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		return false
	})
	agc.SetTotalPower(totalPower)
	agc.SetValidatorPowers(validatorPowers)
	if k.GetLastTotalPower(ctx).Cmp(totalPower) != 0 {
		panic("-")
	}

	//set validatorPower cache
	c.AddCache(cacheItemV(validatorPowers), k)

	agc.PrepareRound(ctx, uint64(ctx.BlockHeight())-1)
	return nil
}
