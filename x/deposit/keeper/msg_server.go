package keeper

import (
	"context"

	deposittype "github.com/ExocoreNetwork/exocore/x/deposit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ deposittype.MsgServer = &Keeper{}

// UpdateParams set `exoCoreLzAppAddress` in the parameters of the deposit module, it can be used to verify whether the caller of precompile contracts is the `exoCoreLzApp` contract.
// This function should be triggered by the governance in the future,and we need to move this function to the `restaking_assets_manage` module to facilitate the query by other modules.
func (k Keeper) UpdateParams(ctx context.Context, params *deposittype.MsgUpdateParams) (*deposittype.MsgUpdateParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	err := k.SetParams(c, &params.Params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
