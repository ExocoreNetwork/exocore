package keeper

import (
	"context"

	deposittype "github.com/ExocoreNetwork/exocore/x/deposit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ deposittype.MsgServer = &Keeper{}

// UpdateParams This function should be triggered by the governance in the future
func (k Keeper) UpdateParams(ctx context.Context, params *deposittype.MsgUpdateParams) (*deposittype.MsgUpdateParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetParams(c, &params.Params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
