package keeper

import (
	"context"
	"strconv"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreatePrice proposes price for new round of specific tokenFeeder
func (k msgServer) CreatePrice(goCtx context.Context, msg *types.MsgCreatePrice) (*types.MsgCreatePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newItem, caches, err := GetAggregatorContext(ctx, k.Keeper).NewCreatePrice(ctx, msg)
	if err != nil {
		return nil, err
	}

	logger := k.Keeper.Logger(ctx)
	logger.Info("add price proposal for aggregation", "feederID", msg.FeederId, "basedBlock", msg.BasedBlock, "proposer", msg.Creator)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreatePrice,
		sdk.NewAttribute(types.AttributeKeyFeederID, strconv.Itoa(int(msg.FeederId))),
		sdk.NewAttribute(types.AttributeKeyBasedBlock, strconv.FormatInt(int64(msg.BasedBlock), 10)),
		sdk.NewAttribute(types.AttributeKeyProposer, msg.Creator),
	),
	)

	if caches != nil {
		if newItem != nil {
			k.AppendPriceTR(ctx, newItem.TokenId, newItem.PriceTR)

			logger.Info("final price aggregation done", "feederID", msg.FeederId, "roundID", newItem.PriceTR.RoundId, "price", newItem.PriceTR.Price)

			ctx.EventManager().EmitEvent(sdk.NewEvent(
				types.EventTypeCreatePrice,
				sdk.NewAttribute(types.AttributeKeyFeederID, strconv.Itoa(int(msg.FeederId))),
				sdk.NewAttribute(types.AttributeKeyRoundID, strconv.FormatInt(int64(newItem.PriceTR.RoundId), 10)),
				sdk.NewAttribute(types.AttributeKeyFinalPrice, newItem.PriceTR.Price),
				sdk.NewAttribute(types.AttributeKeyPriceUpdated, types.AttributeValuePriceUpdatedSuccess),
			),
			)
			cs.RemoveCache(caches)
		} else {
			cs.AddCache(caches)
		}
	}

	return &types.MsgCreatePriceResponse{}, nil
}
