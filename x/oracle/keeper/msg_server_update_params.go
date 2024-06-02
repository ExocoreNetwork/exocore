package keeper

import (
	"context"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ms msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: skip the authority check for test
	//	if ms.authority != msg.Authority {
	//		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	//	}
	p := ms.GetParams(ctx)
	var err error
	height := uint64(ctx.BlockHeight())
	// add sources
	//	if len(msg.Params.Sources) > 1 {
	// if p, err = p.AddSources(msg.Params.Sources[1:]...); err != nil {
	if p, err = p.AddSources(msg.Params.Sources...); err != nil {
		return nil, err
	}
	//	}
	// add chains
	//	if len(msg.Params.Chains) > 1 {
	if p, err = p.AddChains(msg.Params.Chains...); err != nil {
		return nil, err
	}
	//	}
	// add tokens
	//	if len(msg.Params.Tokens) > 1 {
	if p, err = p.UpdateTokens(msg.Params.Tokens...); err != nil {
		return nil, err
	}
	//	}
	// add rules
	//	if len(msg.Params.Rules) > 1 {
	if p, err = p.AddRules(msg.Params.Rules...); err != nil {
		return nil, err
	}
	//	}
	// udpate tokenFeeders
	// if len(msg.Params.TokenFeeders) > 1 {
	// for _, tokenFeeder := range msg.Params.TokenFeeders[1:] {
	for _, tokenFeeder := range msg.Params.TokenFeeders {
		if p, err = p.UpdateTokenFeeder(tokenFeeder, height); err != nil {
			return nil, err
		}
	}
	// }
	// validate params
	if err = p.Validate(); err != nil {
		return nil, err
	}
	// set updated new params
	ms.SetParams(ctx, p)
	return &types.MsgUpdateParamsResponse{}, nil
}
