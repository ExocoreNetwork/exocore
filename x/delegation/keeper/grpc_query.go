package keeper

import (
	"context"

	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ delegationtype.QueryServer = Keeper{}

func (k Keeper) QuerySingleDelegationInfo(ctx context.Context, req *delegationtype.SingleDelegationInfoReq) (*delegationtype.DelegationAmounts, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetSingleDelegationInfo(c, req.StakerID, req.AssetID, req.OperatorAddr)
}

func (k Keeper) QueryDelegationInfo(ctx context.Context, info *delegationtype.DelegationInfoReq) (*delegationtype.QueryDelegationInfoResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetDelegationInfo(c, info.StakerID, info.AssetID)
}

func (k Keeper) QueryOperatorInfo(ctx context.Context, req *delegationtype.QueryOperatorInfoReq) (*delegationtype.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetOperatorInfo(c, req.OperatorAddr)
}
