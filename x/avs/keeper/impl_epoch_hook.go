package keeper

import (
	"strconv"

	sdkmath "cosmossdk.io/math"
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
	taskResList := wrapper.keeper.GetTaskStatisticalEpochEndAVSs(ctx, epochIdentifier, epochNumber)

	if len(taskResList) != 0 {
		groupedTasks := wrapper.keeper.GroupTasksByIDAndAddress(taskResList)
		for _, value := range groupedTasks {
			var signedOperatorList []string
			var taskID uint64
			var taskAddr string
			var avsAddr string
			totalPower := sdkmath.LegacyNewDec(0)
			for _, res := range value {
				// Find signed operators
				if res.BlsSignature != nil {
					signedOperatorList = append(signedOperatorList, res.OperatorAddress)
					if avsAddr == "" {
						avsInfo := wrapper.keeper.GetAVSInfoByTaskAddress(ctx, res.TaskContractAddress)
						avsAddr = avsInfo.AvsAddress
					}
					if taskID == 0 {
						taskID = res.TaskId
					}
					if taskAddr == "" {
						taskAddr = res.TaskContractAddress
					}
					usd, err := wrapper.keeper.operatorKeeper.GetOperatorOptedUSDValue(ctx, avsAddr, res.OperatorAddress)
					if err != nil {
						ctx.Logger().Error("Failed to update task result statistics", "task result", taskAddr, "error", err)
						// Handle the error gracefully, continue to the next
						continue
					}
					if usd.ActiveUSDValue.IsNegative() {
						ctx.Logger().Error("Failed to update task result statistics", "task result", taskAddr, "error", err)
						// Handle the error gracefully, continue to the next
						continue
					}
					totalPower = totalPower.Add(usd.ActiveUSDValue)
				}
			}
			taskInfo, err := wrapper.keeper.GetTaskInfo(ctx, strconv.FormatUint(taskID, 10), taskAddr)
			if err != nil {
				ctx.Logger().Error("Failed to update task result statistics", "task result", taskAddr, "error", err)
				// Handle the error gracefully, continue to the next
				continue
			}
			diff := Difference(taskInfo.OptInOperators, signedOperatorList)
			taskInfo.SignedOperators = signedOperatorList
			taskInfo.NoSignedOperators = diff

			// Calculate actual threshold

			actualThreshold := wrapper.keeper.CalculateActualThreshold(ctx, totalPower, avsAddr)

			taskInfo.ActualThreshold = actualThreshold.BigInt().Uint64()

			// Update the taskInfo in the state
			err = wrapper.keeper.SetTaskInfo(ctx, taskInfo)
			if err != nil {
				ctx.Logger().Error("Failed to update task result statistics", "task result", taskAddr, "error", err)
				// Handle the error gracefully, continue to the next
				continue
			}
		}
	}
}

// BeforeEpochStart is called before an epoch starts.
func (wrapper EpochsHooksWrapper) BeforeEpochStart(
	sdk.Context, string, int64,
) {
}
