package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	operatortypes "github.com/exocore/x/operator/types"
)

var _ operatortypes.QueryServer = Keeper{}

func (k Keeper) QueryOperatorInfo(ctx context.Context, req *operatortypes.QueryOperatorInfoReq) (*operatortypes.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.GetOperatorInfo(c, req.OperatorAddr)
}
