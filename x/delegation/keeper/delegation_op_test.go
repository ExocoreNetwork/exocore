package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	keeper2 "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/exocore/x/operator/types"
)

func (suite *KeeperTestSuite) TestDelegateTo() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzId := uint64(101)

	depositEvent := &keeper.DepositParams{
		ClientChainLzId: clientChainLzId,
		Action:          types.Deposit,
		StakerAddress:   suite.address[:],
		OpAmount:        sdkmath.NewInt(100),
	}
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.app.DepositKeeper.Deposit(suite.ctx, depositEvent)
	suite.NoError(err)

	opAccAddr, err := sdk.AccAddressFromBech32("evmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3h6cprl")
	suite.NoError(err)
	delegationParams := &keeper2.DelegationOrUndelegationParams{
		ClientChainLzId: clientChainLzId,
		Action:          types.DelegateTo,
		AssetsAddress:   usdtAddress[:],
		OperatorAddress: opAccAddr,
		StakerAddress:   suite.address[:],
		OpAmount:        sdkmath.NewInt(50),
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	err = suite.app.DelegationKeeper.DelegateTo(suite.ctx, delegationParams)
	suite.EqualError(err, delegationtype.ErrOperatorNotExist.Error())

	registerReq := &types2.RegisterOperatorReq{
		FromAddress: opAccAddr.String(),
		Info: &types2.OperatorInfo{
			EarningsAddr: opAccAddr.String(),
		},
	}
	_, err = suite.app.OperatorKeeper.RegisterOperator(suite.ctx, registerReq)
	suite.NoError(err)

	err = suite.app.DelegationKeeper.DelegateTo(suite.ctx, delegationParams)
	suite.NoError(err)

	// check delegation states
	stakerId, assetId := types.GetStakeIDAndAssetId(delegationParams.ClientChainLzId, delegationParams.StakerAddress, delegationParams.AssetsAddress)
	restakerState, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:  depositEvent.OpAmount,
		CanWithdrawAmountOrWantChangeValue:   depositEvent.OpAmount.Sub(delegationParams.OpAmount),
		WaitUnbondingAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *restakerState)

	operatorState, err := suite.app.StakingAssetsManageKeeper.GetOperatorSpecifiedAssetInfo(suite.ctx, opAccAddr, assetId)
	suite.NoError(err)
	suite.Equal(types.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue:         delegationParams.OpAmount,
		OperatorOwnAmountOrWantChangeValue:   sdkmath.NewInt(0),
		WaitUnbondingAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *operatorState)

	specifiedDelegationAmount, err := suite.app.DelegationKeeper.GetSingleDelegationInfo(suite.ctx, stakerId, assetId, opAccAddr.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		CanUndelegationAmount:  delegationParams.OpAmount,
		WaitUndelegationAmount: sdkmath.NewInt(0),
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.app.DelegationKeeper.GetStakerDelegationTotalAmount(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(delegationParams.OpAmount, totalDelegationAmount)
}

func (suite *KeeperTestSuite) TestUndelegateFrom() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzId := uint64(101)

	depositEvent := &keeper.DepositParams{
		ClientChainLzId: clientChainLzId,
		Action:          types.Deposit,
		StakerAddress:   suite.address[:],
		OpAmount:        sdkmath.NewInt(100),
	}
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.app.DepositKeeper.Deposit(suite.ctx, depositEvent)
	suite.NoError(err)

	opAccAddr, err := sdk.AccAddressFromBech32("evmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3h6cprl")
	suite.NoError(err)
	delegationEvent := &keeper2.DelegationOrUndelegationParams{
		ClientChainLzId: clientChainLzId,
		Action:          types.DelegateTo,
		AssetsAddress:   usdtAddress[:],
		OperatorAddress: opAccAddr,
		StakerAddress:   suite.address[:],
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
	_, err = suite.app.OperatorKeeper.RegisterOperator(suite.ctx, registerReq)
	suite.NoError(err)

	err = suite.app.DelegationKeeper.DelegateTo(suite.ctx, delegationEvent)
	suite.NoError(err)

	// test Undelegation
	delegationEvent.LzNonce = 1
	err = suite.app.DelegationKeeper.UndelegateFrom(suite.ctx, delegationEvent)
	suite.NoError(err)

	// check state
	stakerId, assetId := types.GetStakeIDAndAssetId(delegationEvent.ClientChainLzId, delegationEvent.StakerAddress, delegationEvent.AssetsAddress)
	restakerState, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:  depositEvent.OpAmount,
		CanWithdrawAmountOrWantChangeValue:   depositEvent.OpAmount.Sub(delegationEvent.OpAmount),
		WaitUnbondingAmountOrWantChangeValue: delegationEvent.OpAmount,
	}, *restakerState)

	operatorState, err := suite.app.StakingAssetsManageKeeper.GetOperatorSpecifiedAssetInfo(suite.ctx, opAccAddr, assetId)
	suite.NoError(err)
	suite.Equal(types.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue:         sdkmath.NewInt(0),
		OperatorOwnAmountOrWantChangeValue:   sdkmath.NewInt(0),
		WaitUnbondingAmountOrWantChangeValue: delegationEvent.OpAmount,
	}, *operatorState)

	specifiedDelegationAmount, err := suite.app.DelegationKeeper.GetSingleDelegationInfo(suite.ctx, stakerId, assetId, opAccAddr.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		CanUndelegationAmount:  sdkmath.NewInt(0),
		WaitUndelegationAmount: delegationEvent.OpAmount,
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.app.DelegationKeeper.GetStakerDelegationTotalAmount(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(delegationEvent.OpAmount, totalDelegationAmount)

	records, err := suite.app.DelegationKeeper.GetStakerUndelegationRecords(suite.ctx, stakerId, assetId, keeper2.PendingRecords)
	suite.NoError(err)
	suite.Equal(1, len(records))
	UndelegationRecord := &delegationtype.UndelegationRecord{
		StakerId:              stakerId,
		AssetId:               assetId,
		OperatorAddr:          delegationEvent.OperatorAddress.String(),
		TxHash:                delegationEvent.TxHash.String(),
		IsPending:             true,
		BlockNumber:           uint64(suite.ctx.BlockHeight()),
		LzTxNonce:             delegationEvent.LzNonce,
		Amount:                delegationEvent.OpAmount,
		ActualCompletedAmount: sdkmath.NewInt(0),
	}
	UndelegationRecord.CompleteBlockNumber = UndelegationRecord.BlockNumber + delegationtype.CanUndelegationDelayHeight
	suite.Equal(UndelegationRecord, records[0])

	suite.ctx.Logger().Info("the complete block number is:", "height", UndelegationRecord.CompleteBlockNumber)
	waitUndelegationRecords, err := suite.app.DelegationKeeper.GetWaitCompleteUndelegationRecords(suite.ctx, UndelegationRecord.CompleteBlockNumber)
	suite.NoError(err)
	suite.Equal(1, len(waitUndelegationRecords))
	suite.Equal(UndelegationRecord, waitUndelegationRecords[0])
}

func (suite *KeeperTestSuite) TestCompleteUndelegation() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzId := uint64(101)

	depositEvent := &keeper.DepositParams{
		ClientChainLzId: clientChainLzId,
		Action:          types.Deposit,
		StakerAddress:   suite.address[:],
		OpAmount:        sdkmath.NewInt(100),
	}
	depositEvent.AssetsAddress = usdtAddress[:]
	err := suite.app.DepositKeeper.Deposit(suite.ctx, depositEvent)
	suite.NoError(err)

	opAccAddr, err := sdk.AccAddressFromBech32("evmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3h6cprl")
	suite.NoError(err)
	delegationEvent := &keeper2.DelegationOrUndelegationParams{
		ClientChainLzId: clientChainLzId,
		Action:          types.DelegateTo,
		AssetsAddress:   usdtAddress[:],
		OperatorAddress: opAccAddr,
		StakerAddress:   suite.address[:],
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
	_, err = suite.app.OperatorKeeper.RegisterOperator(suite.ctx, registerReq)
	suite.NoError(err)

	err = suite.app.DelegationKeeper.DelegateTo(suite.ctx, delegationEvent)
	suite.NoError(err)

	delegationEvent.LzNonce = 1
	err = suite.app.DelegationKeeper.UndelegateFrom(suite.ctx, delegationEvent)
	suite.NoError(err)
	UndelegateHeight := suite.ctx.BlockHeight()
	suite.ctx.Logger().Info("the ctx block height is:", "height", UndelegateHeight)

	// test complete Undelegation
	completeBlockNumber := UndelegateHeight + int64(delegationtype.CanUndelegationDelayHeight)
	suite.ctx = suite.ctx.WithBlockHeight(completeBlockNumber)
	suite.app.DelegationKeeper.EndBlock(suite.ctx, abci.RequestEndBlock{})

	// check state
	stakerId, assetId := types.GetStakeIDAndAssetId(delegationEvent.ClientChainLzId, delegationEvent.StakerAddress, delegationEvent.AssetsAddress)
	restakerState, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:  depositEvent.OpAmount,
		CanWithdrawAmountOrWantChangeValue:   depositEvent.OpAmount,
		WaitUnbondingAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *restakerState)

	operatorState, err := suite.app.StakingAssetsManageKeeper.GetOperatorSpecifiedAssetInfo(suite.ctx, opAccAddr, assetId)
	suite.NoError(err)
	suite.Equal(types.OperatorSingleAssetOrChangeInfo{
		TotalAmountOrWantChangeValue:         sdkmath.NewInt(0),
		OperatorOwnAmountOrWantChangeValue:   sdkmath.NewInt(0),
		WaitUnbondingAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *operatorState)

	specifiedDelegationAmount, err := suite.app.DelegationKeeper.GetSingleDelegationInfo(suite.ctx, stakerId, assetId, opAccAddr.String())
	suite.NoError(err)
	suite.Equal(delegationtype.DelegationAmounts{
		CanUndelegationAmount:  sdkmath.NewInt(0),
		WaitUndelegationAmount: sdkmath.NewInt(0),
	}, *specifiedDelegationAmount)

	totalDelegationAmount, err := suite.app.DelegationKeeper.GetStakerDelegationTotalAmount(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(sdkmath.NewInt(0), totalDelegationAmount)

	records, err := suite.app.DelegationKeeper.GetStakerUndelegationRecords(suite.ctx, stakerId, assetId, keeper2.CompletedRecords)
	suite.NoError(err)
	suite.Equal(1, len(records))
	UndelegationRecord := &delegationtype.UndelegationRecord{
		StakerId:              stakerId,
		AssetId:               assetId,
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

	waitUndelegationRecords, err := suite.app.DelegationKeeper.GetWaitCompleteUndelegationRecords(suite.ctx, UndelegationRecord.CompleteBlockNumber)
	suite.NoError(err)
	suite.Equal(1, len(waitUndelegationRecords))
	suite.Equal(UndelegationRecord, waitUndelegationRecords[0])
}
