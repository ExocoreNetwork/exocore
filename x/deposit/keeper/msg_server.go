package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/exocore/x/deposit/types"
)

var _ types2.MsgServer = &Keeper{}

func (k Keeper) UpdateParams(ctx context.Context, params *types2.MsgUpdateParams) (*types2.MsgUpdateParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetParams(c, &params.Params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
