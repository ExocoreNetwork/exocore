package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/avs/keeper"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
)

func (suite *AVSTestSuite) Test_GroupStatistics() {
	tasks := []avstypes.TaskResultInfo{
		{OperatorAddress: "addr1", TaskResponseHash: "hash1", TaskResponse: []byte("response1"), BlsSignature: []byte("sig1"), TaskContractAddress: "contract1", TaskId: 1},
		{OperatorAddress: "addr2", TaskResponseHash: "hash2", TaskResponse: []byte("response2"), BlsSignature: []byte("sig2"), TaskContractAddress: "contract2", TaskId: 2},
		{OperatorAddress: "addr3", TaskResponseHash: "hash3", TaskResponse: []byte("response3"), BlsSignature: []byte("sig3"), TaskContractAddress: "contract1", TaskId: 1},
	}

	groupedTasks := suite.App.AVSManagerKeeper.GroupTasksByIDAndAddress(tasks)
	suite.Equal(2, len(groupedTasks["contract1_1"]))
}
func (suite *AVSTestSuite) TestEpochEnd_TaskCalculation() {
	suite.TestSubmitTask_OnlyPhaseTwo_Mul()
	err := suite.App.OperatorKeeper.SetAVSUSDValue(suite.Ctx, suite.avsAddr, sdkmath.LegacyNewDec(500))
	for _, operatorAddress := range suite.operatorAddresses {
		delta := operatortypes.DeltaOperatorUSDInfo{
			SelfUSDValue:   sdkmath.LegacyNewDec(100),
			TotalUSDValue:  sdkmath.LegacyNewDec(100),
			ActiveUSDValue: sdkmath.LegacyNewDec(100),
		}
		suite.App.OperatorKeeper.UpdateOperatorUSDValue(suite.Ctx, suite.avsAddr, operatorAddress, delta)
	}

	suite.NoError(err)
	suite.CommitAfter(suite.EpochDuration)

	info, err := suite.App.AVSManagerKeeper.GetTaskInfo(suite.Ctx, strconv.FormatUint(suite.taskId, 10), common.Address(suite.taskAddress.Bytes()).String())
	suite.NoError(err)
	expectInfo := &avstypes.TaskInfo{
		TaskContractAddress:   suite.taskAddress.String(),
		Name:                  "test-avsTask",
		TaskId:                suite.taskId,
		Hash:                  []byte("req-struct"),
		TaskResponsePeriod:    2,
		TaskStatisticalPeriod: 1,
		TaskChallengePeriod:   2,
		ThresholdPercentage:   0,
		StartingEpoch:         20,
		OptInOperators:        suite.operatorAddresses,
		NoSignedOperators:     nil,
		SignedOperators:       suite.operatorAddresses,
		ActualThreshold:       0,
	}
	diff := keeper.Difference(expectInfo.SignedOperators, info.SignedOperators)

	suite.Equal(0, len(diff))
	suite.Equal(expectInfo.NoSignedOperators, info.NoSignedOperators)
	suite.Equal(expectInfo.ActualThreshold, info.ActualThreshold)

}
