package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
)

func (suite *DelegationTestSuite) TestTokensFromShares() {
	testCases := []struct {
		// input
		totalShare     sdkmath.LegacyDec
		operatorAmount sdkmath.Int
		stakerShare    sdkmath.LegacyDec
		// output
		stakerAmount sdkmath.Int
		innerError   error
	}{
		// error cases
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(50),
			stakerShare:    sdkmath.LegacyNewDec(51),
			innerError:     delegationtypes.ErrInsufficientShares,
			stakerAmount:   sdkmath.NewInt(0),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(0),
			operatorAmount: sdkmath.NewInt(50),
			stakerShare:    sdkmath.LegacyNewDec(0),
			innerError:     delegationtypes.ErrDivisorIsZero,
			stakerAmount:   sdkmath.NewInt(0),
		},

		// the share will be equal to the amount if there isn't a slash event
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(50),
			stakerShare:    sdkmath.LegacyNewDec(0),
			innerError:     nil,
			stakerAmount:   sdkmath.NewInt(0),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(50),
			stakerShare:    sdkmath.LegacyMustNewDecFromStr("3.4"),
			innerError:     nil,
			stakerAmount:   sdkmath.NewInt(3),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(50),
			stakerShare:    sdkmath.LegacyNewDec(50),
			innerError:     nil,
			stakerAmount:   sdkmath.NewInt(50),
		},

		// the share will be greater than the amount if there is a slash event
		{
			totalShare:     sdkmath.LegacyNewDec(70),
			operatorAmount: sdkmath.NewInt(50),
			stakerShare:    sdkmath.LegacyNewDec(0),
			innerError:     nil,
			stakerAmount:   sdkmath.NewInt(0),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(70),
			operatorAmount: sdkmath.NewInt(50),
			stakerShare:    sdkmath.LegacyMustNewDecFromStr("3.4"),
			innerError:     nil,
			stakerAmount:   sdkmath.NewInt(2),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(70),
			operatorAmount: sdkmath.NewInt(50),
			stakerShare:    sdkmath.LegacyNewDec(70),
			innerError:     nil,
			stakerAmount:   sdkmath.NewInt(50),
		},
	}

	for _, testCase := range testCases {
		amount, err := keeper.TokensFromShares(testCase.stakerShare, testCase.totalShare, testCase.operatorAmount)
		if testCase.innerError != nil {
			suite.ErrorContains(err, testCase.innerError.Error())
		} else {
			suite.NoError(err)
			suite.Equal(testCase.stakerAmount, amount)
		}
	}
}

func (suite *DelegationTestSuite) TestSharesFromTokens() {
	testCases := []struct {
		// input
		totalShare     sdkmath.LegacyDec
		operatorAmount sdkmath.Int
		stakerAmount   sdkmath.Int

		// output
		stakerShare sdkmath.LegacyDec
		innerError  error
	}{
		// error cases
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(50),
			stakerAmount:   sdkmath.NewInt(51),
			innerError:     delegationtypes.ErrInsufficientAssetAmount,
			stakerShare:    sdkmath.LegacyNewDec(0),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(0),
			stakerAmount:   sdkmath.NewInt(0),
			innerError:     delegationtypes.ErrDivisorIsZero,
			stakerShare:    sdkmath.LegacyNewDec(0),
		},

		// the share will be equal to the amount if there isn't a slash event
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(50),
			stakerAmount:   sdkmath.NewInt(0),
			innerError:     nil,
			stakerShare:    sdkmath.LegacyNewDec(0),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(50),
			stakerAmount:   sdkmath.NewInt(3),
			innerError:     nil,
			stakerShare:    sdkmath.LegacyNewDec(3),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(50),
			operatorAmount: sdkmath.NewInt(50),
			stakerAmount:   sdkmath.NewInt(50),
			innerError:     nil,
			stakerShare:    sdkmath.LegacyNewDec(50),
		},

		// the share will be greater than the amount if there is a slash event
		{
			totalShare:     sdkmath.LegacyNewDec(70),
			operatorAmount: sdkmath.NewInt(50),
			stakerAmount:   sdkmath.NewInt(0),
			innerError:     nil,
			stakerShare:    sdkmath.LegacyNewDec(0),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(70),
			operatorAmount: sdkmath.NewInt(50),
			stakerAmount:   sdkmath.NewInt(2),
			innerError:     nil,
			stakerShare:    sdkmath.LegacyMustNewDecFromStr("2.8"),
		},
		{
			totalShare:     sdkmath.LegacyNewDec(70),
			operatorAmount: sdkmath.NewInt(50),
			stakerAmount:   sdkmath.NewInt(50),
			innerError:     nil,
			stakerShare:    sdkmath.LegacyNewDec(70),
		},
	}

	for _, testCase := range testCases {
		share, err := keeper.SharesFromTokens(testCase.totalShare, testCase.stakerAmount, testCase.operatorAmount)
		if testCase.innerError != nil {
			suite.ErrorContains(err, testCase.innerError.Error())
		} else {
			suite.NoError(err)
			suite.Equal(testCase.stakerShare, share)
		}
	}
}

func (suite *DelegationTestSuite) TestCalculateShare() {
	assetAmount := sdkmath.NewInt(10)
	// test the case that the operator doesn't exist
	_, assetID := assetstype.GetStakeIDAndAssetID(suite.clientChainLzID, nil, suite.assetAddr[:])
	share, err := suite.App.DelegationKeeper.CalculateShare(suite.Ctx, suite.opAccAddr, assetID, assetAmount)
	suite.NoError(err)
	suite.Equal(sdkmath.LegacyNewDecFromBigInt(assetAmount.BigInt()), share)

	// test the case that the asset amount of operator is zero
	err = suite.App.AssetsKeeper.UpdateOperatorAssetState(suite.Ctx, suite.opAccAddr, assetID, assetstype.DeltaOperatorSingleAsset{
		TotalAmount: sdkmath.NewInt(0),
		TotalShare:  sdkmath.LegacyNewDec(0),
	})
	share, err = suite.App.DelegationKeeper.CalculateShare(suite.Ctx, suite.opAccAddr, assetID, assetAmount)
	suite.NoError(err)
	suite.Equal(sdkmath.LegacyNewDecFromBigInt(assetAmount.BigInt()), share)

	// test normal cases
	err = suite.App.AssetsKeeper.UpdateOperatorAssetState(suite.Ctx, suite.opAccAddr, assetID, assetstype.DeltaOperatorSingleAsset{
		TotalAmount: sdkmath.NewInt(50),
		TotalShare:  sdkmath.LegacyNewDec(60),
	})
	share, err = suite.App.DelegationKeeper.CalculateShare(suite.Ctx, suite.opAccAddr, assetID, assetAmount)
	suite.NoError(err)
	suite.Equal(sdkmath.LegacyNewDec(12), share)
}

func (suite *DelegationTestSuite) TestValidateUndeleagtionAmount() {
	suite.prepareDeposit()
	suite.prepareDelegation()
	stakerID, assetID := assetstype.GetStakeIDAndAssetID(suite.clientChainLzID, suite.Address[:], suite.assetAddr[:])

	undelegationAmount := sdkmath.NewInt(0)
	share, err := suite.App.DelegationKeeper.ValidateUndeleagtionAmount(suite.Ctx, suite.opAccAddr, stakerID, assetID, undelegationAmount)
	suite.NoError(err)
	suite.Equal(sdkmath.LegacyNewDec(0), share)

	undelegationAmount = sdkmath.NewInt(10)
	share, err = suite.App.DelegationKeeper.ValidateUndeleagtionAmount(suite.Ctx, suite.opAccAddr, stakerID, assetID, undelegationAmount)
	suite.NoError(err)
	suite.Equal(sdkmath.LegacyNewDecFromBigInt(undelegationAmount.BigInt()), share)

	undelegationAmount = suite.delegationAmount.Add(sdkmath.NewInt(1))
	share, err = suite.App.DelegationKeeper.ValidateUndeleagtionAmount(suite.Ctx, suite.opAccAddr, stakerID, assetID, undelegationAmount)
	suite.Error(err, delegationtypes.ErrInsufficientShares)
}

func (suite *DelegationTestSuite) TestCalculateSlashShare() {
	suite.prepareDeposit()
	suite.prepareDelegation()
	stakerID, assetID := assetstype.GetStakeIDAndAssetID(suite.clientChainLzID, suite.Address[:], suite.assetAddr[:])
	slashAmount := sdkmath.NewInt(0)
	suite.App.DelegationKeeper.CalculateSlashShare(suite.Ctx, suite.opAccAddr, stakerID, assetID, slashAmount)
}
