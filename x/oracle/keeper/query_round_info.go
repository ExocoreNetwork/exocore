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

func (k Keeper) RoundInfoAll(goCtx context.Context, req *types.QueryAllRoundInfoRequest) (*types.QueryAllRoundInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var roundInfos []types.RoundInfo
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	roundInfoStore := prefix.NewStore(store, types.KeyPrefix(types.RoundInfoKeyPrefix))

	pageRes, err := query.Paginate(roundInfoStore, req.Pagination, func(key []byte, value []byte) error {
		var roundInfo types.RoundInfo
		if err := k.cdc.Unmarshal(value, &roundInfo); err != nil {
			return err
		}

		roundInfos = append(roundInfos, roundInfo)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRoundInfoResponse{RoundInfo: roundInfos, Pagination: pageRes}, nil
}

func (k Keeper) RoundInfo(goCtx context.Context, req *types.QueryGetRoundInfoRequest) (*types.QueryGetRoundInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetRoundInfo(
		ctx,
		req.TokenId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRoundInfoResponse{RoundInfo: val}, nil
}
