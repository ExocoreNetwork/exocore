package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
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
	prevParams := k.GetDogfoodParams(c)
	nextParams := msg.Params
	logger := k.Logger(c)
	if nextParams.EpochsUntilUnbonded == 0 {
		logger.Info(
			"UpdateParams",
			"overriding EpochsUntilUnbonded with value", prevParams.EpochsUntilUnbonded,
		)
		nextParams.EpochsUntilUnbonded = prevParams.EpochsUntilUnbonded
	}
	if nextParams.MaxValidators == 0 {
		logger.Info(
			"UpdateParams",
			"overriding MaxValidators with value", prevParams.MaxValidators,
		)
		nextParams.MaxValidators = prevParams.MaxValidators
	}
	if err := epochstypes.ValidateEpochIdentifierInterface(
		nextParams.EpochIdentifier,
	); err != nil {
		logger.Info(
			"UpdateParams",
			"overriding EpochIdentifier with value", prevParams.EpochIdentifier,
		)
		nextParams.EpochIdentifier = prevParams.EpochIdentifier
	}
	if nextParams.HistoricalEntries == 0 {
		logger.Info(
			"UpdateParams",
			"overriding HistoricalEntries with value", prevParams.HistoricalEntries,
		)
		nextParams.HistoricalEntries = prevParams.HistoricalEntries
	}
	if len(nextParams.AssetIDs) == 0 {
		logger.Error(
			"UpdateParams",
			"overriding AssetIDs with value", prevParams.AssetIDs,
		)
		nextParams.AssetIDs = prevParams.AssetIDs
	}
	if nextParams.MinSelfDelegation.IsNil() || nextParams.MinSelfDelegation.IsNegative() {
		logger.Error(
			"UpdateParams",
			"overriding MinSelfDelegation with value", prevParams.MinSelfDelegation,
		)
		nextParams.MinSelfDelegation = prevParams.MinSelfDelegation
	}
	// now do stateful validations
	if _, found := k.epochsKeeper.GetEpochInfo(c, nextParams.EpochIdentifier); !found {
		logger.Info(
			"UpdateParams",
			"overriding EpochIdentifier with value", prevParams.EpochIdentifier,
		)
		nextParams.EpochIdentifier = prevParams.EpochIdentifier
	}
	override := false
	for _, assetID := range nextParams.AssetIDs {
		if !k.restakingKeeper.IsStakingAsset(c, strings.ToLower(assetID)) {
			override = true
			break
		}
	}
	if override {
		logger.Info(
			"UpdateParams",
			"overriding AssetIDs with value", prevParams.AssetIDs,
		)
		nextParams.AssetIDs = prevParams.AssetIDs
	}
	k.SetParams(c, nextParams)
	return &types.MsgUpdateParamsResponse{}, nil
}
