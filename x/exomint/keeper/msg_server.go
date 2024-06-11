package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/x/exomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"
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
	logger := c.Logger().With(types.ModuleName)
	if len(nextParams.MintDenom) == 0 {
		logger.Info("UpdateParams", "overriding MintDenom with value", prevParams.MintDenom)
		nextParams.MintDenom = prevParams.MintDenom
	}
	if nextParams.EpochReward.IsNil() || !nextParams.EpochReward.IsPositive() {
		logger.Info("UpdateParams", "overriding EpochReward with value", prevParams.EpochReward)
		nextParams.EpochReward = prevParams.EpochReward
	}
	if err := epochstypes.ValidateEpochIdentifierInterface(
		nextParams.EpochIdentifier,
	); err != nil {
		logger.Info("UpdateParams", "overriding EpochIdentifier with value", prevParams.EpochIdentifier)
		nextParams.EpochIdentifier = prevParams.EpochIdentifier
	}
	k.SetParams(c, msg.Params)
	return nil, nil
}
