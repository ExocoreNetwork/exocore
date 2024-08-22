package keeper_test

import (
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
)

func (suite *AVSTestSuite) TestSubmitTask_PhaseAll1() {
	suite.TestSubmitTask_OnlyPhaseTwo_Mul()
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
		ThresholdPercentage:   60,
		StartingEpoch:         20,
		OptInOperators:        suite.operatorAddresses,
		NoSignedOperators:     nil,
		SignedOperators:       nil,
		ActualThreshold:       0,
	}
	suite.Equal(*expectInfo, *info)
}
