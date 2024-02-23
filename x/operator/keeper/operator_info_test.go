package keeper_test

import (
	operatortype "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
)

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

	getOperatorInfo, err := suite.app.OperatorKeeper.GetOperatorInfo(suite.ctx, &operatortype.GetOperatorInfoReq{OperatorAddr: suite.accAddress.String()})
	suite.NoError(err)
	suite.Equal(*info, *getOperatorInfo)
}

func (suite *KeeperTestSuite) TestHistoricalOperatorInfo() {
	height := suite.ctx.BlockHeight()
	info := &operatortype.OperatorInfo{
		EarningsAddr:     suite.accAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &operatortype.ClientChainEarningAddrList{
			EarningInfoList: nil,
		},
	}
	err := suite.app.OperatorKeeper.SetOperatorInfo(suite.ctx, suite.accAddress.String(), info)
	suite.NoError(err)
	suite.NextBlock()
	suite.Equal(height+1, suite.ctx.BlockHeight(), "nexBlock failed")

	newInfo := *info
	newInfo.OperatorMetaInfo = "new operator"
	err = suite.app.OperatorKeeper.SetOperatorInfo(suite.ctx, suite.accAddress.String(), &newInfo)
	suite.NoError(err)

	//get historical operator info
	historicalQueryCtx, err := types.ContextForHistoricalState(suite.ctx, height)
	suite.NoError(err)
	getInfo, err := suite.app.OperatorKeeper.GetOperatorInfo(historicalQueryCtx, &operatortype.GetOperatorInfoReq{
		OperatorAddr: suite.accAddress.String(),
	})
	suite.NoError(err)
	suite.Equal(info.OperatorMetaInfo, getInfo.OperatorMetaInfo)

	getInfo, err = suite.app.OperatorKeeper.GetOperatorInfo(suite.ctx, &operatortype.GetOperatorInfoReq{
		OperatorAddr: suite.accAddress.String(),
	})
	suite.NoError(err)
	suite.Equal(newInfo.OperatorMetaInfo, getInfo.OperatorMetaInfo)
}
