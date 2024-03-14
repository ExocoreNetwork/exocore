package keeper_test
delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"

import delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"

func (suite *DelegationTestSuite) TestOperatorInfo() {
	info := &delegationtype.OperatorInfo{
		EarningsAddr:     suite.AccAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &delegationtype.ClientChainEarningAddrList{
			EarningInfoList: []*delegationtype.ClientChainEarningAddrInfo{
				{101, "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
			},
		},
	}
	err := suite.App.DelegationKeeper.SetOperatorInfo(suite.Ctx, suite.AccAddress.String(), info)
	suite.NoError(err)

	getOperatorInfo, err := suite.App.DelegationKeeper.GetOperatorInfo(suite.Ctx, suite.AccAddress.String())
	suite.NoError(err)
	suite.Equal(*info, *getOperatorInfo)
}
