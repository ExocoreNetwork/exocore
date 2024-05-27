package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// func (k Keeper) PricesAll(goCtx context.Context, req *types.QueryAllPricesRequest) (*types.QueryAllPricesResponse, error) {
//	if req == nil {
//		return nil, status.Error(codes.InvalidArgument, "invalid request")
//	}
//
//	var pricess []types.Prices
//	ctx := sdk.UnwrapSDKContext(goCtx)
//
//	store := ctx.KVStore(k.storeKey)
//	pricesStore := prefix.NewStore(store, types.KeyPrefix(types.PricesKeyPrefix))
//	pricesTokenStore := prefix.NewStore(pricesStore, types.PricesKey(tokenID))
//
//	pageRes, err := query.Paginate(pricesTokenStore, req.Pagination, func(key []byte, value []byte) error {
//		var prices types.Prices
//		if err := k.cdc.Unmarshal(value, &prices); err != nil {
//			return err
//		}
//
//		pricess = append(pricess, prices)
//		return nil
//	})
//
//	if err != nil {
//		return nil, status.Error(codes.Internal, err.Error())
//	}
//
//	return &types.QueryAllPricesResponse{Prices: pricess, Pagination: pageRes}, nil
//}

func (k Keeper) Prices(goCtx context.Context, req *types.QueryGetPricesRequest) (*types.QueryGetPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetPrices(
		ctx,
		req.TokenId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPricesResponse{Prices: val}, nil
}

// TODO: LatestPrice(tokenID)
