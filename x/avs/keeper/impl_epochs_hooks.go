package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
