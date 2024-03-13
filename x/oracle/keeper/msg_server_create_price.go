package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePrice(goCtx context.Context, msg *types.MsgCreatePrice) (*types.MsgCreatePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	/**
		1. aggregator.rInfo.Tokenid->status == 0(1 ignore and return)
		2. basedBlock is valid [roundInfo.basedBlock, *+5], each base only allow for one submit each validator, window for submition is 5 blocks while every validator only allowed to submit at most 3 transactions each round
		3. check the rule fulfilled(sources check), check the decimal of the 1st mathc the params' definition(among prices the decimal had been checked in ante stage), timestamp:later than previous block's timestamp, [not future than now(+1s), this is checked in anteHandler], timestamp verification is not necessary
	**/

	newItem, caches, _ := GetAggregatorContext(ctx, k.Keeper).NewCreatePrice(ctx, msg)

	if caches != nil {
		if newItem != nil {
			k.AppendPriceTR(ctx, newItem.TokenId, newItem.PriceTR)
			//TODO: move related caches
			cs.RemoveCache(caches, k)
		} else {
			cs.AddCache(caches, k)
		}
	}

	return &types.MsgCreatePriceResponse{}, nil
}
