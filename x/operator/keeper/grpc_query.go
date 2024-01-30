package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	operatortypes "github.com/exocore/x/operator/types"
)

var _ operatortypes.QueryServer = &Keeper{}

func (k *Keeper) GetOperatorInfo(ctx context.Context, req *operatortypes.GetOperatorInfoReq) (*operatortypes.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.OperatorInfo(c, req.OperatorAddr)
}
