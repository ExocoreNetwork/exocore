package keeper_test

import (
	"fmt"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"

	"cosmossdk.io/math"
)

func (suite *StakingAssetsTestSuite) TestUpdateStakerAssetsState() {
	stakerID := fmt.Sprintf("%s_%s", suite.Address, "0")
	ethUniAssetID := fmt.Sprintf("%s_%s", "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984", "101")
	ethUniInitialChangeValue := assetstype.StakerSingleAssetChangeInfo{
		ChangeForTotalDeposit: math.NewInt(1000),
		ChangeForWithdrawable: math.NewInt(1000),
	}

	// test the initial storage of statker assets state
	err := suite.App.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.Ctx, stakerID, ethUniAssetID, ethUniInitialChangeValue)
	suite.Require().NoError(err)

	// test that the retrieved value is correct
	getInfo, err := suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, ethUniAssetID)
	suite.Require().NoError(err)
	suite.Require().True(ethUniInitialChangeValue.ChangeForTotalDeposit.Equal(getInfo.TotalDepositAmount))
	suite.Require().True(ethUniInitialChangeValue.ChangeForWithdrawable.Equal(getInfo.WithdrawableAmount))

	// test valid increase of staker asset state
	ethUniInitialChangeValue.ChangeForTotalDeposit = math.NewInt(500)
	ethUniInitialChangeValue.ChangeForWithdrawable = math.NewInt(500)
	err = suite.App.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.Ctx, stakerID, ethUniAssetID, ethUniInitialChangeValue)
	suite.Require().NoError(err)

	getInfo, err = suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, ethUniAssetID)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmount.Equal(math.NewInt(1500)))
	suite.Require().True(getInfo.WithdrawableAmount.Equal(math.NewInt(1500)))

	// test valid decrease of staker asset state
	ethUniInitialChangeValue.ChangeForTotalDeposit = math.NewInt(-500)
	ethUniInitialChangeValue.ChangeForWithdrawable = math.NewInt(-500)
	err = suite.App.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.Ctx, stakerID, ethUniAssetID, ethUniInitialChangeValue)
	suite.Require().NoError(err)
	getInfo, err = suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, ethUniAssetID)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmount.Equal(math.NewInt(1000)))
	suite.Require().True(getInfo.WithdrawableAmount.Equal(math.NewInt(1000)))

	// test the decreased amount is bigger than original state
	ethUniInitialChangeValue.ChangeForTotalDeposit = math.NewInt(-2000)
	ethUniInitialChangeValue.ChangeForWithdrawable = math.NewInt(-500)
	err = suite.App.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.Ctx, stakerID, ethUniAssetID, ethUniInitialChangeValue)
	suite.Require().Error(err, assetstype.ErrSubAmountIsMoreThanOrigin)
	getInfo, err = suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, ethUniAssetID)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmount.Equal(math.NewInt(1000)))
	suite.Require().True(getInfo.WithdrawableAmount.Equal(math.NewInt(1000)))

	ethUniInitialChangeValue.ChangeForTotalDeposit = math.NewInt(-500)
	ethUniInitialChangeValue.ChangeForWithdrawable = math.NewInt(-2000)
	err = suite.App.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.Ctx, stakerID, ethUniAssetID, ethUniInitialChangeValue)
	suite.Require().Error(err, assetstype.ErrSubAmountIsMoreThanOrigin)
	getInfo, err = suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, ethUniAssetID)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmount.Equal(math.NewInt(1000)))
	suite.Require().True(getInfo.WithdrawableAmount.Equal(math.NewInt(1000)))

	// test the storage of multiple assets state
	ethUsdtAssetID := fmt.Sprintf("%s_%s", "0xdac17f958d2ee523a2206206994597c13d831ec7", "101")
	ethUsdtInitialChangeValue := assetstype.StakerSingleAssetChangeInfo{
		ChangeForTotalDeposit: math.NewInt(2000),
		ChangeForWithdrawable: math.NewInt(2000),
	}
	err = suite.App.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.Ctx, stakerID, ethUsdtAssetID, ethUsdtInitialChangeValue)
	suite.Require().NoError(err)
	getInfo, err = suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, ethUsdtAssetID)
	suite.Require().NoError(err)
	suite.Require().True(getInfo.TotalDepositAmount.Equal(math.NewInt(2000)))
	suite.Require().True(getInfo.WithdrawableAmount.Equal(math.NewInt(2000)))
}

func (suite *StakingAssetsTestSuite) TestGetStakerAssetInfos() {
	stakerID := fmt.Sprintf("%s_%s", suite.Address, "0")
	ethUniAssetID := fmt.Sprintf("%s_%s", "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984", "101")
	ethUsdtAssetID := fmt.Sprintf("%s_%s", "0xdac17f958d2ee523a2206206994597c13d831ec7", "101")
	ethUniInitialChangeValue := assetstype.StakerSingleAssetChangeInfo{
		ChangeForTotalDeposit: math.NewInt(1000),
		ChangeForWithdrawable: math.NewInt(1000),
	}
	ethUsdtInitialChangeValue := assetstype.StakerSingleAssetChangeInfo{
		ChangeForTotalDeposit: math.NewInt(2000),
		ChangeForWithdrawable: math.NewInt(2000),
	}
	err := suite.App.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.Ctx, stakerID, ethUniAssetID, ethUniInitialChangeValue)
	suite.Require().NoError(err)
	err = suite.App.StakingAssetsManageKeeper.UpdateStakerAssetState(suite.Ctx, stakerID, ethUsdtAssetID, ethUsdtInitialChangeValue)
	suite.Require().NoError(err)

	// test get all assets state of staker
	assetsInfo, err := suite.App.StakingAssetsManageKeeper.GetStakerAssetInfos(suite.Ctx, stakerID)
	suite.Require().NoError(err)
	uniState, isExist := assetsInfo[ethUniAssetID]
	suite.Require().True(isExist)
	suite.Require().True(uniState.TotalDepositAmount.Equal(math.NewInt(1000)))
	suite.Require().True(uniState.WithdrawableAmount.Equal(math.NewInt(1000)))

	usdtState, isExist := assetsInfo[ethUsdtAssetID]
	suite.Require().True(isExist)
	suite.Require().True(usdtState.TotalDepositAmount.Equal(math.NewInt(2000)))
	suite.Require().True(usdtState.WithdrawableAmount.Equal(math.NewInt(2000)))
}
