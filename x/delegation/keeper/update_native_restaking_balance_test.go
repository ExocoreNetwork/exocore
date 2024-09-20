package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *DelegationTestSuite) TestUpdateNativeRestakingBalance() {
	// test case: slash 80
	// withdrawable: 60 slash: 60 -> 0
	// undelegation: 10 slash: 10 -> 0
	// delegated amount:
	//	defaultOperator: 10 slash: 1/3*10 -> 6
	//  anotherOperator: 20 slash: 1/3*20 -> 13
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
	slashAmount := sdkmath.NewInt(80)
	stakerID, assetID := types.GetStakeIDAndAssetID(suite.clientChainLzID, suite.Address[:], suite.assetAddr.Bytes())
	err = suite.App.DelegationKeeper.UpdateNativeRestakingBalance(suite.Ctx, stakerID, assetID, slashAmount.Neg())
	suite.NoError(err)

}
