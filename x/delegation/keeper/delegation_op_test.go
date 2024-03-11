package keeper_test

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/x/assets/types"
	keeper2 "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	types2 "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func (suite *DelegationTestSuite) TestDelegateTo() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzID := uint64(101)

	depositEvent := &keeper.DepositParams{
		ClientChainLzID: clientChainLzID,
		Action:          types.Deposit,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(100),
	}
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.App.DepositKeeper.Deposit(suite.Ctx, depositEvent)
	suite.NoError(err)

	opAccAddr, err := sdk.AccAddressFromBech32("exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr")
	suite.NoError(err)
	delegationParams := &keeper2.DelegationOrUndelegationParams{
		ClientChainLzID: clientChainLzID,
		Action:          types.DelegateTo,
		AssetsAddress:   usdtAddress[:],
		OperatorAddress: opAccAddr,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(50),
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, delegationParams)
	suite.EqualError(err, errorsmod.Wrap(delegationtype.ErrOperatorNotExist, fmt.Sprintf("input opreatorAddr is:%s", delegationParams.OperatorAddress)).Error())

	registerReq := &types2.RegisterOperatorReq{
		FromAddress: opAccAddr.String(),
		Info: &types2.OperatorInfo{
			EarningsAddr: opAccAddr.String(),
		},
	}
	_, err = suite.App.OperatorKeeper.RegisterOperator(suite.Ctx, registerReq)
	suite.NoError(err)

	err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, delegationParams)
	suite.NoError(err)

	// check delegation states
	stakerID, assetID := types.GetStakeIDAndAssetID(delegationParams.ClientChainLzID, delegationParams.StakerAddress, delegationParams.AssetsAddress)
	restakerState, err := suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetInfo{
		TotalDepositAmount:  depositEvent.OpAmount,
		WithdrawableAmount:  depositEvent.OpAmount.Sub(delegationParams.OpAmount),
		WaitUnbondingAmount: sdkmath.NewInt(0),
	}, *restakerState)

	operatorState, err := suite.App.StakingAssetsManageKeeper.GetOperatorSpecifiedAssetInfo(suite.Ctx, opAccAddr, assetID)
	suite.NoError(err)
	suite.Equal(types.OperatorSingleAssetInfo{
		TotalAmount:                        delegationParams.OpAmount,
		OperatorOwnAmount:                  sdkmath.NewInt(0),
		WaitUnbondingAmount:                sdkmath.NewInt(0),
		OperatorUnbondingAmount:            sdkmath.NewInt(0),
		OperatorUnbondableAmountAfterSlash: sdkmath.NewInt(0),
	}, *operatorState)

	specifiedDelegationAmount, err := suite.App.DelegationKeeper.GetSingleDelegationInfo(suite.Ctx, stakerID, assetID, opAccAddr.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		UndelegatableAmount:     delegationParams.OpAmount,
		WaitUndelegationAmount:  sdkmath.NewInt(0),
		UndelegatableAfterSlash: sdkmath.NewInt(0),
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.App.DelegationKeeper.GetStakerDelegationTotalAmount(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(delegationParams.OpAmount, totalDelegationAmount)
}

func (suite *DelegationTestSuite) TestUndelegateFrom() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzID := uint64(101)

	depositEvent := &keeper.DepositParams{
		ClientChainLzID: clientChainLzID,
		Action:          types.Deposit,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(100),
	}
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.App.DepositKeeper.Deposit(suite.Ctx, depositEvent)
	suite.NoError(err)

	opAccAddr, err := sdk.AccAddressFromBech32("exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr")
	suite.NoError(err)
	delegationEvent := &keeper2.DelegationOrUndelegationParams{
		ClientChainLzID: clientChainLzID,
		Action:          types.DelegateTo,
		AssetsAddress:   usdtAddress[:],
		OperatorAddress: opAccAddr,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(50),
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	registerReq := &types2.RegisterOperatorReq{
		FromAddress: opAccAddr.String(),
		Info: &types2.OperatorInfo{
			EarningsAddr: opAccAddr.String(),
		},
	}
	_, err = suite.App.OperatorKeeper.RegisterOperator(suite.Ctx, registerReq)
	suite.NoError(err)

	err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, delegationEvent)
	suite.NoError(err)

	// test Undelegation
	delegationEvent.LzNonce = 1
	err = suite.App.DelegationKeeper.UndelegateFrom(suite.Ctx, delegationEvent)
	suite.NoError(err)

	// check state
	stakerID, assetID := types.GetStakeIDAndAssetID(delegationEvent.ClientChainLzID, delegationEvent.StakerAddress, delegationEvent.AssetsAddress)
	restakerState, err := suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetInfo{
		TotalDepositAmount:  depositEvent.OpAmount,
		WithdrawableAmount:  depositEvent.OpAmount.Sub(delegationEvent.OpAmount),
		WaitUnbondingAmount: delegationEvent.OpAmount,
	}, *restakerState)

	operatorState, err := suite.App.StakingAssetsManageKeeper.GetOperatorSpecifiedAssetInfo(suite.Ctx, opAccAddr, assetID)
	suite.NoError(err)
	suite.Equal(types.OperatorSingleAssetInfo{
		TotalAmount:                        sdkmath.NewInt(0),
		OperatorOwnAmount:                  sdkmath.NewInt(0),
		WaitUnbondingAmount:                delegationEvent.OpAmount,
		OperatorUnbondingAmount:            sdkmath.NewInt(0),
		OperatorUnbondableAmountAfterSlash: sdkmath.NewInt(0),
	}, *operatorState)

	specifiedDelegationAmount, err := suite.App.DelegationKeeper.GetSingleDelegationInfo(suite.Ctx, stakerID, assetID, opAccAddr.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		UndelegatableAmount:     sdkmath.NewInt(0),
		WaitUndelegationAmount:  delegationEvent.OpAmount,
		UndelegatableAfterSlash: delegationEvent.OpAmount,
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.App.DelegationKeeper.GetStakerDelegationTotalAmount(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(0), totalDelegationAmount)

	records, err := suite.App.DelegationKeeper.GetStakerUndelegationRecords(suite.Ctx, stakerID, assetID, keeper2.PendingRecords)
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
		ActualCompletedAmount: sdkmath.NewInt(0),
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
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzID := uint64(101)

	depositEvent := &keeper.DepositParams{
		ClientChainLzID: clientChainLzID,
		Action:          types.Deposit,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(100),
	}
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.App.DepositKeeper.Deposit(suite.Ctx, depositEvent)
	suite.NoError(err)

	opAccAddr, err := sdk.AccAddressFromBech32("exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr")
	suite.NoError(err)
	delegationEvent := &keeper2.DelegationOrUndelegationParams{
		ClientChainLzID: clientChainLzID,
		Action:          types.DelegateTo,
		AssetsAddress:   usdtAddress[:],
		OperatorAddress: opAccAddr,
		StakerAddress:   suite.Address[:],
		OpAmount:        sdkmath.NewInt(50),
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	registerReq := &types2.RegisterOperatorReq{
		FromAddress: opAccAddr.String(),
		Info: &types2.OperatorInfo{
			EarningsAddr: opAccAddr.String(),
		},
	}
	_, err = suite.App.OperatorKeeper.RegisterOperator(suite.Ctx, registerReq)
	suite.NoError(err)

	err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, delegationEvent)
	suite.NoError(err)

	delegationEvent.LzNonce = 1
	err = suite.App.DelegationKeeper.UndelegateFrom(suite.Ctx, delegationEvent)
	suite.NoError(err)
	UndelegateHeight := suite.Ctx.BlockHeight()
	suite.Ctx.Logger().Info("the ctx block height is:", "height", UndelegateHeight)

	// test complete Undelegation
	completeBlockNumber := UndelegateHeight + int64(delegationtype.CanUndelegationDelayHeight)
	suite.Ctx = suite.Ctx.WithBlockHeight(completeBlockNumber)
	suite.App.DelegationKeeper.EndBlock(suite.Ctx, abci.RequestEndBlock{})

	// check state
	stakerID, assetID := types.GetStakeIDAndAssetID(delegationEvent.ClientChainLzID, delegationEvent.StakerAddress, delegationEvent.AssetsAddress)
	restakerState, err := suite.App.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetInfo{
		TotalDepositAmount:  depositEvent.OpAmount,
		WithdrawableAmount:  depositEvent.OpAmount,
		WaitUnbondingAmount: sdkmath.NewInt(0),
	}, *restakerState)

	operatorState, err := suite.App.StakingAssetsManageKeeper.GetOperatorSpecifiedAssetInfo(suite.Ctx, opAccAddr, assetID)
	suite.NoError(err)
	suite.Equal(types.OperatorSingleAssetInfo{
		TotalAmount:                        sdkmath.NewInt(0),
		OperatorOwnAmount:                  sdkmath.NewInt(0),
		WaitUnbondingAmount:                sdkmath.NewInt(0),
		OperatorUnbondingAmount:            sdkmath.NewInt(0),
		OperatorUnbondableAmountAfterSlash: sdkmath.NewInt(0),
	}, *operatorState)

	specifiedDelegationAmount, err := suite.App.DelegationKeeper.GetSingleDelegationInfo(suite.Ctx, stakerID, assetID, opAccAddr.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		UndelegatableAmount:     sdkmath.NewInt(0),
		WaitUndelegationAmount:  sdkmath.NewInt(0),
		UndelegatableAfterSlash: sdkmath.NewInt(0),
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.App.DelegationKeeper.GetStakerDelegationTotalAmount(suite.Ctx, stakerID, assetID)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(0), totalDelegationAmount)

	records, err := suite.App.DelegationKeeper.GetStakerUndelegationRecords(suite.Ctx, stakerID, assetID, keeper2.CompletedRecords)
	suite.NoError(err)
	suite.Equal(1, len(records))
	UndelegationRecord := &delegationtype.UndelegationRecord{
		StakerID:              stakerID,
		AssetID:               assetID,
		OperatorAddr:          delegationEvent.OperatorAddress.String(),
		TxHash:                delegationEvent.TxHash.String(),
		IsPending:             false,
		BlockNumber:           uint64(UndelegateHeight),
		LzTxNonce:             delegationEvent.LzNonce,
		Amount:                delegationEvent.OpAmount,
		ActualCompletedAmount: delegationEvent.OpAmount,
		CompleteBlockNumber:   uint64(completeBlockNumber),
	}
	suite.Equal(UndelegationRecord, records[0])

	waitUndelegationRecords, err := suite.App.DelegationKeeper.GetWaitCompleteUndelegationRecords(suite.Ctx, UndelegationRecord.CompleteBlockNumber)
	suite.NoError(err)
	suite.Equal(1, len(waitUndelegationRecords))
	suite.Equal(UndelegationRecord, waitUndelegationRecords[0])
}
