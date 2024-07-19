package keeper_test

import tasktype "github.com/ExocoreNetwork/exocore/x/avs/types"

func (suite *AVSTestSuite) TestTaskInfo() {
	info := &tasktype.RegisterAVSTaskReq{
		FromAddress: suite.AccAddress.String(),
		Task: &tasktype.TaskInfo{
			Name:                "test-avstask-01",
			TaskId:              "avstask up",
			TaskContractAddress: "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8",
			Data:                []byte("active"),
		},
	}
	err := suite.App.AVSManagerKeeper.SetAVSTaskInfo(suite.Ctx, info)
	suite.NoError(err)

	getTaskInfo, err := suite.App.AVSManagerKeeper.GetAVSTaskInfo(suite.Ctx, "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8")
	suite.NoError(err)
	suite.Equal(*info.Task, *getTaskInfo)
}

func (suite *AVSTestSuite) TestOperator_pubkey() {
	err := suite.App.AVSManagerKeeper.SetOperatorPubKey(suite.Ctx, "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8", []byte("pubkey"))
	suite.NoError(err)

	pub2, err := suite.App.AVSManagerKeeper.GetOperatorPubKey(suite.Ctx, "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8")
	suite.NoError(err)
	suite.Equal([]byte("pubkey"), pub2)
}
