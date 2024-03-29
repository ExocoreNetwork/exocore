package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	"github.com/ExocoreNetwork/exocore/x/reward/keeper"
	rewardtype "github.com/ExocoreNetwork/exocore/x/reward/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *KeeperTestSuite) TestClaimWithdrawRequest() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	event := &keeper.RewardParams{
		ClientChainLzID:       101,
		Action:                types.WithDrawReward,
		WithdrawRewardAddress: suite.address[:],
		OpAmount:              sdkmath.NewInt(10),
	}

	// test the case that the deposit asset hasn't registered
	event.AssetsAddress = usdcAddress[:]
	err := suite.app.RewardKeeper.RewardForWithdraw(suite.ctx, event)
	suite.ErrorContains(err, rewardtype.ErrRewardAssetNotExist.Error())

	// test the normal case
	event.AssetsAddress = usdtAddress[:]
	err = suite.app.RewardKeeper.RewardForWithdraw(suite.ctx, event)
	suite.NoError(err)

	// check state after reward
	stakerID, assetID := types.GetStakeIDAndAssetID(event.ClientChainLzID, event.WithdrawRewardAddress, event.AssetsAddress)
	info, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerID, assetID)
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
