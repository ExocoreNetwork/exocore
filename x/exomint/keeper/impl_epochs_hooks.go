package keeper

import (
	"fmt"
	"strings"

	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	types "github.com/ExocoreNetwork/exocore/x/exomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EpochsHooksWrapper is the wrapper structure that implements the epochs hooks for the
// keeper.
type EpochsHooksWrapper struct {
	keeper *Keeper
}

// Interface guard
var _ epochstypes.EpochHooks = EpochsHooksWrapper{}

// EpochsHooks returns the epochs hooks wrapper.
func (k *Keeper) EpochsHooks() EpochsHooksWrapper {
	return EpochsHooksWrapper{k}
}

// AfterEpochEnd is called after an epoch ends. It is called during the BeginBlock function.
func (wrapper EpochsHooksWrapper) AfterEpochEnd(
	ctx sdk.Context, identifier string, number int64,
) {
	params := wrapper.keeper.GetParams(ctx)
	if strings.Compare(identifier, params.EpochIdentifier) == 0 {
		mintedCoin := sdk.NewCoin(
			params.MintDenom,
			params.EpochReward,
		)
		mintedCoins := sdk.NewCoins(mintedCoin)

		logger := wrapper.keeper.Logger(ctx)
		err := wrapper.keeper.MintCoins(ctx, mintedCoins)
		if err != nil {
			logger.Error(
				"AfterEpochEnd",
				"could not mint coins", err,
			)
			return
		}

		err = wrapper.keeper.AddCollectedFees(ctx, mintedCoins)
		if err != nil {
			logger.Error(
				"AfterEpochEnd",
				"could not transfer coins", err,
			)
			return
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMint,
				sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
				sdk.NewAttribute(types.AttributeEpochIdentifier, identifier),
				sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprintf("%d", number)),
			),
		)

		logger.Info(
			"AfterEpochEnd",
			"minted successfully", mintedCoins.String(),
		)
	}
}

// BeforeEpochStart is called before an epoch starts.
func (wrapper EpochsHooksWrapper) BeforeEpochStart(
	sdk.Context, string, int64,
) {
	// no-op
}
