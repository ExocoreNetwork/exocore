package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	operatorKeeper "github.com/ExocoreNetwork/exocore/x/operator/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *OperatorTestSuite) TestCalculateUSDValue() {
	suite.prepare()
	price, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
	suite.NoError(err)
	usdValue := operatorKeeper.CalculateUSDValue(suite.delegationAmount, price.Value, suite.assetDecimal, price.Decimal)
	expectedValue := sdkmath.LegacyNewDecFromBigInt(suite.delegationAmount.BigInt()).QuoInt(sdkmath.NewIntWithDecimal(1, int(suite.assetDecimal)))
	suite.Equal(expectedValue, usdValue)
	suite.Equal(int64(0), usdValue.TruncateInt64())
	float64Value, err := usdValue.Float64()
	suite.NoError(err)
	suite.Equal(5e-05, float64Value)
}

func (suite *OperatorTestSuite) TestCalculatedUSDValueOverflow() {
	price := assetstype.MaxAssetTotalSupply
	priceDecimal := uint8(assetstype.MaxDecimal)
	amount := assetstype.MaxAssetTotalSupply
	assetDecimal := uint32(assetstype.MaxDecimal)
	usdValue := operatorKeeper.CalculateUSDValue(amount, price, assetDecimal, priceDecimal)
	expectedValue := sdkmath.LegacyNewDecFromBigInt(sdkmath.NewIntWithDecimal(1, 2*assetstype.MaxDecForTotalSupply-2*assetstype.MaxDecimal).BigInt())
	suite.Equal(expectedValue, usdValue)

	priceDecimal = uint8(0)
	assetDecimal = uint32(0)
	usdValue = operatorKeeper.CalculateUSDValue(amount, price, assetDecimal, priceDecimal)
	expectedValue = sdkmath.LegacyNewDecFromBigInt(sdkmath.NewIntWithDecimal(1, 2*assetstype.MaxDecForTotalSupply).BigInt())
	suite.Equal(expectedValue, usdValue)

	price = sdkmath.NewInt(1)
	priceDecimal = uint8(assetstype.MaxDecimal)
	amount = sdkmath.NewInt(1)
	assetDecimal = uint32(assetstype.MaxDecimal)
	usdValue = operatorKeeper.CalculateUSDValue(amount, price, assetDecimal, priceDecimal)
	expectedValue = sdkmath.LegacyNewDec(0)
	suite.Equal(expectedValue.String(), usdValue.String())

	price = sdkmath.NewInt(1)
	priceDecimal = uint8(0)
	amount = sdkmath.NewInt(1)
	assetDecimal = uint32(assetstype.MaxDecimal)
	usdValue = operatorKeeper.CalculateUSDValue(amount, price, assetDecimal, priceDecimal)
	expectedValue = sdkmath.LegacyNewDecFromBigIntWithPrec(amount.BigInt(), sdkmath.LegacyPrecision)
	suite.Equal(expectedValue, usdValue)
	float64Value, err := usdValue.Float64()
	suite.NoError(err)
	suite.Equal(1e-18, float64Value)
}

func (suite *OperatorTestSuite) TestAVSUSDValue() {
	suite.prepare()
	err := suite.App.OperatorKeeper.OptIn(suite.Ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	usdtPrice, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
	suite.NoError(err)
	usdtValue := operatorKeeper.CalculateUSDValue(suite.delegationAmount, usdtPrice.Value, suite.assetDecimal, usdtPrice.Decimal)
	// deposit and delegate another asset to the operator
	usdcAddr := common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	usdcClientChainAsset := assetstype.AssetInfo{
		Name:             "USD coin",
		Symbol:           "USDC",
		Address:          usdcAddr.String(),
		Decimals:         6,
		TotalSupply:      sdkmath.NewInt(1e18),
		LayerZeroChainID: 101,
		MetaInfo:         "USDC",
	}
	_, err = suite.App.AssetsKeeper.RegisterAsset(
		suite.Ctx,
		&assetstype.RegisterAssetReq{
			FromAddress: suite.AccAddress.String(),
			Info:        &usdcClientChainAsset,
		})
	suite.NoError(err)
	suite.prepareDeposit(usdcAddr, sdkmath.NewInt(1e8))
	usdcPrice, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
	suite.NoError(err)
	delegatedAmount := sdkmath.NewIntWithDecimal(8, 7)
	suite.prepareDelegation(usdcAddr, delegatedAmount)

	// updating the new voting power
	suite.NoError(err)
	usdcValue := operatorKeeper.CalculateUSDValue(suite.delegationAmount, usdcPrice.Value, suite.assetDecimal, usdcPrice.Decimal)
	expectedUSDvalue := usdcValue.Add(usdtValue)
	suite.App.OperatorKeeper.EndBlock(suite.Ctx, abci.RequestEndBlock{})
	avsUSDValue, err := suite.App.OperatorKeeper.GetAVSUSDValue(suite.Ctx, suite.avsAddr)
	suite.NoError(err)
	suite.Equal(expectedUSDvalue, avsUSDValue)
	operatorUSDValue, err := suite.App.OperatorKeeper.GetOperatorUSDValue(suite.Ctx, suite.avsAddr, suite.operatorAddr.String())
	suite.NoError(err)
	suite.Equal(expectedUSDvalue, operatorUSDValue)
}
