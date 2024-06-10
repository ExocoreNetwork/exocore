package keeper

import (
	"strings"

	types "github.com/ExocoreNetwork/exocore/x/exomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"
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
	ctx sdk.Context, identifier string, epoch int64,
) {
	params := wrapper.keeper.GetParams(ctx)
	if strings.Compare(identifier, params.EpochIdentifier) == 0 {
		mintedCoin := sdk.NewCoin(
			params.MintDenom,
			params.EpochReward,
		)
		mintedCoins := sdk.NewCoins(mintedCoin)

		err := wrapper.keeper.MintCoins(ctx, mintedCoins)
		if err != nil {
			ctx.Logger().With(types.ModuleName).Error(
				"AfterEpochEnd",
				"could not mint coins", err,
			)
			return
		}

		err = wrapper.keeper.AddCollectedFees(ctx, mintedCoins)
		if err != nil {
			ctx.Logger().With(types.ModuleName).Error(
				"AfterEpochEnd",
				"could not transfer coins", err,
			)
			return
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeMint,
				sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
			),
		)

		ctx.Logger().With(types.ModuleName).Info(
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
