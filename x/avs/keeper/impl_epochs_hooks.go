package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"
)

// BeforeEpochStart : noop, We don't need to do anything here
func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {
}

// AfterEpochEnd Processing avs epoch end
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	logger := k.Logger(ctx)

	k.IteratAVSInfo(ctx, func(_ int64, avsInfo types.AVSInfo) (stop bool) {
		if epochIdentifier == avsInfo.EpochIdentifier && epochNumber > avsInfo.EffectiveCurrentEpoch {
			{
				logger.Info("Process business logic during avs epoch end", "identifier", epochIdentifier)
				// When the avs epoch ends, execute the business logic such as reward
			}

			return false
		}

		if err := ctx.EventManager().EmitTypedEvent(&types.AVSInfo{}); err != nil {
			logger.Error("Failed to emit event", "error", err)
		}
		return false
	})
}

// Hooks wrapper struct for avs keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Hooks Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// BeforeEpochStart epochs hooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
