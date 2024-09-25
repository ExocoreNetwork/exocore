package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	assettypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *DelegationTestSuite) TestUpdateNativeRestakingBalance() {
	// test case: slash 80
	// withdrawable: 60 slash: 60 -> 0
	// undelegation: 10 slash: 10 -> 0
	// delegated amount:
	//	defaultOperator: 10 slash: 1/3*10 -> 2/3*10
	//  anotherOperator: 20 slash: 1/3*20 -> 2/3*20
	depositAmount := sdkmath.NewInt(100)
	delegateAmountToDefaultOperator := sdkmath.NewInt(20)
	undelegateAmountFromDefaultOperator := sdkmath.NewInt(10)
	delegateAmountToAnotherOperator := sdkmath.NewInt(20)
	anotherOperatorAddr, err := sdk.AccAddressFromBech32("exo18cggcpvwspnd5c6ny8wrqxpffj5zmhklprtnph")
	suite.NoError(err)

	suite.basicPrepare()
	suite.prepareDeposit(depositAmount)

	delegationEvent := suite.prepareDelegation(delegateAmountToDefaultOperator, suite.opAccAddr)
	delegationEvent.LzNonce = 1
	delegationEvent.OpAmount = undelegateAmountFromDefaultOperator
	err = suite.App.DelegationKeeper.UndelegateFrom(suite.Ctx, delegationEvent)
	suite.NoError(err)
	suite.prepareDelegation(delegateAmountToAnotherOperator, anotherOperatorAddr)

	// update negative balance
	// The actual slash amount is 79, not 80; this is due to precision loss.
	slashAmount := sdkmath.NewInt(80)
	actualSlashAmount := sdkmath.NewInt(79)
	stakerID, assetID := assettypes.GetStakerIDAndAssetID(suite.clientChainLzID, suite.Address[:], suite.assetAddr.Bytes())
	err = suite.App.DelegationKeeper.UpdateNativeRestakingBalance(suite.Ctx, stakerID, assetID, slashAmount.Neg())
	suite.NoError(err)

	// check the asset state for the slashed staker
	stakerAssetInfo, err := suite.App.AssetsKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	expectAssetInfo := assettypes.StakerAssetInfo{
		TotalDepositAmount: depositAmount.Sub(actualSlashAmount),
		WithdrawableAmount: sdkmath.NewInt(0),
		// it will be decreased when the undelegation is completed.
		PendingUndelegationAmount: undelegateAmountFromDefaultOperator,
	}
	suite.Equal(expectAssetInfo, *stakerAssetInfo)

	// check the undelegation state after slashing
	records, err := suite.App.DelegationKeeper.GetStakerUndelegationRecords(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(0), records[0].ActualCompletedAmount)

	// check the delegated share for two operators
	delegationForDefaultOperator, err := suite.App.DelegationKeeper.GetSingleDelegationInfo(suite.Ctx, stakerID, assetID, suite.opAccAddr.String())
	suite.NoError(err)
	slashProportion := sdkmath.LegacyNewDec(1).Sub((sdkmath.LegacyNewDec(1)).Quo(sdkmath.LegacyNewDec(3)))
	expectedShareForDefaultOperator := sdkmath.LegacyNewDec(10).Mul(slashProportion)
	suite.Equal(expectedShareForDefaultOperator, delegationForDefaultOperator.UndelegatableShare)

	delegationForAnotherOperator, err := suite.App.DelegationKeeper.GetSingleDelegationInfo(suite.Ctx, stakerID, assetID, anotherOperatorAddr.String())
	suite.NoError(err)
	expectedShareForAnotherOperator := sdkmath.LegacyNewDec(20).Mul(slashProportion)
	suite.Equal(expectedShareForAnotherOperator, delegationForAnotherOperator.UndelegatableShare)

	// check the asset states of two operators
	defaultOperatorAsset, err := suite.App.AssetsKeeper.GetOperatorSpecifiedAssetInfo(suite.Ctx, suite.opAccAddr, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(7), defaultOperatorAsset.TotalAmount)
	suite.Equal(expectedShareForDefaultOperator, defaultOperatorAsset.TotalShare)

	anotherOperatorAsset, err := suite.App.AssetsKeeper.GetOperatorSpecifiedAssetInfo(suite.Ctx, anotherOperatorAddr, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(14), anotherOperatorAsset.TotalAmount)
	suite.Equal(expectedShareForAnotherOperator, anotherOperatorAsset.TotalShare)
}
