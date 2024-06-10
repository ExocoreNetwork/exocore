package keeper

import (
	"context"
	"strconv"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CreatePrice proposes price for new round of specific tokenFeeder
func (ms msgServer) CreatePrice(goCtx context.Context, msg *types.MsgCreatePrice) (*types.MsgCreatePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newItem, caches, err := GetAggregatorContext(ctx, ms.Keeper).NewCreatePrice(ctx, msg)
	if err != nil {
		return nil, err
	}

	logger := ms.Keeper.Logger(ctx)
	logger.Info("add price proposal for aggregation", "feederID", msg.FeederID, "basedBlock", msg.BasedBlock, "proposer", msg.Creator)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeCreatePrice,
		sdk.NewAttribute(types.AttributeKeyFeederID, strconv.FormatUint(msg.FeederID, 10)),
		sdk.NewAttribute(types.AttributeKeyBasedBlock, strconv.FormatUint(msg.BasedBlock, 10)),
		sdk.NewAttribute(types.AttributeKeyProposer, msg.Creator),
	),
	)

	if caches == nil {
		return &types.MsgCreatePriceResponse{}, nil
	}
	if newItem != nil {
		ms.AppendPriceTR(ctx, newItem.TokenID, newItem.PriceTR)

		logger.Info("final price aggregation done", "feederID", msg.FeederID, "roundID", newItem.PriceTR.RoundID, "price", newItem.PriceTR.Price)

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeCreatePrice,
			sdk.NewAttribute(types.AttributeKeyFeederID, strconv.FormatUint(msg.FeederID, 10)),
			sdk.NewAttribute(types.AttributeKeyRoundID, strconv.FormatUint(newItem.PriceTR.RoundID, 10)),
			sdk.NewAttribute(types.AttributeKeyFinalPrice, newItem.PriceTR.Price),
			sdk.NewAttribute(types.AttributeKeyPriceUpdated, types.AttributeValuePriceUpdatedSuccess)),
		)
		if !ctx.IsCheckTx() {
			cs.RemoveCache(caches)
		}
	} else if !ctx.IsCheckTx() {
		cs.AddCache(caches)
	}

	return &types.MsgCreatePriceResponse{}, nil
}
