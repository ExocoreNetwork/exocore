package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	delegationtype "github.com/exocore/x/delegation/types"
)

var _ delegationtype.QueryServer = Keeper{}

func (k Keeper) QuerySingleDelegationInfo(ctx context.Context, req *delegationtype.SingleDelegationInfoReq) (*delegationtype.DelegationAmounts, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetSingleDelegationInfo(c, req.StakerId, req.AssetId, req.OperatorAddr)
}

func (k Keeper) QueryDelegationInfo(ctx context.Context, info *delegationtype.DelegationInfoReq) (*delegationtype.QueryDelegationInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetDelegationInfo(c, info.StakerId, info.AssetId)
}

func (k Keeper) QueryOperatorInfo(ctx context.Context, req *delegationtype.QueryOperatorInfoReq) (*delegationtype.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetOperatorInfo(c, req.OperatorAddr)
}
