package cache

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var zeroBig = big.NewInt(0)

type (
	CacheItemV map[string]*big.Int
	CacheItemP *common.Params
	CacheItemM struct {
		FeederId  uint64
		PSources  []*types.PriceWithSource
		Validator string
	}
)

type Cache struct {
	msg        cacheMsgs
	validators *cacheValidator
	params     *cacheParams
}

type cacheMsgs map[uint64][]*CacheItemM

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

func (c cacheMsgs) add(item *CacheItemM) {
	if ims, ok := c[item.FeederId]; ok {
		for _, im := range ims {
			if im.Validator == item.Validator {
				for _, p := range im.PSources {
					for _, pInput := range item.PSources {
						if p.SourceId == pInput.SourceId {
							p.Prices = append(p.Prices, pInput.Prices...)
							return
						}
					}
				}
				im.PSources = append(im.PSources, item.PSources...)
				return
			}
		}
	}
	c[item.FeederId] = append(c[item.FeederId], item)
}

func (c cacheMsgs) remove(item *CacheItemM) {
	delete(c, item.FeederId)
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
				FeederId:  msg.FeederId,
				PSources:  msg.PSources,
				Validator: msg.Validator,
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
	// params' update is triggered when params is actually updated, so no need to do comparison here, just udpate and mark the flag
	// TODO: add comparison check, that's something should be done for validation
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
	// remove and append for KVStore
	k.SetIndexRecentParams(ctx, index)
	index.Index = append(index.Index, block)
	k.SetIndexRecentParams(ctx, index)
}

// memory cache
// func (c *Cache) AddCache(i any, k common.KeeperOracle) {
func (c *Cache) AddCache(i any) {
	switch item := i.(type) {
	case *CacheItemM:
		c.msg.add(item)
		//	case *params:
	case CacheItemP:
		c.params.add(item)
	case CacheItemV:
		c.validators.add(item)
	default:
		panic("no other types are support")
	}
}

// func (c *Cache) RemoveCache(i any, k common.KeeperOracle) {
func (c *Cache) RemoveCache(i any) {
	switch item := i.(type) {
	case *CacheItemM:
		c.msg.remove(item)
	default:
	}
}

func (c *Cache) GetCache(i any) bool {
	switch item := i.(type) {
	case CacheItemV:
		if item == nil {
			return false
		}
		for addr, power := range c.validators.validators {
			item[addr] = power
		}
	case CacheItemP:
		if item == nil {
			return false
		}
		*item = *(c.params.params)
	case *[]*CacheItemM:
		if item == nil {
			return false
		}
		tmp := make([]*CacheItemM, 0, len(c.msg))
		for _, msgs := range c.msg {
			tmp = append(tmp, msgs...)
		}
		*item = tmp
	default:
		return false
	}
	return true
}

func (c *Cache) CommitCache(ctx sdk.Context, reset bool, k common.KeeperOracle) {
	if len(c.msg) > 0 {
		c.msg.commit(ctx, k)
		c.msg = make(map[uint64][]*CacheItemM)
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
		msg: make(map[uint64][]*CacheItemM),
		validators: &cacheValidator{
			validators: make(map[string]*big.Int),
		},
		params: &cacheParams{
			params: &common.Params{},
		},
	}
}
