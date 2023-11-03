package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/delegation/types"
)

var _ types2.QueryServer = Keeper{}

func (k Keeper) QuerySingleDelegationInfo(ctx context.Context, req *types2.SingleDelegationInfoReq) (*types2.ValueField, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetSingleDelegationInfo(c, req.StakerId, req.AssetId, req.OperatorAddr)
}

func (k Keeper) QueryDelegationInfo(ctx context.Context, info *types2.DelegationInfoReq) (*types2.QueryDelegationInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetDelegationInfo(c, info.StakerId, info.AssetId)
}
func (k Keeper) QueryOperatorInfo(ctx context.Context, req *types2.QueryOperatorInfoReq) (*types2.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetOperatorInfo(c, req.OperatorAddr)
}
