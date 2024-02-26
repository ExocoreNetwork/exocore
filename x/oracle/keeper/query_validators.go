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

func (k Keeper) ValidatorsAll(goCtx context.Context, req *types.QueryAllValidatorsRequest) (*types.QueryAllValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var validatorss []types.Validators
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	validatorsStore := prefix.NewStore(store, types.KeyPrefix(types.ValidatorsKeyPrefix))

	pageRes, err := query.Paginate(validatorsStore, req.Pagination, func(key []byte, value []byte) error {
		var validators types.Validators
		if err := k.cdc.Unmarshal(value, &validators); err != nil {
			return err
		}

		validatorss = append(validatorss, validators)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllValidatorsResponse{Validators: validatorss, Pagination: pageRes}, nil
}

func (k Keeper) Validators(goCtx context.Context, req *types.QueryGetValidatorsRequest) (*types.QueryGetValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetValidators(
		ctx,
		req.Block,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetValidatorsResponse{Validators: val}, nil
}
