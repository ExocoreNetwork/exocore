package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"

	restakingtype "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
)

func (suite *KeeperTestSuite) TestUpdateStakerAssetsState() {
	stakerID := fmt.Sprintf("%s_%s", suite.address, "0")
	ethUniAssetId := fmt.Sprintf("%s_%s", "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984", "101")
	ethUniInitialChangeValue := restakingtype.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: math.NewInt(1000),
		CanWithdrawAmountOrWantChangeValue:  math.NewInt(1000),
	}

	// test the initial storage of statker assets state
	err := suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUniAssetId, ethUniInitialChangeValue)
	suite.Require().NoError(err)

	// test that the retrieved value is correct
	getInfo, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, ethUniAssetId)
	suite.Require().NoError(err)
	suite.Require().True(ethUniInitialChangeValue.TotalDepositAmountOrWantChangeValue.Equal(getInfo.TotalDepositAmountOrWantChangeValue))
	suite.Require().True(ethUniInitialChangeValue.CanWithdrawAmountOrWantChangeValue.Equal(getInfo.CanWithdrawAmountOrWantChangeValue))

	// test ErrInputUpdateStateIsZero
	/*	ethUniInitialChangeValue.TotalDepositAmountOrWantChangeValue = math.NewInt(0)
		ethUniInitialChangeValue.CanWithdrawAmountOrWantChangeValue = math.NewInt(0)
		err = suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUniAssetId, ethUniInitialChangeValue)
		suite.Require().Error(err, types2.ErrInputUpdateStateIsZero)*/

	// test valid increase of staker asset state
	ethUniInitialChangeValue.TotalDepositAmountOrWantChangeValue = math.NewInt(500)
	ethUniInitialChangeValue.CanWithdrawAmountOrWantChangeValue = math.NewInt(500)
	err = suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUniAssetId, ethUniInitialChangeValue)
	suite.Require().NoError(err)

	getInfo, err = suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, ethUniAssetId)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmountOrWantChangeValue.Equal(math.NewInt(1500)))
	suite.Require().True(getInfo.CanWithdrawAmountOrWantChangeValue.Equal(math.NewInt(1500)))

	// test valid decrease of staker asset state
	ethUniInitialChangeValue.TotalDepositAmountOrWantChangeValue = math.NewInt(-500)
	ethUniInitialChangeValue.CanWithdrawAmountOrWantChangeValue = math.NewInt(-500)
	err = suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUniAssetId, ethUniInitialChangeValue)
	suite.Require().NoError(err)
	getInfo, err = suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, ethUniAssetId)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmountOrWantChangeValue.Equal(math.NewInt(1000)))
	suite.Require().True(getInfo.CanWithdrawAmountOrWantChangeValue.Equal(math.NewInt(1000)))

	// test the decreased amount is bigger than original state
	ethUniInitialChangeValue.TotalDepositAmountOrWantChangeValue = math.NewInt(-2000)
	ethUniInitialChangeValue.CanWithdrawAmountOrWantChangeValue = math.NewInt(-500)
	err = suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUniAssetId, ethUniInitialChangeValue)
	suite.Require().Error(err, restakingtype.ErrSubAmountIsMoreThanOrigin)
	getInfo, err = suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, ethUniAssetId)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmountOrWantChangeValue.Equal(math.NewInt(1000)))
	suite.Require().True(getInfo.CanWithdrawAmountOrWantChangeValue.Equal(math.NewInt(1000)))

	ethUniInitialChangeValue.TotalDepositAmountOrWantChangeValue = math.NewInt(-500)
	ethUniInitialChangeValue.CanWithdrawAmountOrWantChangeValue = math.NewInt(-2000)
	err = suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUniAssetId, ethUniInitialChangeValue)
	suite.Require().Error(err, restakingtype.ErrSubAmountIsMoreThanOrigin)
	getInfo, err = suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, ethUniAssetId)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmountOrWantChangeValue.Equal(math.NewInt(1000)))
	suite.Require().True(getInfo.CanWithdrawAmountOrWantChangeValue.Equal(math.NewInt(1000)))

	// test the storage of multiple assets state
	ethUsdtAssetId := fmt.Sprintf("%s_%s", "0xdac17f958d2ee523a2206206994597c13d831ec7", "101")
	ethUsdtInitialChangeValue := restakingtype.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: math.NewInt(2000),
		CanWithdrawAmountOrWantChangeValue:  math.NewInt(2000),
	}
	err = suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUsdtAssetId, ethUsdtInitialChangeValue)
	suite.Require().NoError(err)
	getInfo, err = suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, ethUsdtAssetId)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmountOrWantChangeValue.Equal(math.NewInt(2000)))
	suite.Require().True(getInfo.CanWithdrawAmountOrWantChangeValue.Equal(math.NewInt(2000)))
}

func (suite *KeeperTestSuite) TestGetStakerAssetInfos() {
	stakerID := fmt.Sprintf("%s_%s", suite.address, "0")
	ethUniAssetId := fmt.Sprintf("%s_%s", "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984", "101")
	ethUsdtAssetId := fmt.Sprintf("%s_%s", "0xdac17f958d2ee523a2206206994597c13d831ec7", "101")
	ethUniInitialChangeValue := restakingtype.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: math.NewInt(1000),
		CanWithdrawAmountOrWantChangeValue:  math.NewInt(1000),
	}
	ethUsdtInitialChangeValue := restakingtype.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue: math.NewInt(2000),
		CanWithdrawAmountOrWantChangeValue:  math.NewInt(2000),
	}
	err := suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUniAssetId, ethUniInitialChangeValue)
	suite.Require().NoError(err)
	err = suite.app.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.ctx, stakerID, ethUsdtAssetId, ethUsdtInitialChangeValue)
	suite.Require().NoError(err)

	// test get all assets state of staker
	assetsInfo, err := suite.app.StakingAssetsManageKeeper.GetStakerAssetInfos(suite.ctx, stakerID)
	suite.Require().NoError(err)
	uniState, isExist := assetsInfo[ethUniAssetId]
	suite.Require().True(isExist)
	suite.Require().True(uniState.TotalDepositAmountOrWantChangeValue.Equal(math.NewInt(1000)))
	suite.Require().True(uniState.CanWithdrawAmountOrWantChangeValue.Equal(math.NewInt(1000)))

	usdtState, isExist := assetsInfo[ethUsdtAssetId]
	suite.Require().True(isExist)
	suite.Require().True(usdtState.TotalDepositAmountOrWantChangeValue.Equal(math.NewInt(2000)))
	suite.Require().True(usdtState.CanWithdrawAmountOrWantChangeValue.Equal(math.NewInt(2000)))
}
