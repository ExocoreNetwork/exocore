package keeper_test

import (
	"fmt"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func (suite *AVSTestSuite) TestSubmitTask1() {
	suite.prepare()
	taskRes := avstypes.TaskResponse{TaskID: 1, NumberSum: big.NewInt(100)}
	jsonData, err := avstypes.MarshalTaskResponse(taskRes)
	suite.NoError(err)
	_ = crypto.Keccak256Hash(jsonData)

	// pub, err := suite.App.AVSManagerKeeper.GetOperatorPubKey(suite.Ctx, suite.operatorAddr.String())
	suite.NoError(err)

	msg, _ := avstypes.GetTaskResponseDigest(taskRes)
	msgBytes := msg[:]
	sig := suite.blsKey.Sign(msgBytes)

	info := &avstypes.TaskResultInfo{
		TaskContractAddress: suite.taskAddress.String(),
		OperatorAddress:     suite.operatorAddr.String(),
		TaskId:              suite.taskId,
		TaskResponseHash:    "",
		TaskResponse:        nil,
		BlsSignature:        sig.Marshal(),
		Stage:               avstypes.TwoPhaseCommitOne,
	}
	err = suite.App.AVSManagerKeeper.SetTaskResultInfo(suite.Ctx, suite.operatorAddr.String(), info)
	suite.NoError(err)

}

func (suite *AVSTestSuite) TestAVSUSDValue1() {
	suite.prepare()
	avsUSDValue, err := suite.App.OperatorKeeper.GetAVSUSDValue(suite.Ctx, suite.avsAddr)
	suite.NoError(err)
	optedUSDValues, err := suite.App.OperatorKeeper.GetOperatorOptedUSDValue(suite.Ctx, suite.avsAddr, suite.operatorAddr.String())
	suite.NoError(err)
	fmt.Println(avsUSDValue)
	fmt.Println(optedUSDValues)

}
