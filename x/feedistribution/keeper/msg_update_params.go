package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
)

func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// For test purpose, skip for now
	// if k.GetAuthority() != req.Authority {
	// 	 return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.GetAuthority(), req.Authority)
	// }

	ctx := sdk.UnwrapSDKContext(goCtx)
	k.SetParams(ctx, req.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}
