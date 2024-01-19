package keeper_test

import operatortype "github.com/exocore/x/operator/types"

func (suite *KeeperTestSuite) TestOperatorInfo() {
	info := &operatortype.OperatorInfo{
		EarningsAddr:     suite.accAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &operatortype.ClientChainEarningAddrList{
			EarningInfoList: []*operatortype.ClientChainEarningAddrInfo{
				{101, "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
			},
		},
	}
	err := suite.app.OperatorKeeper.SetOperatorInfo(suite.ctx, suite.accAddress.String(), info)
	suite.NoError(err)

	getOperatorInfo, err := suite.app.OperatorKeeper.GetOperatorInfo(suite.ctx, suite.accAddress.String())
	suite.NoError(err)
	suite.Equal(*info, *getOperatorInfo)
}
