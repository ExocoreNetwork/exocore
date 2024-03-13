package keeper

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var zeroBig = big.NewInt(0)
var caches *cache

type cacheItemV map[string]*big.Int
type cacheItemP *params
type cacheItemM struct {
	feederId  int32
	pSources  []*types.PriceWithSource
	validator string
}

type cache struct {
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
	params *params
	update bool
}

func (c cacheMsgs) add(item *cacheItemM) {
	c[item.feederId] = append(c[item.feederId], item)
}
func (c cacheMsgs) remove(item *cacheItemM) {
	delete(c, item.feederId)
}

func (c cacheMsgs) commit(ctx sdk.Context, k *Keeper) {
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
		if b >= block-maxNonce {
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

func (c *cacheValidator) commit(ctx sdk.Context, k *Keeper) {
	block := uint64(ctx.BlockHeight())
	k.SetValidatorUpdateBlock(ctx, types.ValidatorUpdateBlock{Block: block})
}

func (c *cacheParams) add(p *params) {
	//params' update is triggered when params is actually updated, so no need to do comparison here, just udpate and mark the flag
	//TODO: add comparison check, that's something should be done for validation
	c.params = p
	c.update = true
}

func (c *cacheParams) commit(ctx sdk.Context, k *Keeper) {
	block := uint64(ctx.BlockHeight())
	index, _ := k.GetIndexRecentParams(ctx)
	for i, b := range index.Index {
		if b >= block-maxNonce {
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
func (k *Keeper) addCache(i any) {
	switch item := i.(type) {
	case *cacheItemM:
		caches.msg.add(item)
		//	case *params:
	case cacheItemP:
		caches.params.add(item)
	case cacheItemV:
		caches.validators.add(item)
	default:
		panic("no other types are support")
	}
}
func (k *Keeper) removeCache(i any) {
	switch item := i.(type) {
	case *cacheItemM:
		caches.msg.remove(item)
	default:
	}
}

func (k *Keeper) commitCache(ctx sdk.Context, reset bool) {
	if len(caches.msg) > 0 {
		caches.msg.commit(ctx, k)
		caches.msg = make(map[int32][]*cacheItemM)
	}

	if caches.validators.update {
		caches.validators.commit(ctx, k)
		caches.validators.update = false
	}

	if caches.params.update {
		caches.params.commit(ctx, k)
		caches.params.update = false
	}
	if reset {
		k.resetCaches(ctx)
	}
}

func (k *Keeper) resetCaches(ctx sdk.Context) {
	caches = &cache{
		msg: make(map[int32][]*cacheItemM),
		validators: &cacheValidator{
			validators: make(map[string]*big.Int),
		},
		params: &cacheParams{
			params: &params{},
		},
	}
}

func (k *Keeper) recacheAggregatorContext(ctx sdk.Context, agc *aggregatorContext) bool {

	from := uint64(ctx.BlockHeight()) - maxNonce
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
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingTypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(validator.GetBondedTokens()))
		addr := string(validator.GetOperator())
		agc.validatorsPower[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		return false
	})
	agc.totalPower = k.stakingKeeper.GetLastTotalPower(ctx).BigInt()
	//TODO: test only
	if agc.totalPower.Cmp(totalPower) != 0 {
		panic("something wrong when get validatorsPower from staking module")
	}

	//reset validators
	k.addCache(cacheItemV(agc.validatorsPower))

	recentMsgs := k.GetAllRecentMsgAsMap(ctx)

	for ; from < to; from++ {
		//fill params
		for b, recentParams := range recentParamsMap {
			prev := uint64(0)
			if b <= from && b > prev {
				pTmp := params(*recentParams)
				agc.params = &pTmp
				if prev > 0 {
					//TODO: safe delete
					delete(recentParamsMap, prev)
				}
				prev = b
			}
		}

		agc.prepareRound(ctx, from)

		if msgs := recentMsgs[from+1]; msgs != nil {
			for _, msg := range msgs {
				//these messages are retreived for recache, just skip the validation check and fill the memory cache
				agc.fillPrice(&types.MsgCreatePrice{
					Creator:  msg.Validator,
					FeederId: msg.FeederId,
					Prices:   msg.PSources,
				})
			}
		}
		agc.sealRound(ctx)
	}

	//fill params cache
	k.addCache(cacheItemP(agc.params))

	agc.prepareRound(ctx, to)

	return true
}

func (k *Keeper) initAggregatorContext(ctx sdk.Context, agc *aggregatorContext) error {
	//set params
	p := k.GetParams(ctx)
	m := make(map[uint64]*types.Params)
	m[uint64(ctx.BlockHeight())] = &p
	//	k.setParams4CacheRecover(m) //used to trace tokenFeeder's update during cache recover
	pTmp := params(p)
	agc.params = &pTmp
	//set params cache
	k.addCache(cacheItemP(agc.params))

	totalPower := big.NewInt(0)
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingTypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(validator.GetBondedTokens()))
		addr := string(validator.GetOperator())
		agc.validatorsPower[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		return false
	})
	agc.totalPower = totalPower

	//set validatorPower cache
	k.addCache(cacheItemV(agc.validatorsPower))

	agc.prepareRound(ctx, uint64(ctx.BlockHeight())-1)
	return nil
}
