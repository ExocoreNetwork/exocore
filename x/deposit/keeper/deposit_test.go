package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	deposittype "github.com/ExocoreNetwork/exocore/x/deposit/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *DepositTestSuite) TestDeposit() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	params := &keeper.DepositParams{
		ClientChainLzID: 101,
		Action:          types.Deposit,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(100),
	}

	// test the case that the deposit asset hasn't registered
	params.AssetsAddress = usdcAddress[:]
	err := suite.App.DepositKeeper.Deposit(suite.Ctx, params)
	suite.ErrorContains(err, deposittype.ErrDepositAssetNotExist.Error())

	assets, err := suite.App.AssetsKeeper.GetAllStakingAssetsInfo(suite.Ctx)
	suite.NoError(err)
	suite.App.Logger().Info("the assets is:", "assets", assets)

	// test the normal case
	params.AssetsAddress = usdtAddress[:]
	err = suite.App.DepositKeeper.Deposit(suite.Ctx, params)
	suite.NoError(err)

	// check state after deposit
	stakerID, assetID := types.GetStakeIDAndAssetID(params.ClientChainLzID, params.StakerAddress, params.AssetsAddress)
	info, err := suite.App.AssetsKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerAssetInfo{
		TotalDepositAmount:  params.OpAmount,
		WithdrawableAmount:  params.OpAmount,
		WaitUnbondingAmount: sdkmath.NewInt(0),
	}, *info)

	assetInfo, err := suite.App.AssetsKeeper.GetStakingAssetInfo(suite.Ctx, assetID)
	suite.NoError(err)
	suite.Equal(params.OpAmount.Add(assets[assetID].StakingTotalAmount), assetInfo.StakingTotalAmount)
}
