package keeper_test

import tasktype "github.com/ExocoreNetwork/exocore/x/avstask/types"

func (suite *AvsTaskTestSuite) TestTaskInfo() {
	info := &tasktype.RegisterAVSTaskReq{
		FromAddress: suite.AccAddress.String(),
		Task: &tasktype.TaskContractInfo{
			Name:                "test-avstask-01",
			MetaInfo:            "avstask up",
			TaskContractAddress: "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8",
			Status:              "active",
		},
	}
	err := suite.App.TaskKeeper.SetAVSTaskInfo(suite.Ctx, info)
	suite.NoError(err)

	getTaskInfo, err := suite.App.TaskKeeper.GetAVSTaskInfo(suite.Ctx, "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8")
	suite.NoError(err)
	suite.Equal(*info.Task, *getTaskInfo)
}

func (suite *AvsTaskTestSuite) TestOperator_pubkey() {
	err := suite.App.TaskKeeper.SetOperatorPubKey(suite.Ctx, "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8", []byte("pubkey"))
	suite.NoError(err)

	pub2, err := suite.App.TaskKeeper.GetOperatorPubKey(suite.Ctx, "exo1j9ly7f0jynscjgvct0enevaa659te58k3xztc8")
	suite.NoError(err)
	suite.Equal([]byte("pubkey"), pub2)
}
