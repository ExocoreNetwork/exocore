package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/slash/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint: unused // Implementation of the msgServer (via proto) to be done.
type msgServer struct {
	Keeper
}

func (k Keeper) UpdateParams(ctx context.Context, params *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetParams(c, &params.Params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
