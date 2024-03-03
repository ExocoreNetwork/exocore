package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RoundDataAll(goCtx context.Context, req *types.QueryAllRoundDataRequest) (*types.QueryAllRoundDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var roundDatas []types.RoundData
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	roundDataStore := prefix.NewStore(store, types.KeyPrefix(types.RoundDataKeyPrefix))

	pageRes, err := query.Paginate(roundDataStore, req.Pagination, func(key []byte, value []byte) error {
		var roundData types.RoundData
		if err := k.cdc.Unmarshal(value, &roundData); err != nil {
			return err
		}

		roundDatas = append(roundDatas, roundData)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRoundDataResponse{RoundData: roundDatas, Pagination: pageRes}, nil
}

func (k Keeper) RoundData(goCtx context.Context, req *types.QueryGetRoundDataRequest) (*types.QueryGetRoundDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetRoundData(
		ctx,
		req.TokenId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRoundDataResponse{RoundData: val}, nil
}
