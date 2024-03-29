package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	depositKeeper "github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/ExocoreNetwork/exocore/x/slash/keeper"
	slashtype "github.com/ExocoreNetwork/exocore/x/slash/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *KeeperTestSuite) TestSlash() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	event := &keeper.SlashParams{
		ClientChainLzID: 101,
		Action:          types.Slash,
		StakerAddress:   suite.address[:],
		OpAmount:        sdkmath.NewInt(90),
	}

	depositEvent := &depositKeeper.DepositParams{
		ClientChainLzID: 101,
		Action:          types.Deposit,
		StakerAddress:   suite.address[:],
		OpAmount:        sdkmath.NewInt(100),
	}

	// deposit firstly
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.app.DepositKeeper.Deposit(suite.ctx, depositEvent)
	suite.NoError(err)

	// test the case that the slash  hasn't registered
	event.AssetsAddress = usdcAddress[:]
	err = suite.app.ExoSlashKeeper.Slash(suite.ctx, event)
	suite.ErrorContains(err, slashtype.ErrSlashAssetNotExist.Error())

	assets, err := suite.app.StakingAssetsManageKeeper.GetAllStakingAssetsInfo(suite.ctx)
	suite.NoError(err)
	suite.app.Logger().Info("the assets is:", "assets", assets)

	stakerID, assetID := types.GetStakeIDAndAssetID(depositEvent.ClientChainLzID, depositEvent.StakerAddress, depositEvent.AssetsAddress)
	info, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     depositEvent.OpAmount,
		CanWithdrawAmountOrWantChangeValue:      depositEvent.OpAmount,
		WaitUndelegationAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *info)

	// test the normal case
	event.AssetsAddress = usdtAddress[:]
	err = suite.app.ExoSlashKeeper.Slash(suite.ctx, event)
	suite.NoError(err)

	// check state after slash
	stakerID, assetID = types.GetStakeIDAndAssetID(event.ClientChainLzID, event.StakerAddress, event.AssetsAddress)
	info, err = suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     sdkmath.NewInt(10),
		CanWithdrawAmountOrWantChangeValue:      sdkmath.NewInt(10),
		WaitUndelegationAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *info)

	assetInfo, err := suite.app.StakingAssetsManageKeeper.GetStakingAssetInfo(suite.ctx, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(10), assetInfo.StakingTotalAmount)
}
