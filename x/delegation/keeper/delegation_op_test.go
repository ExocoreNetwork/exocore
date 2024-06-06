package keeper_test

import (
	"fmt"

	assetskeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	operatortype "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *DelegationTestSuite) prepare() {
	suite.assetAddr = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	suite.clientChainLzID = uint64(101)
	opAccAddr, err := sdk.AccAddressFromBech32("exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr")
	suite.NoError(err)
	suite.opAccAddr = opAccAddr
	suite.depositAmount = sdkmath.NewInt(100)
	suite.delegationAmount = sdkmath.NewInt(50)
}

func (suite *DelegationTestSuite) prepareDeposit() *assetskeeper.DepositWithdrawParams {
	suite.prepare()
	depositEvent := &assetskeeper.DepositWithdrawParams{
		ClientChainLzID: suite.clientChainLzID,
		Action:          types.Deposit,
		StakerAddress:   suite.Address[:],
		OpAmount:        suite.depositAmount,
	}
	depositEvent.AssetsAddress = suite.assetAddr[:]
	err := suite.App.AssetsKeeper.PerformDepositOrWithdraw(suite.Ctx, depositEvent)
	suite.NoError(err)
	return depositEvent
}

func (suite *DelegationTestSuite) prepareDelegation() *delegationtype.DelegationOrUndelegationParams {
	delegationEvent := &delegationtype.DelegationOrUndelegationParams{
		ClientChainLzID: suite.clientChainLzID,
		Action:          types.DelegateTo,
		AssetsAddress:   suite.assetAddr.Bytes(),
		OperatorAddress: suite.opAccAddr,
		StakerAddress:   suite.Address[:],
		OpAmount:        suite.delegationAmount,
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	registerReq := &operatortype.RegisterOperatorReq{
		FromAddress: suite.opAccAddr.String(),
		Info: &operatortype.OperatorInfo{
			EarningsAddr: suite.opAccAddr.String(),
		},
	}
	_, err := suite.App.OperatorKeeper.RegisterOperator(suite.Ctx, registerReq)
	suite.NoError(err)

	err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, delegationEvent)
	suite.NoError(err)
	return delegationEvent
}

func (suite *DelegationTestSuite) TestDelegateTo() {
	suite.prepareDeposit()
	opAccAddr, err := sdk.AccAddressFromBech32("exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr")
	suite.NoError(err)
	delegationParams := &delegationtype.DelegationOrUndelegationParams{
		ClientChainLzID: suite.clientChainLzID,
		Action:          types.DelegateTo,
		AssetsAddress:   suite.assetAddr.Bytes(),
		OperatorAddress: opAccAddr,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(50),
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, delegationParams)
	suite.EqualError(err, errorsmod.Wrap(delegationtype.ErrOperatorNotExist, fmt.Sprintf("input operatorAddr is:%s", delegationParams.OperatorAddress)).Error())

	registerReq := &operatortype.RegisterOperatorReq{
		FromAddress: opAccAddr.String(),
		Info: &operatortype.OperatorInfo{
			EarningsAddr: opAccAddr.String(),
		},
	}
	_, err = suite.App.OperatorKeeper.RegisterOperator(suite.Ctx, registerReq)
	suite.NoError(err)

	err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, delegationParams)
	suite.NoError(err)

	// check delegation states
	stakerID, assetID := types.GetStakeIDAndAssetID(delegationParams.ClientChainLzID, delegationParams.StakerAddress, delegationParams.AssetsAddress)
	restakerState, err := suite.App.AssetsKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerAssetInfo{
		TotalDepositAmount:  suite.depositAmount,
		WithdrawableAmount:  suite.depositAmount.Sub(delegationParams.OpAmount),
		WaitUnbondingAmount: sdkmath.NewInt(0),
	}, *restakerState)

	operatorState, err := suite.App.AssetsKeeper.GetOperatorSpecifiedAssetInfo(suite.Ctx, opAccAddr, assetID)
	suite.NoError(err)
	suite.Equal(types.OperatorAssetInfo{
		TotalAmount:             delegationParams.OpAmount,
		OperatorAmount:          sdkmath.NewInt(0),
		WaitUnbondingAmount:     sdkmath.NewInt(0),
		OperatorUnbondingAmount: sdkmath.NewInt(0),
		TotalShare:              sdkmath.LegacyNewDecFromBigInt(delegationParams.OpAmount.BigInt()),
		OperatorShare:           sdkmath.LegacyNewDec(0),
	}, *operatorState)

	specifiedDelegationAmount, err := suite.App.DelegationKeeper.GetSingleDelegationInfo(suite.Ctx, stakerID, assetID, opAccAddr.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		UndelegatableShare:     sdkmath.LegacyNewDecFromBigInt(delegationParams.OpAmount.BigInt()),
		WaitUndelegationAmount: sdkmath.NewInt(0),
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.App.DelegationKeeper.StakerDelegatedTotalAmount(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(delegationParams.OpAmount, totalDelegationAmount)
}

func (suite *DelegationTestSuite) TestUndelegateFrom() {
	suite.prepareDeposit()
	delegationEvent := suite.prepareDelegation()
	// test Undelegation
	delegationEvent.LzNonce = 1
	err := suite.App.DelegationKeeper.UndelegateFrom(suite.Ctx, delegationEvent)
	suite.NoError(err)

	// check state
	stakerID, assetID := types.GetStakeIDAndAssetID(delegationEvent.ClientChainLzID, delegationEvent.StakerAddress, delegationEvent.AssetsAddress)
	restakerState, err := suite.App.AssetsKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerAssetInfo{
		TotalDepositAmount:  suite.depositAmount,
		WithdrawableAmount:  suite.depositAmount.Sub(delegationEvent.OpAmount),
		WaitUnbondingAmount: delegationEvent.OpAmount,
	}, *restakerState)

	operatorState, err := suite.App.AssetsKeeper.GetOperatorSpecifiedAssetInfo(suite.Ctx, delegationEvent.OperatorAddress, assetID)
	suite.NoError(err)
	suite.Equal(types.OperatorAssetInfo{
		TotalAmount:             sdkmath.NewInt(0),
		OperatorAmount:          sdkmath.NewInt(0),
		WaitUnbondingAmount:     delegationEvent.OpAmount,
		OperatorUnbondingAmount: sdkmath.NewInt(0),
		TotalShare:              sdkmath.LegacyNewDec(0),
		OperatorShare:           sdkmath.LegacyNewDec(0),
	}, *operatorState)

	specifiedDelegationAmount, err := suite.App.DelegationKeeper.GetSingleDelegationInfo(suite.Ctx, stakerID, assetID, delegationEvent.OperatorAddress.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		WaitUndelegationAmount: delegationEvent.OpAmount,
		UndelegatableShare:     sdkmath.LegacyNewDec(0),
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.App.DelegationKeeper.StakerDelegatedTotalAmount(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(0), totalDelegationAmount)

	records, err := suite.App.DelegationKeeper.GetStakerUndelegationRecords(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(1, len(records))
	UndelegationRecord := &delegationtype.UndelegationRecord{
		StakerID:              stakerID,
		AssetID:               assetID,
		OperatorAddr:          delegationEvent.OperatorAddress.String(),
		TxHash:                delegationEvent.TxHash.String(),
		IsPending:             true,
		BlockNumber:           uint64(suite.Ctx.BlockHeight()),
		LzTxNonce:             delegationEvent.LzNonce,
		Amount:                delegationEvent.OpAmount,
		ActualCompletedAmount: delegationEvent.OpAmount,
	}
	UndelegationRecord.CompleteBlockNumber = UndelegationRecord.BlockNumber + delegationtype.CanUndelegationDelayHeight
	suite.Equal(UndelegationRecord, records[0])

	suite.Ctx.Logger().Info("the complete block number is:", "height", UndelegationRecord.CompleteBlockNumber)
	waitUndelegationRecords, err := suite.App.DelegationKeeper.GetWaitCompleteUndelegationRecords(suite.Ctx, UndelegationRecord.CompleteBlockNumber)
	suite.NoError(err)
	suite.Equal(1, len(waitUndelegationRecords))
	suite.Equal(UndelegationRecord, waitUndelegationRecords[0])
}

func (suite *DelegationTestSuite) TestCompleteUndelegation() {
	suite.prepareDeposit()
	delegationEvent := suite.prepareDelegation()

	delegationEvent.LzNonce = 1
	err := suite.App.DelegationKeeper.UndelegateFrom(suite.Ctx, delegationEvent)
	suite.NoError(err)
	UndelegateHeight := suite.Ctx.BlockHeight()
	suite.Ctx.Logger().Info("the ctx block height is:", "height", UndelegateHeight)

	// test complete Undelegation
	completeBlockNumber := UndelegateHeight + int64(delegationtype.CanUndelegationDelayHeight)
	suite.Ctx = suite.Ctx.WithBlockHeight(completeBlockNumber)
	suite.App.DelegationKeeper.EndBlock(suite.Ctx, abci.RequestEndBlock{})

	// check state
	stakerID, assetID := types.GetStakeIDAndAssetID(delegationEvent.ClientChainLzID, delegationEvent.StakerAddress, delegationEvent.AssetsAddress)
	restakerState, err := suite.App.AssetsKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerAssetInfo{
		TotalDepositAmount:  suite.depositAmount,
		WithdrawableAmount:  suite.depositAmount,
		WaitUnbondingAmount: sdkmath.NewInt(0),
	}, *restakerState)

	operatorState, err := suite.App.AssetsKeeper.GetOperatorSpecifiedAssetInfo(suite.Ctx, delegationEvent.OperatorAddress, assetID)
	suite.NoError(err)
	suite.Equal(types.OperatorAssetInfo{
		TotalAmount:             sdkmath.NewInt(0),
		OperatorAmount:          sdkmath.NewInt(0),
		WaitUnbondingAmount:     sdkmath.NewInt(0),
		OperatorUnbondingAmount: sdkmath.NewInt(0),
		TotalShare:              sdkmath.LegacyNewDec(0),
		OperatorShare:           sdkmath.LegacyNewDec(0),
	}, *operatorState)

	specifiedDelegationAmount, err := suite.App.DelegationKeeper.GetSingleDelegationInfo(suite.Ctx, stakerID, assetID, delegationEvent.OperatorAddress.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		UndelegatableShare:     sdkmath.LegacyNewDec(0),
		WaitUndelegationAmount: sdkmath.NewInt(0),
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.App.DelegationKeeper.StakerDelegatedTotalAmount(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(0), totalDelegationAmount)

	records, err := suite.App.DelegationKeeper.GetStakerUndelegationRecords(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(0, len(records))

	waitUndelegationRecords, err := suite.App.DelegationKeeper.GetWaitCompleteUndelegationRecords(suite.Ctx, uint64(completeBlockNumber))
	suite.NoError(err)
	suite.Equal(0, len(waitUndelegationRecords))
}
