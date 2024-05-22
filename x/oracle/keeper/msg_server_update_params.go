package keeper

import (
	"context"
	"errors"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ms msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: skip the authority check for test
	//	if ms.authority != msg.Authority {
	//		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	//	}

	// store params
	//
	//	ms.SetParams(ctx, msg.Params)
	// validation check on tokenfeeders
	p := ms.GetParams(ctx)
	height := uint64(ctx.BlockHeight())
	for _, feeder := range msg.Params.TokenFeeders {
		if feeder.StartBaseBlock > 0 && feeder.StartBaseBlock <= height {
			return nil, errors.New("startBaseBlock for tokenFeeder invalid: history block")
		}
		if fIDs := p.GetFeederIDsByTokenID(feeder.TokenID); len(fIDs) > 0 {
			// update exist tokenFeeder(startBlock, endBlock)
			// the latest feeder: 1. not start, 2. stopped, 3. running

			f := p.TokenFeeders[fIDs[len(fIDs)-1]]

			// latest feeder for this token has stopped
			if f.EndBlock <= height {
				// this should be a feeder set continue from the last one
				// startroundid valid
				lastRoundID := (f.EndBlock-f.StartBaseBlock)/f.Interval + f.StartRoundID
				if feeder.StartRoundID != lastRoundID+1 {
					return nil, errors.New("startRoundID for tokenFeeder invalid: should be last roundid+1 for this token")
				}
				if feeder.StartBaseBlock < 1 || feeder.Interval < 1 {
					return nil, errors.New("interval/startBaseblock for tokenFeeder invalid: should be bigger than 1")
				}
				if feeder.RuleID >= uint64(len(p.Rules)) {
					return nil, errors.New("ruleID for tokenFeeder invalid: rule doesn't exist")
				}
				// add a new feeder, restart a feeder for one token
				p.TokenFeeders = append(p.TokenFeeders, feeder)
				continue
			}

			// latest feeder is running
			if f.StartBaseBlock <= height {
				// update can only be used to set the EndBlock
				if feeder.RuleID != 0 || feeder.StartRoundID != 0 || feeder.StartBaseBlock != 0 || feeder.Interval != 0 {
					return nil, errors.New("fields invalid for tokenFeeder: only EndBlock should could be set when try to update an running feeder")

				}
				if feeder.EndBlock <= f.StartBaseBlock || (feeder.EndBlock-f.StartRoundID)%f.Interval < 3 {
					return nil, errors.New("endBlock for tokenFeeder invalid: when update endblock for exist feeder")
				}
				f.EndBlock = feeder.EndBlock
				continue
			}

			// feeder not start yet
			if f.StartBaseBlock > height {
				// update: startBaseBlock, endBlock
				if feeder.RuleID != 0 || feeder.StartRoundID != 0 || feeder.Interval != 0 {
					return nil, errors.New("fields invalid for tokenFeeder: only EndBlock should could be set when try to update an running feeder")
				}
				if feeder.StartBaseBlock > 0 {
					if feeder.EndBlock > 0 {
						if (feeder.EndBlock-feeder.StartBaseBlock)%f.Interval < 3 {
							return nil, errors.New("update startBaseBlock&endBlock for tokenFeeder invalid")
						}
						f.EndBlock = feeder.EndBlock
					}
					f.StartBaseBlock = feeder.StartBaseBlock
				} else {
					if feeder.EndBlock > 0 {
						if (feeder.EndBlock-feeder.StartBaseBlock)%f.Interval < 3 {
							return nil, errors.New("update startBaseBlock&endBlock for tokenFeeder invalid: endBlock invalid")

						}
						f.EndBlock = feeder.EndBlock
						continue
					} else {
						return nil, errors.New("update startBaseBlock&endBlock for tokenFeeder invalid: at least one field is set")
					}
				}
				continue
			}
		} else {
			// create feeder for a token (first time for this token)
			if feeder.TokenID >= uint64(len(p.Tokens)) || feeder.RuleID >= uint64(len(p.Rules)) || feeder.StartBaseBlock <= height || feeder.Interval < 1 || feeder.StartRoundID != 1 {
				return nil, errors.New("invalid input for create new feeder for a token the first time")
			}
			p.TokenFeeders = append(p.TokenFeeders, feeder)
		}
	}

	// TODO: validation check on chains, tokens, rules, sources, and cross verification
	// 	ms.SetUpdateParams(ctx, msg.Params)
	ms.SetParams(ctx, msg.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}
