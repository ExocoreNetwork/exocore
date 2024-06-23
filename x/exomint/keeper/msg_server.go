package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/x/exomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateParams is used to trigger a params update. It must be signed by the authority.
func (k Keeper) UpdateParams(
	ctx context.Context, msg *types.MsgUpdateParams,
) (*types.MsgUpdateParamsResponse, error) {
	c := sdk.UnwrapSDKContext(ctx)
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			k.authority, msg.Authority,
		)
	}
	prevParams := k.GetParams(c)
	nextParams := msg.Params
	// stateless validations
	overParams := nextParams.OverrideIfUnset(prevParams, k.Logger(c))
	// stateful validations
	// no need to check if MintDenom is registered in BankKeeper, since it does not itself perform such checks.
	// the reward is already guaranteed to be positive and fits in the bit length.
	// so, we just have to check epoch here.
	if _, found := k.epochsKeeper.GetEpochInfo(c, overParams.EpochIdentifier); !found {
		k.Logger(c).Info("UpdateParams", "overriding EpochIdentifier with value", prevParams.EpochIdentifier)
		overParams.EpochIdentifier = prevParams.EpochIdentifier
	}
	k.SetParams(c, overParams)
	return &types.MsgUpdateParamsResponse{}, nil
}
