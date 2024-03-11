package keeper_test

import (
	"github.com/ExocoreNetwork/exocore/x/assets"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
)

func (suite *StakingAssetsTestSuite) TestGenesisClientChainAndAssetInfo() {
	defaultGensisState := assets.DefaultGenesisState()

	// test the client chains getting
	clientChains, err := suite.App.StakingAssetsManageKeeper.GetAllClientChainInfo(suite.Ctx)
	suite.NoError(err)
	suite.Ctx.Logger().Info("the clientChains is:", "info", clientChains)
	for _, clientChain := range defaultGensisState.DefaultSupportedClientChains {
		info, ok := clientChains[clientChain.LayerZeroChainID]
		suite.True(ok)
		suite.Equal(info, clientChain)
	}

	chainInfo, err := suite.App.StakingAssetsManageKeeper.GetClientChainInfoByIndex(suite.Ctx, 101)
	suite.NoError(err)
	suite.Equal(clientChains[101], chainInfo)

	// test the client chain assets getting
	assets, err := suite.App.StakingAssetsManageKeeper.GetAllStakingAssetsInfo(suite.Ctx)
	suite.NoError(err)
	for _, asset := range defaultGensisState.DefaultSupportedClientChainTokens {
		_, assetID := assetstype.GetStakeIDAndAssetIDFromStr(asset.LayerZeroChainID, "", asset.Address)
		suite.Ctx.Logger().Info("the asset id is:", "assetID", assetID)
		info, ok := assets[assetID]
		suite.True(ok)
		suite.Equal(asset, info.AssetBasicInfo)
	}

	usdtAsset := defaultGensisState.DefaultSupportedClientChainTokens[0]
	_, assetID := assetstype.GetStakeIDAndAssetIDFromStr(usdtAsset.LayerZeroChainID, "", usdtAsset.Address)
	assetInfo, err := suite.App.StakingAssetsManageKeeper.GetStakingAssetInfo(suite.Ctx, assetID)
	suite.NoError(err)
	suite.Equal(usdtAsset, assetInfo.AssetBasicInfo)
}
