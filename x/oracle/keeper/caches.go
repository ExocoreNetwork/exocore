package keeper

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type cacheItem struct {
	feederId  int32
	pSources  []*types.PriceWithSource
	validator string
	// power     *big.Int
}

// used to track validator change
type cacheValidator struct {
	validators []*types.Validator
	update     bool
}

// used to track params change
type cacheParams struct {
	params *params
	update bool
}

// memory cache
func (k *Keeper) addCache(item *cacheItem) {

}

// memory cache
func (k *Keeper) removeCache(item *cacheItem) {

}

func (k *Keeper) commitCache() {

}

// persist cache into KV
func (k *Keeper) updateRecentTxs() {
	//TODO: add

	//TODO: delete, use map in the KVStore, so actually dont need to implement the exactly 'delete' at firt, just remove all blocks before maxDistance
}

func (k *Keeper) updateCacheValidators() bool {
	return false
}
func (k *Keeper) updateCacheParams() bool {
	return false
}

// from KVStore
func (k *Keeper) getCaches(ctx sdk.Context) map[uint64][]*cacheItem {
	return nil
}

func (k *Keeper) getParams4CacheRecover(ctx sdk.Context) map[uint64]*params {
	return nil
}

func (k *Keeper) clearCaches(ctx sdk.Context) {

}

// from KVStore
func (k *Keeper) getCacheValidators(ctx sdk.Context) map[string]*big.Int {
	return nil
}

func (k *Keeper) recacheAggregatorContext(ctx sdk.Context, agc *aggregatorContext) bool {

	from := uint64(ctx.BlockHeight()) - maxNonce
	to := uint64(ctx.BlockHeight()) - 1

	//fill validators
	validators, ok := k.GetValidators(ctx)
	recenetParams := k.getParams4CacheRecover(ctx)
	if !ok || recenetParams == nil {
		//no cache, this is the very first running, so go to initial proces instead
		return false
	}
	h := validators.Block
	if h > from {
		from = h
	}
	for _, v := range validators.ValidatorList {
		agc.validatorsPower[v.Operator], _ = new(big.Int).SetString(v.Power, 10)
	}

	recentMsgs := k.getCaches(ctx)

	for ; from < to; from++ {
		//fill params
		for h, p := range recenetParams {
			prev := uint64(0)
			if h <= from && h > prev {
				pTmp := params(*p)
				agc.params = &pTmp
				if prev > 0 {
					//TODO: safe delete
					delete(recenetParams, prev)
				}
				prev = h
			}
		}

		agc.prepareRound(ctx, from)

		if msgs := recentMsgs[from+1]; msgs != nil {
			for _, msg := range msgs {
				//these messages are retreived for recache, just skip the validation check and fill the memory cache
				agc.fillPrice(&types.MsgCreatePrice{
					Creator:  msg.validator,
					FeederId: msg.feederId,
					Prices:   msg.pSources,
				})
			}
		}
		agc.sealRound(ctx)
	}

	agc.prepareRound(ctx, to)

	return true
}

func (k *Keeper) initAggregatorContext(ctx sdk.Context, agc *aggregatorContext) error {
	//set params
	p := k.GetParams(ctx)
	m := make(map[uint64]*types.Params)
	m[uint64(ctx.BlockHeight())] = &p
	k.setParams4CacheRecover(m) //used to trace tokenFeeder's update during cache recover
	pTmp := params(p)
	agc.params = &pTmp

	totalPower := big.NewInt(0)
	vList := make([]*types.Validator, 0)
	k.stakingKeeper.IterateBondedValidatorsByPower(ctx, func(index int64, validator stakingTypes.ValidatorI) bool {
		power := big.NewInt(validator.GetConsensusPower(validator.GetBondedTokens()))
		addr := string(validator.GetOperator())
		agc.validatorsPower[addr] = power
		totalPower = new(big.Int).Add(totalPower, power)
		vList = append(vList, &types.Validator{Operator: addr, Power: power.String()})
		return false
	})
	agc.totalPower = totalPower
	k.SetValidators(ctx, types.Validators{
		Block:         uint64(ctx.BlockHeight()),
		ValidatorList: vList,
	})

	agc.prepareRound(ctx, uint64(ctx.BlockHeight())-1)
	return nil
}

// this is a mock serves as a placeholder
func (k *Keeper) setParams4CacheRecover(p map[uint64]*types.Params) {}
