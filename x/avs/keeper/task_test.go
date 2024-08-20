package keeper_test

import (
	types "github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *AVSTestSuite) TestTaskInfo() {
	info := &types.TaskInfo{
		TaskContractAddress: common.Address(suite.AccAddress.Bytes()).String(),
		Name:                "test-avstask-01",
		TaskId:              3,
		Hash:                []byte("active"),
		TaskResponsePeriod:  10000,
		TaskChallengePeriod: 5000,
		ThresholdPercentage: 60,
	}
	err := suite.App.AVSManagerKeeper.SetTaskInfo(suite.Ctx, info)
	suite.NoError(err)

	getTaskInfo, err := suite.App.AVSManagerKeeper.GetTaskInfo(suite.Ctx, "avstask01", common.Address(suite.AccAddress.Bytes()).String())
	suite.NoError(err)
	suite.Equal(*info, *getTaskInfo)
}

func (suite *AVSTestSuite) TestOperator_pubkey() {
	blsPub := &types.BlsPubKeyInfo{
		Operator: "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8",
		PubKey:   []byte("pubkey"),
		Name:     "pubkey",
	}

	err := suite.App.AVSManagerKeeper.SetOperatorPubKey(suite.Ctx, blsPub)
	suite.NoError(err)

	pub, err := suite.App.AVSManagerKeeper.GetOperatorPubKey(suite.Ctx, "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8")
	suite.NoError(err)
	suite.Equal([]byte("pubkey"), pub.PubKey)
}
func (suite *AVSTestSuite) TestGetTaskId() {
	addr := common.Address(suite.AccAddress.Bytes())

	taskId := suite.App.AVSManagerKeeper.GetTaskID(suite.Ctx, addr)
	suite.Equal(uint64(1), taskId)

	taskId = suite.App.AVSManagerKeeper.GetTaskID(suite.Ctx, addr)
	suite.Equal(uint64(2), taskId)
	taskId = suite.App.AVSManagerKeeper.GetTaskID(suite.Ctx, addr)
	suite.Equal(uint64(3), taskId)

	addr = common.Address(suite.avsAddress.Bytes())

	taskId = suite.App.AVSManagerKeeper.GetTaskID(suite.Ctx, addr)
	suite.Equal(uint64(1), taskId)

	taskId = suite.App.AVSManagerKeeper.GetTaskID(suite.Ctx, addr)
	suite.Equal(uint64(2), taskId)
}
