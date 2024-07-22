package keeper_test

import types "github.com/ExocoreNetwork/exocore/x/avs/types"

func (suite *AVSTestSuite) TestTaskInfo() {
	info := &types.TaskInfo{
		TaskContractAddress: suite.AccAddress.String(),
		Name:                "test-avstask-01",
		TaskId:              "avstask01",
		Data:                []byte("active"),
		TaskResponsePeriod:  10000,
		TaskChallengePeriod: 5000,
		ThresholdPercentage: 60,
	}
	err := suite.App.AVSManagerKeeper.SetTaskInfo(suite.Ctx, info)
	suite.NoError(err)

	getTaskInfo, err := suite.App.AVSManagerKeeper.GetTaskInfo(suite.Ctx, "avstask01", suite.AccAddress.String())
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
