package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	depositKeeper "github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	"github.com/ExocoreNetwork/exocore/x/withdraw/keeper"
	withdrawtype "github.com/ExocoreNetwork/exocore/x/withdraw/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *WithdrawTestSuite) TestClaimWithdrawRequest() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	event := &keeper.WithdrawParams{
		ClientChainLzID: 101,
		Action:          types.WithdrawPrinciple,
		WithdrawAddress: suite.Address[:],
		OpAmount:        sdkmath.NewInt(90),
	}

	depositEvent := &depositKeeper.DepositParams{
		ClientChainLzID: 101,
		Action:          types.Deposit,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(100),
	}

	// deposit firstly
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.App.DepositKeeper.Deposit(suite.Ctx, depositEvent)
	suite.NoError(err)

	// test the case that the withdraw asset hasn't registered
	event.AssetsAddress = usdcAddress[:]
	err = suite.App.WithdrawKeeper.Withdraw(suite.Ctx, event)
	suite.ErrorContains(err, withdrawtype.ErrWithdrawAssetNotExist.Error())

	assets, err := suite.App.StakingAssetsManageKeeper.GetAllStakingAssetsInfo(suite.Ctx)
	suite.NoError(err)
	suite.App.Logger().Info("the assets is:", "assets", assets)

	stakerID, assetID := types.GetStakeIDAndAssetID(depositEvent.ClientChainLzID, depositEvent.StakerAddress, depositEvent.AssetsAddress)
	info, err := suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetInfo{
		TotalDepositAmount:  depositEvent.OpAmount,
		WithdrawableAmount:  depositEvent.OpAmount,
		WaitUnbondingAmount: sdkmath.NewInt(0),
	}, *info)
	// test the normal case
	event.AssetsAddress = usdtAddress[:]
	err = suite.App.WithdrawKeeper.Withdraw(suite.Ctx, event)
	suite.NoError(err)

	// check state after withdraw
	stakerID, assetID = types.GetStakeIDAndAssetID(event.ClientChainLzID, event.WithdrawAddress, event.AssetsAddress)
	info, err = suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetInfo{
		TotalDepositAmount:  sdkmath.NewInt(10),
		WithdrawableAmount:  sdkmath.NewInt(10),
		WaitUnbondingAmount: sdkmath.NewInt(0),
	}, *info)

	assetInfo, err := suite.App.StakingAssetsManageKeeper.GetStakingAssetInfo(suite.Ctx, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(10), assetInfo.StakingTotalAmount)
}
