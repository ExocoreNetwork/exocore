package keeper_test

import delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"

func (suite *KeeperTestSuite) TestOperatorInfo() {
	info := &delegationtype.OperatorInfo{
		EarningsAddr:     suite.accAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &delegationtype.ClientChainEarningAddrList{
			EarningInfoList: []*delegationtype.ClientChainEarningAddrInfo{
				{101, "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
			},
		},
	}
	err := suite.app.DelegationKeeper.SetOperatorInfo(suite.ctx, suite.accAddress.String(), info)
	suite.NoError(err)

	getOperatorInfo, err := suite.app.DelegationKeeper.GetOperatorInfo(suite.ctx, suite.accAddress.String())
	suite.NoError(err)
	suite.Equal(*info, *getOperatorInfo)
}
