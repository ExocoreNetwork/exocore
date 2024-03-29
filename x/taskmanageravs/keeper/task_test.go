package keeper_test

import tasktype "github.com/ExocoreNetwork/exocore/x/taskmanageravs/types"

func (suite *KeeperTestSuite) TestTaskInfo() {
	info := &tasktype.RegisterAVSTaskReq{
		AVSAddress: suite.accAddress.String(),
		Task: &tasktype.TaskContractInfo{
			Name:                "test-task-01",
			MetaInfo:            "task up",
			TaskContractAddress: "0x0000000000000000000000000000000000000901",
			TaskContractId:      99,
			Status:              "active",
		},
	}
	index, err := suite.app.TaskKeeper.SetAvsTaskInfo(suite.ctx, info)
	suite.NoError(err)

	getOperatorInfo, err := suite.app.TaskKeeper.GetAvsTaskInfo(suite.ctx, index)
	suite.NoError(err)
	suite.Equal(*info.Task, *getOperatorInfo)
}
