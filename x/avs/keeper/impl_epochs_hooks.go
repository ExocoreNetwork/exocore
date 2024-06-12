package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"
	"github.com/evmos/evmos/v14/x/inflation/types"
)

func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {
}

// AfterEpochEnd mints and allocates coins at the end of each epoch end
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {

	if epochIdentifier != epochstypes.DayEpochID {
		return
	}

	expEpochID := k.GetEpochIdentifier(ctx)
	if epochIdentifier != expEpochID {
		return
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeEpochNumber, fmt.Sprintf("%d", epochNumber)),
		),
	)
}
