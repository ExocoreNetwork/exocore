package keeper

import (
	"context"
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePrice(goCtx context.Context, msg *types.MsgCreatePrice) (*types.MsgCreatePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO: Handling the message
	_ = ctx

	/**
		1. aggregator.rInfo.Tokenid->status == 0(1 ignore and return)
		2. basedBlock is valid [roundInfo.basedBlock, *+5], each base only allow for one submit each validator, window for submition is 5 blocks while every validator only allowed to submit at most 3 transactions each round
		3. check the rule fulfilled(sources check), check the decimal of the 1st mathc the params' definition(among prices the decimal had been checked in ante stage), timestamp:later than previous block's timestamp, [not future than now(+1s), this is checked in anteHandler], timestamp verification is not necessary
	**/

	newItem, caches, _ := getAggregatorContext(ctx, k.Keeper).newCreatePrice(ctx, msg)

	if caches != nil {
		k.addCache(caches)
	}

	if newItem != nil {
		k.AppendPriceTR(ctx, newItem.tokenId, newItem.priceTR)
		//TODO: move related caches
		k.removeCache(nil)
	}

	return &types.MsgCreatePriceResponse{}, nil
}

func getAggregatorContext(ctx sdk.Context, k Keeper) *aggregatorContext {
	if agc != nil {
		return agc
	}

	//initialize the aggregatorContext, normally triggered when node restart
	k.clearCaches(ctx)
	agc = &aggregatorContext{
		validatorsPower: make(map[string]*big.Int),
		totalPower:      big.NewInt(0),
		rounds:          make(map[int32]*roundInfo),
		aggregators:     make(map[int32]*worker),
	}
	if validators := k.getCacheValidators(ctx); validators != nil {
		agc.validatorsPower = validators
		for _, v := range validators {
			agc.totalPower = new(big.Int).Add(agc.totalPower, v)
		}
	}

	p := params(k.GetParams(ctx))
	agc.params = &p

	//replay the recentMsgs to recover the cache
	agc.prepareRound(ctx, uint64(ctx.BlockHeight())-uint64(maxNonce))
	agc.recache(uint64(ctx.BlockHeight())-uint64(maxNonce), uint64(ctx.BlockHeight())-1, k)

	//TODO: 1. prepare roundInfo with feeder and ctx, prepare response for status: 2->1
	recentItems := k.getCaches(ctx)

	for _, item := range recentItems {

	}

	return agc
}
