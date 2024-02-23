package keeper

import (
	"context"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ operatortypes.QueryServer = &Keeper{}

func (k *Keeper) GetOperatorInfo(ctx context.Context, req *operatortypes.GetOperatorInfoReq) (*operatortypes.OperatorInfo, error) {
	c := sdk.UnwrapSDKContext(ctx)
	return k.OperatorInfo(c, req.OperatorAddr)
}
