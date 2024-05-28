package cache

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var zeroBig = big.NewInt(0)

type (
	ItemV map[string]*big.Int
	ItemP *common.Params
	ItemM types.MsgItem
)

type Cache struct {
	msg        *cacheMsgs
	validators *cacheValidator
	params     *cacheParams
}

// type cacheMsgs map[uint64][]*ItemM
type cacheMsgs []*ItemM

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

func (c *cacheMsgs) add(item *ItemM) {
	*c = append(*c, item)
}

// remove removes all items with the same feederID
func (c *cacheMsgs) remove(item *ItemM) {
	var newCache []*ItemM
	for _, msg := range *c {
		if msg.FeederID != item.FeederID {
			newCache = append(newCache, msg)
		}
	}
	*c = newCache
}

func (c cacheMsgs) commit(ctx sdk.Context, k common.KeeperOracle) {
	block := uint64(ctx.BlockHeight())

	recentMsgs := types.RecentMsg{
		Block: block,
		Msgs:  make([]*types.MsgItem, 0),
	}

	for _, msg := range c {
		msgTmp := types.MsgItem(*msg)
		recentMsgs.Msgs = append(recentMsgs.Msgs, &msgTmp)
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
	i := 0
	for ; i < len(index.Index); i++ {
		b := index.Index[i]
		if b >= block-common.MaxNonce {
			index.Index = index.Index[i:]
			break
		}
		k.RemoveRecentParams(ctx, b)
	}
	index.Index = index.Index[i:]
	// remove and append for KVStore
	index.Index = append(index.Index, block)
	k.SetIndexRecentParams(ctx, index)

	p := types.Params(*c.params)
	k.SetRecentParams(ctx, types.RecentParams{
		Block:  block,
		Params: &p,
	})
}

// memory cache
func (c *Cache) AddCache(i any) {
	switch item := i.(type) {
	case *ItemM:
		c.msg.add(item)
	case ItemP:
		c.params.add(item)
	case ItemV:
		c.validators.add(item)
	default:
		panic("no other types are support")
	}
}

// RemoveCache removes all cached msgs with the same feederID
func (c *Cache) RemoveCache(i any) {
	if item, isItemM := i.(*ItemM); isItemM {
		c.msg.remove(item)
	}
}

func (c *Cache) GetCache(i any) bool {
	switch item := i.(type) {
	case ItemV:
		if item == nil {
			return false
		}
		for addr, power := range c.validators.validators {
			item[addr] = power
		}
		return c.validators.update
	case ItemP:
		if item == nil {
			return false
		}
		*item = *(c.params.params)
		return c.params.update
	case *([]*ItemM):
		if item == nil {
			return false
		}
		*item = *c.msg
		return len(*c.msg) > 0
	default:
		return false
	}
}

// SkipCommit skip real commit by setting the updage flag to false
func (c *Cache) SkipCommit() {
	c.validators.update = false
	c.params.update = false
}

func (c *Cache) CommitCache(ctx sdk.Context, reset bool, k common.KeeperOracle) {
	if len(*(c.msg)) > 0 {
		c.msg.commit(ctx, k)
		*(c.msg) = make([]*ItemM, 0)
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
		msg: new(cacheMsgs),
		validators: &cacheValidator{
			validators: make(map[string]*big.Int),
		},
		params: &cacheParams{
			params: &common.Params{},
		},
	}
}
