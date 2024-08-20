package keeper

import (
	"fmt"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EpochsHooksWrapper is the wrapper structure that implements the epochs hooks for the avs
// keeper.
type EpochsHooksWrapper struct {
	keeper *Keeper
}

// Interface guard
var _ epochstypes.EpochHooks = EpochsHooksWrapper{}

// EpochsHooks returns the epochs hooks wrapper. It follows the "accept interfaces, return
// concretes" pattern.
func (k *Keeper) EpochsHooks() EpochsHooksWrapper {
	return EpochsHooksWrapper{k}
}

// AfterEpochEnd is called after an epoch ends. It is called during the BeginBlock function.
func (wrapper EpochsHooksWrapper) AfterEpochEnd(
	ctx sdk.Context, epochIdentifier string, epochNumber int64,
) {
	// get all the task info bypass the epoch end
	// threshold calculation, signature verification, nosig quantity statistics
	// todo: need to consider the calling order
	taskList := wrapper.keeper.GetTaskStatisticalEpochEndAVSs(ctx, epochIdentifier, epochNumber)
	for _, task := range taskList {
		fmt.Print(task)
		// task address should be hex
		//err := wrapper.keeper.UpdateVotingPower(ctx, task)
		//if err != nil {
		//	ctx.Logger().Error("Failed to update voting power", "task", task, "error", err)
		//	// Handle the error gracefully, continue to the next AVS
		//	continue
		//}
	}
}

// BeforeEpochStart is called before an epoch starts.
func (wrapper EpochsHooksWrapper) BeforeEpochStart(
	sdk.Context, string, int64,
) {
}
