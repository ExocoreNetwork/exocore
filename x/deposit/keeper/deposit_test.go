package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/exocore/x/deposit/keeper"
	types2 "github.com/exocore/x/deposit/types"
	"github.com/exocore/x/restaking_assets_manage/types"
)

func (suite *KeeperTestSuite) TestDeposit() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	event := &keeper.DepositParams{
		ClientChainLzId: 101,
		Action:          types.Deposit,
		StakerAddress:   suite.address[:],
		OpAmount:        sdkmath.NewInt(100),
	}

	//test the case that the deposit asset hasn't registered
	event.AssetsAddress = usdcAddress[:]
	err := suite.app.DepositKeeper.Deposit(suite.ctx, event)
	suite.ErrorContains(err, types2.ErrDepositAssetNotExist.Error())

	assets, err := suite.app.StakingAssetsManageKeeper.GetAllStakingAssetsInfo(suite.ctx)
	suite.NoError(err)
	suite.app.Logger().Info("the assets is:", "assets", assets)

	//test the normal case
	event.AssetsAddress = usdtAddress[:]
	err = suite.app.DepositKeeper.Deposit(suite.ctx, event)
	suite.NoError(err)

	//check state after deposit
	stakerId, assetId := types.GetStakeIDAndAssetId(event.ClientChainLzId, event.StakerAddress, event.AssetsAddress)
	info, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     event.OpAmount,
		CanWithdrawAmountOrWantChangeValue:      event.OpAmount,
		WaitUnDelegationAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *info)

	assetInfo, err := suite.app.StakingAssetsManageKeeper.GetStakingAssetInfo(suite.ctx, assetId)
	suite.NoError(err)
	suite.Equal(event.OpAmount, assetInfo.StakingTotalAmount)
}
