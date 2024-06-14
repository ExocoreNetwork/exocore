package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/avs/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"
)

func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {
}

// AfterEpochEnd Record epoch end avs
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {

	var avsList []types.AVSInfo
	k.IteratEpochEndAVSInfo(ctx, func(_ int64, epochEndAVSInfo types.AVSInfo) (stop bool) {
		avsList = append(avsList, epochEndAVSInfo)
		if epochIdentifier == epochEndAVSInfo.EpochIdentifier {
			return true
		}
		return false
	})

	if epochIdentifier != epochstypes.DayEpochID {
		return
	}

	expEpochID := k.GetEpochIdentifier(ctx)
	if epochIdentifier != expEpochID {
		return
	}

	ctx.EventManager().EmitTypedEvent(
		&types.AVSInfo{},
	)
}
