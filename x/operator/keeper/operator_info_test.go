package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	operatortype "github.com/ExocoreNetwork/exocore/x/operator/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *OperatorTestSuite) TestOperatorInfo() {
	info := &operatortype.OperatorInfo{
		EarningsAddr:     suite.AccAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &operatortype.ClientChainEarningAddrList{
			EarningInfoList: []*operatortype.ClientChainEarningAddrInfo{
				{101, "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
			},
		},
		Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
		MinSelfDelegation: math.NewInt(0),
	}
	err := suite.App.OperatorKeeper.SetOperatorInfo(suite.Ctx, suite.AccAddress.String(), info)
	suite.NoError(err)

	getOperatorInfo, err := suite.App.OperatorKeeper.GetOperatorInfo(suite.Ctx, &operatortype.GetOperatorInfoReq{OperatorAddr: suite.AccAddress.String()})
	suite.NoError(err)
	suite.Equal(*info, *getOperatorInfo)
}

func (suite *OperatorTestSuite) TestHistoricalOperatorInfo() {
	height := suite.Ctx.BlockHeight()
	info := &operatortype.OperatorInfo{
		EarningsAddr:     suite.AccAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &operatortype.ClientChainEarningAddrList{
			EarningInfoList: nil,
		},
	}
	err := suite.App.OperatorKeeper.SetOperatorInfo(suite.Ctx, suite.AccAddress.String(), info)
	suite.NoError(err)
	suite.NextBlock()
	suite.Equal(height+1, suite.Ctx.BlockHeight(), "nexBlock failed")

	newInfo := *info
	newInfo.OperatorMetaInfo = "new operator"
	err = suite.App.OperatorKeeper.SetOperatorInfo(suite.Ctx, suite.AccAddress.String(), &newInfo)
	suite.NoError(err)

	// get historical operator info
	historicalQueryCtx, err := types.ContextForHistoricalState(suite.Ctx, height)
	suite.NoError(err)
	getInfo, err := suite.App.OperatorKeeper.GetOperatorInfo(historicalQueryCtx, &operatortype.GetOperatorInfoReq{
		OperatorAddr: suite.AccAddress.String(),
	})
	suite.NoError(err)
	suite.Equal(info.OperatorMetaInfo, getInfo.OperatorMetaInfo)

	getInfo, err = suite.App.OperatorKeeper.GetOperatorInfo(suite.Ctx, &operatortype.GetOperatorInfoReq{
		OperatorAddr: suite.AccAddress.String(),
	})
	suite.NoError(err)
	suite.Equal(newInfo.OperatorMetaInfo, getInfo.OperatorMetaInfo)
}
