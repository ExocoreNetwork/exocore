package keeper_test

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	delegationKeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	"github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	operatorKeeper "github.com/ExocoreNetwork/exocore/x/operator/keeper"
	operatorTypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	restakingTypes "github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type StateForCheck struct {
	OptedInfo        *operatorTypes.OptedInfo
	AVSTotalShare    sdkmath.LegacyDec
	AVSOperatorShare sdkmath.LegacyDec
	AssetState       *operatorTypes.OptedInAssetState
	OperatorShare    sdkmath.LegacyDec
	StakerShare      sdkmath.LegacyDec
}

func (suite *OperatorTestSuite) prepare() {
	opAccAddr, err := sdk.AccAddressFromBech32("exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr")
	suite.NoError(err)
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzID := uint64(101)

	suite.avsAddr = "avsTestAddr"
	suite.operatorAddr = opAccAddr
	suite.assetAddr = usdtAddress
	suite.assetDecimal = 6
	suite.clientChainLzID = clientChainLzID
	suite.depositAmount = sdkmath.NewInt(100)
	suite.delegationAmount = sdkmath.NewInt(50)
	suite.updatedAmountForOptIn = sdkmath.NewInt(20)
	suite.stakerID, suite.assetID = restakingTypes.GetStakeIDAndAssetID(suite.clientChainLzID, suite.Address[:], suite.assetAddr[:])

	// staking assets
	depositParam := &keeper.DepositParams{
		ClientChainLzID: suite.clientChainLzID,
		Action:          restakingTypes.Deposit,
		StakerAddress:   suite.Address[:],
		OpAmount:        suite.depositAmount,
	}
	depositParam.AssetsAddress = suite.assetAddr[:]
	err = suite.App.DepositKeeper.Deposit(suite.Ctx, depositParam)
	suite.NoError(err)

	// register operator
	registerReq := &operatorTypes.RegisterOperatorReq{
		FromAddress: suite.operatorAddr.String(),
		Info: &operatorTypes.OperatorInfo{
			EarningsAddr: suite.operatorAddr.String(),
		},
	}
	_, err = suite.App.OperatorKeeper.RegisterOperator(suite.Ctx, registerReq)
	suite.NoError(err)

	// delegate to operator
	delegationParam := &delegationKeeper.DelegationOrUndelegationParams{
		ClientChainLzID: suite.clientChainLzID,
		Action:          restakingTypes.DelegateTo,
		AssetsAddress:   suite.assetAddr[:],
		OperatorAddress: suite.operatorAddr,
		StakerAddress:   suite.Address[:],
		OpAmount:        suite.delegationAmount,
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, delegationParam)
	suite.NoError(err)
}

func (suite *OperatorTestSuite) CheckState(expectedState *StateForCheck) {
	// check opted info
	optInfo, err := suite.App.OperatorKeeper.GetOptedInfo(suite.Ctx, suite.operatorAddr.String(), suite.avsAddr)
	if expectedState.OptedInfo == nil {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(*expectedState.OptedInfo, *optInfo)
	}
	// check total USD value for AVS and operator
	value, err := suite.App.OperatorKeeper.GetAVSShare(suite.Ctx, suite.avsAddr)
	if expectedState.AVSTotalShare.IsNil() {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(expectedState.AVSTotalShare, value)
	}

	value, err = suite.App.OperatorKeeper.GetOperatorShare(suite.Ctx, suite.avsAddr, suite.operatorAddr.String())
	if expectedState.AVSOperatorShare.IsNil() {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(expectedState.AVSOperatorShare, value)
	}

	// check assets state for AVS and operator
	assetState, err := suite.App.OperatorKeeper.GetAssetState(suite.Ctx, suite.assetID, suite.avsAddr, suite.operatorAddr.String())
	if expectedState.AssetState == nil {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(*expectedState.AssetState, *assetState)
	}

	// check asset USD share for staker and operator
	operatorShare, err := suite.App.OperatorKeeper.GetStakerShare(suite.Ctx, suite.avsAddr, "", suite.operatorAddr.String())
	if expectedState.OperatorShare.IsNil() {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(expectedState.OperatorShare, operatorShare)
	}
	stakerShare, err := suite.App.OperatorKeeper.GetStakerShare(suite.Ctx, suite.avsAddr, suite.stakerID, suite.operatorAddr.String())
	if expectedState.StakerShare.IsNil() {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(expectedState.StakerShare, stakerShare)
	}
}

func (suite *OperatorTestSuite) TestOptIn() {
	suite.prepare()
	err := suite.App.OperatorKeeper.OptIn(suite.Ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	// check if the related state is correct
	price, decimal, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
	share := operatorKeeper.CalculateShare(suite.delegationAmount, price, suite.assetDecimal, decimal)
	expectedState := &StateForCheck{
		OptedInfo: &operatorTypes.OptedInfo{
			OptedInHeight:  uint64(suite.Ctx.BlockHeight()),
			OptedOutHeight: operatorTypes.DefaultOptedOutHeight,
		},
		AVSTotalShare:    share,
		AVSOperatorShare: share,
		AssetState: &operatorTypes.OptedInAssetState{
			Amount: suite.delegationAmount,
			Value:  share,
		},
		OperatorShare: sdkmath.LegacyDec{},
		StakerShare:   share,
	}
	suite.CheckState(expectedState)
}

func (suite *OperatorTestSuite) TestOptOut() {
	suite.prepare()
	err := suite.App.OperatorKeeper.OptOut(suite.Ctx, suite.operatorAddr, suite.avsAddr)
	suite.EqualError(err, operatorTypes.ErrNotOptedIn.Error())

	err = suite.App.OperatorKeeper.OptIn(suite.Ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	optInHeight := suite.Ctx.BlockHeight()
	suite.NextBlock()

	err = suite.App.OperatorKeeper.OptOut(suite.Ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)

	expectedState := &StateForCheck{
		OptedInfo: &operatorTypes.OptedInfo{
			OptedInHeight:  uint64(optInHeight),
			OptedOutHeight: uint64(suite.Ctx.BlockHeight()),
		},
		AVSTotalShare:    sdkmath.LegacyNewDec(0),
		AVSOperatorShare: sdkmath.LegacyDec{},
		AssetState:       nil,
		OperatorShare:    sdkmath.LegacyDec{},
		StakerShare:      sdkmath.LegacyDec{},
	}
	suite.CheckState(expectedState)
}

func (suite *OperatorTestSuite) TestCalculateShare() {
	suite.prepare()
	price, decimal, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
	suite.NoError(err)
	share := operatorKeeper.CalculateShare(suite.delegationAmount, price, suite.assetDecimal, decimal)
	suite.Equal(sdkmath.LegacyNewDecWithPrec(5000, int64(operatorTypes.USDValueDefaultDecimal)), share)
}

func (suite *OperatorTestSuite) TestUpdateOptedInAssetsState() {
	suite.prepare()
	err := suite.App.OperatorKeeper.OptIn(suite.Ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	optInHeight := suite.Ctx.BlockHeight()
	suite.NextBlock()

	err = suite.App.OperatorKeeper.UpdateOptedInAssetsState(suite.Ctx, suite.stakerID, suite.assetID, suite.operatorAddr.String(), suite.updatedAmountForOptIn)
	suite.NoError(err)

	price, decimal, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
	oldShare := operatorKeeper.CalculateShare(suite.delegationAmount, price, suite.assetDecimal, decimal)
	addShare := operatorKeeper.CalculateShare(suite.updatedAmountForOptIn, price, suite.assetDecimal, decimal)
	newShare := oldShare.Add(addShare)

	expectedState := &StateForCheck{
		OptedInfo: &operatorTypes.OptedInfo{
			OptedInHeight:  uint64(optInHeight),
			OptedOutHeight: operatorTypes.DefaultOptedOutHeight,
		},
		AVSTotalShare:    newShare,
		AVSOperatorShare: newShare,
		AssetState: &operatorTypes.OptedInAssetState{
			Amount: suite.delegationAmount.Add(suite.updatedAmountForOptIn),
			Value:  newShare,
		},
		OperatorShare: sdkmath.LegacyDec{},
		StakerShare:   newShare,
	}
	suite.CheckState(expectedState)
}

func (suite *OperatorTestSuite) TestSlash() {
	suite.prepare()
	err := suite.App.OperatorKeeper.OptIn(suite.Ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	optInHeight := suite.Ctx.BlockHeight()

	// run to the block at specified height
	runToHeight := optInHeight + 10
	for i := optInHeight; i < runToHeight; i++ {
		suite.NextBlock()
	}
	suite.Equal(runToHeight, suite.Ctx.BlockHeight())
}
