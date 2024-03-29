package keeper_test

import (
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
)

func (suite *KeeperTestSuite) TestGenesisClientChainAndAssetInfo() {
	defaultGensisState := restaking_assets_manage.DefaultGenesisState()

	// test the client chains getting
	clientChains, err := suite.app.StakingAssetsManageKeeper.GetAllClientChainInfo(suite.ctx)
	suite.NoError(err)
	suite.ctx.Logger().Info("the clientChains is:", "info", clientChains)
	for _, clientChain := range defaultGensisState.DefaultSupportedClientChains {
		info, ok := clientChains[clientChain.LayerZeroChainID]
		suite.True(ok)
		suite.Equal(info, clientChain)
	}

	chainInfo, err := suite.app.StakingAssetsManageKeeper.GetClientChainInfoByIndex(suite.ctx, 101)
	suite.NoError(err)
	suite.Equal(clientChains[101], chainInfo)

	// test the client chain assets getting
	assets, err := suite.app.StakingAssetsManageKeeper.GetAllStakingAssetsInfo(suite.ctx)
	suite.NoError(err)
	for _, asset := range defaultGensisState.DefaultSupportedClientChainTokens {
		_, assetID := types.GetStakeIDAndAssetIDFromStr(asset.LayerZeroChainID, "", asset.Address)
		suite.ctx.Logger().Info("the asset id is:", "assetID", assetID)
		info, ok := assets[assetID]
		suite.True(ok)
		suite.Equal(asset, info.AssetBasicInfo)
	}

	usdtAsset := defaultGensisState.DefaultSupportedClientChainTokens[0]
	_, assetID := types.GetStakeIDAndAssetIDFromStr(usdtAsset.LayerZeroChainID, "", usdtAsset.Address)
	assetInfo, err := suite.app.StakingAssetsManageKeeper.GetStakingAssetInfo(suite.ctx, assetID)
	suite.NoError(err)
	suite.Equal(usdtAsset, assetInfo.AssetBasicInfo)
}
