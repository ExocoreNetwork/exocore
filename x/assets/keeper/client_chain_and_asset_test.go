package keeper_test

import (
	"cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *StakingAssetsTestSuite) TestGenesisClientChainAndAssetInfo() {
	ethClientChain := assetstype.ClientChainInfo{
		Name:               "ethereum",
		MetaInfo:           "ethereum blockchain",
		ChainId:            1,
		FinalizationBlocks: 10,
		LayerZeroChainID:   101,
		AddressLength:      20,
	}
	usdtClientChainAsset := assetstype.AssetInfo{
		Name:             "Tether USD",
		Symbol:           "USDT",
		Address:          "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		Decimals:         6,
		LayerZeroChainID: ethClientChain.LayerZeroChainID,
		MetaInfo:         "Tether USD token",
	}
	totalSupply, _ := sdk.NewIntFromString("40022689732746729")
	usdtClientChainAsset.TotalSupply = totalSupply
	stakingInfo := assetstype.StakingAssetInfo{
		AssetBasicInfo:     &usdtClientChainAsset,
		StakingTotalAmount: math.NewInt(0),
	}
	defaultGensisState := assetstype.NewGenesis(
		assetstype.DefaultParams(),
		[]assetstype.ClientChainInfo{ethClientChain},
		[]assetstype.StakingAssetInfo{stakingInfo},
		[]assetstype.DepositsByStaker{},
	)

	// test the client chains getting
	clientChains, err := suite.App.AssetsKeeper.GetAllClientChainInfo(suite.Ctx)
	suite.NoError(err)
	suite.Ctx.Logger().Info("the clientChains is:", "info", clientChains)
	for _, clientChain := range defaultGensisState.ClientChains {
		suite.Contains(clientChains, clientChain)
	}

	chainInfo, err := suite.App.AssetsKeeper.GetClientChainInfoByIndex(suite.Ctx, 101)
	suite.NoError(err)
	suite.Contains(clientChains, *chainInfo)

	// test the client chain assets getting
	assets, err := suite.App.AssetsKeeper.GetAllStakingAssetsInfo(suite.Ctx)
	suite.NoError(err)
	for _, assetX := range defaultGensisState.Tokens {
		asset := assetX.AssetBasicInfo
		_, assetID := assetstype.GetStakeIDAndAssetIDFromStr(asset.LayerZeroChainID, "", asset.Address)
		suite.Ctx.Logger().Info("the asset id is:", "assetID", assetID)
		info, ok := assets[assetID]
		suite.True(ok)
		suite.Equal(asset, info.AssetBasicInfo)
	}

	usdtAssetX := defaultGensisState.Tokens[0]
	usdtAsset := usdtAssetX.AssetBasicInfo
	_, assetID := assetstype.GetStakeIDAndAssetIDFromStr(usdtAsset.LayerZeroChainID, "", usdtAsset.Address)
	assetInfo, err := suite.App.AssetsKeeper.GetStakingAssetInfo(suite.Ctx, assetID)
	suite.NoError(err)
	suite.Equal(usdtAsset, assetInfo.AssetBasicInfo)
}
