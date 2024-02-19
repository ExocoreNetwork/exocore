package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/withdraw/types"
)

// nolint: unused // To be implemented when creating the requests.
type msgServer struct {
	Keeper
}

func (k Keeper) UpdateParams(context.Context, *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// c := sdk.UnwrapSDKContext(ctx)
	// err := k.SetParams(c, &params.Params)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
