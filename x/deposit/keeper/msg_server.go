package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	deposittype "github.com/exocore/x/deposit/types"
)

var _ deposittype.MsgServer = &Keeper{}

func (k Keeper) UpdateParams(ctx context.Context, params *deposittype.MsgUpdateParams) (*deposittype.MsgUpdateParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetParams(c, &params.Params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
