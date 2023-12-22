package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	depositKeeper "github.com/exocore/x/deposit/keeper"
	"github.com/exocore/x/restaking_assets_manage/types"
	"github.com/exocore/x/withdraw/keeper"
	withdrawtype "github.com/exocore/x/withdraw/types"
)

func (suite *KeeperTestSuite) TestClaimWithdrawRequest() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	event := &keeper.WithdrawParams{
		ClientChainLzId: 101,
		Action:          types.WithdrawPrinciple,
		WithdrawAddress: suite.address[:],
		OpAmount:        sdkmath.NewInt(90),
	}

	depositEvent := &depositKeeper.DepositParams{
		ClientChainLzId: 101,
		Action:          types.Deposit,
		StakerAddress:   suite.address[:],
		OpAmount:        sdkmath.NewInt(100),
	}

	// deposit firstly
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.app.DepositKeeper.Deposit(suite.ctx, depositEvent)
	suite.NoError(err)

	// test the case that the withdraw asset hasn't registered
	event.AssetsAddress = usdcAddress[:]
	err = suite.app.WithdrawKeeper.Withdraw(suite.ctx, event)
	suite.ErrorContains(err, withdrawtype.ErrWithdrawAssetNotExist.Error())

	assets, err := suite.app.StakingAssetsManageKeeper.GetAllStakingAssetsInfo(suite.ctx)
	suite.NoError(err)
	suite.app.Logger().Info("the assets is:", "assets", assets)

	stakerId, assetId := types.GetStakeIDAndAssetId(depositEvent.ClientChainLzId, depositEvent.StakerAddress, depositEvent.AssetsAddress)
	info, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     depositEvent.OpAmount,
		CanWithdrawAmountOrWantChangeValue:      depositEvent.OpAmount,
		WaitUndelegationAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *info)
	// test the normal case
	event.AssetsAddress = usdtAddress[:]
	err = suite.app.WithdrawKeeper.Withdraw(suite.ctx, event)
	suite.NoError(err)

	// check state after withdraw
	stakerId, assetId = types.GetStakeIDAndAssetId(event.ClientChainLzId, event.WithdrawAddress, event.AssetsAddress)
	info, err = suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     sdkmath.NewInt(10),
		CanWithdrawAmountOrWantChangeValue:      sdkmath.NewInt(10),
		WaitUndelegationAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *info)

	assetInfo, err := suite.app.StakingAssetsManageKeeper.GetStakingAssetInfo(suite.ctx, assetId)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(10), assetInfo.StakingTotalAmount)
}
