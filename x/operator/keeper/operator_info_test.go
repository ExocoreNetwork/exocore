package keeper_test

import (
	operatortype "github.com/exocore/x/operator/types"
)

func (s *KeeperTestSuite) TestOperatorInfo() {
	info := &operatortype.OperatorInfo{
		EarningsAddr:     s.accAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &operatortype.ClientChainEarningAddrList{
			EarningInfoList: []*operatortype.ClientChainEarningAddrInfo{
				{101, "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
			},
		},
	}
	err := s.app.OperatorKeeper.SetOperatorInfo(s.ctx, s.accAddress.String(), info)
	s.NoError(err)

	getOperatorInfo, err := s.app.OperatorKeeper.GetOperatorInfo(s.ctx, &operatortype.GetOperatorInfoReq{OperatorAddr: s.accAddress.String()})
	s.NoError(err)
	s.Equal(*info, *getOperatorInfo)
}

func (s *KeeperTestSuite) TestHistoricalOperatorInfo() {
	height := s.ctx.BlockHeight()
	info := &operatortype.OperatorInfo{
		EarningsAddr:     s.accAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &operatortype.ClientChainEarningAddrList{
			EarningInfoList: nil,
		},
	}
	err := s.app.OperatorKeeper.SetOperatorInfo(s.ctx, s.accAddress.String(), info)
	s.NoError(err)
	s.NextBlock()
	s.Equal(height+1, s.ctx.BlockHeight(), "nexBlock failed")

	newInfo := *info
	newInfo.OperatorMetaInfo = "new operator"
	err = s.app.OperatorKeeper.SetOperatorInfo(s.ctx, s.accAddress.String(), &newInfo)
	s.NoError(err)

	//get historical operator info
	historicalQueryCtx, err := s.app.CreateQueryContext(height, false)
	s.NoError(err)
	getInfo, err := s.app.OperatorKeeper.GetOperatorInfo(historicalQueryCtx, &operatortype.GetOperatorInfoReq{
		OperatorAddr: s.accAddress.String(),
	})
	s.NoError(err)
	s.Equal(info.OperatorMetaInfo, getInfo.OperatorMetaInfo)

	getInfo, err = s.app.OperatorKeeper.GetOperatorInfo(s.ctx, &operatortype.GetOperatorInfoReq{
		OperatorAddr: s.accAddress.String(),
	})
	s.NoError(err)
	s.Equal(getInfo.OperatorMetaInfo, newInfo.OperatorMetaInfo)
}
