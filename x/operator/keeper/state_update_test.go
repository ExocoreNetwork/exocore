package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	delegationKeeper "github.com/exocore/x/delegation/keeper"
	"github.com/exocore/x/deposit/keeper"
	operatorKeeper "github.com/exocore/x/operator/keeper"
	operatorTypes "github.com/exocore/x/operator/types"
	restakingTypes "github.com/exocore/x/restaking_assets_manage/types"
	"strings"
)

type StateForCheck struct {
	OptedInfo        *operatorTypes.OptedInfo
	AVSTotalShare    sdkmath.LegacyDec
	AVSOperatorShare sdkmath.LegacyDec
	AssetState       *operatorTypes.AssetOptedInState
	OperatorShare    sdkmath.LegacyDec
	StakerShare      sdkmath.LegacyDec
}

func (suite *KeeperTestSuite) prepare() {
	opAccAddr, err := sdk.AccAddressFromBech32("evmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3h6cprl")
	suite.NoError(err)
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzId := uint64(101)

	suite.avsAddr = "avsTestAddr"
	suite.operatorAddr = opAccAddr
	suite.assetAddr = usdtAddress
	suite.assetDecimal = 6
	suite.clientChainLzId = clientChainLzId
	suite.depositAmount = sdkmath.NewInt(100)
	suite.delegationAmount = sdkmath.NewInt(50)
	suite.updatedAmountForOptIn = sdkmath.NewInt(20)
	suite.stakerId, suite.assetId = restakingTypes.GetStakeIDAndAssetId(suite.clientChainLzId, suite.address[:], suite.assetAddr[:])

	//staking assets
	depositParam := &keeper.DepositParams{
		ClientChainLzId: suite.clientChainLzId,
		Action:          restakingTypes.Deposit,
		StakerAddress:   suite.address[:],
		OpAmount:        suite.depositAmount,
	}
	depositParam.AssetsAddress = suite.assetAddr[:]
	err = suite.app.DepositKeeper.Deposit(suite.ctx, depositParam)
	suite.NoError(err)

	//register operator
	registerReq := &operatorTypes.RegisterOperatorReq{
		FromAddress: suite.operatorAddr.String(),
		Info: &operatorTypes.OperatorInfo{
			EarningsAddr: suite.operatorAddr.String(),
		},
	}
	_, err = suite.app.OperatorKeeper.RegisterOperator(suite.ctx, registerReq)
	suite.NoError(err)

	//delegate to operator
	delegationParam := &delegationKeeper.DelegationOrUndelegationParams{
		ClientChainLzId: suite.clientChainLzId,
		Action:          restakingTypes.DelegateTo,
		AssetsAddress:   suite.assetAddr[:],
		OperatorAddress: suite.operatorAddr,
		StakerAddress:   suite.address[:],
		OpAmount:        suite.delegationAmount,
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	err = suite.app.DelegationKeeper.DelegateTo(suite.ctx, delegationParam)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) CheckState(expectedState *StateForCheck) {
	//check opted info
	optInfo, err := suite.app.OperatorKeeper.GetOptedInfo(suite.ctx, suite.operatorAddr.String(), suite.avsAddr)
	if expectedState.OptedInfo == nil {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(*expectedState.OptedInfo, *optInfo)
	}
	//check total USD value for AVS and operator
	value, err := suite.app.OperatorKeeper.GetAVSShare(suite.ctx, suite.avsAddr)
	if expectedState.AVSTotalShare.IsNil() {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(expectedState.AVSTotalShare, value)
	}

	value, err = suite.app.OperatorKeeper.GetOperatorShare(suite.ctx, suite.avsAddr, suite.operatorAddr.String())
	if expectedState.AVSOperatorShare.IsNil() {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(expectedState.AVSOperatorShare, value)
	}

	//check assets state for AVS and operator
	assetState, err := suite.app.OperatorKeeper.GetAssetState(suite.ctx, suite.assetId, suite.avsAddr, suite.operatorAddr.String())
	if expectedState.AssetState == nil {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(*expectedState.AssetState, *assetState)
	}

	//check asset USD share for staker and operator
	operatorShare, err := suite.app.OperatorKeeper.GetStakerShare(suite.ctx, suite.avsAddr, "", suite.operatorAddr.String())
	if expectedState.OperatorShare.IsNil() {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(expectedState.OperatorShare, operatorShare)
	}
	stakerShare, err := suite.app.OperatorKeeper.GetStakerShare(suite.ctx, suite.avsAddr, suite.stakerId, suite.operatorAddr.String())
	if expectedState.StakerShare.IsNil() {
		suite.True(strings.Contains(err.Error(), operatorTypes.ErrNoKeyInTheStore.Error()))
	} else {
		suite.NoError(err)
		suite.Equal(expectedState.StakerShare, stakerShare)
	}
}

func (suite *KeeperTestSuite) TestOptIn() {
	suite.prepare()
	err := suite.app.OperatorKeeper.OptIn(suite.ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	//check if the related state is correct
	price, decimal, err := suite.app.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.ctx, suite.assetId)
	share := operatorKeeper.CalculateShare(suite.delegationAmount, price, suite.assetDecimal, decimal)
	expectedState := &StateForCheck{
		OptedInfo: &operatorTypes.OptedInfo{
			OptedInHeight:  uint64(suite.ctx.BlockHeight()),
			OptedOutHeight: operatorTypes.DefaultOptedOutHeight,
		},
		AVSTotalShare:    share,
		AVSOperatorShare: share,
		AssetState: &operatorTypes.AssetOptedInState{
			Amount: suite.delegationAmount,
			Value:  share,
		},
		OperatorShare: sdkmath.LegacyDec{},
		StakerShare:   share,
	}
	suite.CheckState(expectedState)
}

func (suite *KeeperTestSuite) TestOptOut() {
	suite.prepare()
	err := suite.app.OperatorKeeper.OptOut(suite.ctx, suite.operatorAddr, suite.avsAddr)
	suite.EqualError(err, operatorTypes.ErrNotOptedIn.Error())

	err = suite.app.OperatorKeeper.OptIn(suite.ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	optInHeight := suite.ctx.BlockHeight()
	suite.NextBlock()

	err = suite.app.OperatorKeeper.OptOut(suite.ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)

	expectedState := &StateForCheck{
		OptedInfo: &operatorTypes.OptedInfo{
			OptedInHeight:  uint64(optInHeight),
			OptedOutHeight: uint64(suite.ctx.BlockHeight()),
		},
		AVSTotalShare:    sdkmath.LegacyNewDec(0),
		AVSOperatorShare: sdkmath.LegacyDec{},
		AssetState:       nil,
		OperatorShare:    sdkmath.LegacyDec{},
		StakerShare:      sdkmath.LegacyDec{},
	}
	suite.CheckState(expectedState)
}

func (suite *KeeperTestSuite) TestCalculateShare() {
	suite.prepare()
	price, decimal, err := suite.app.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.ctx, suite.assetId)
	suite.NoError(err)
	share := operatorKeeper.CalculateShare(suite.delegationAmount, price, suite.assetDecimal, decimal)
	suite.Equal(sdkmath.LegacyNewDecWithPrec(5000, int64(operatorTypes.USDValueDefaultDecimal)), share)
}

func (suite *KeeperTestSuite) TestUpdateOptedInAssetsState() {
	suite.prepare()
	err := suite.app.OperatorKeeper.OptIn(suite.ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	optInHeight := suite.ctx.BlockHeight()
	suite.NextBlock()

	err = suite.app.OperatorKeeper.UpdateOptedInAssetsState(suite.ctx, suite.stakerId, suite.assetId, suite.operatorAddr.String(), suite.updatedAmountForOptIn)
	suite.NoError(err)

	price, decimal, err := suite.app.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.ctx, suite.assetId)
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
		AssetState: &operatorTypes.AssetOptedInState{
			Amount: suite.delegationAmount.Add(suite.updatedAmountForOptIn),
			Value:  newShare,
		},
		OperatorShare: sdkmath.LegacyDec{},
		StakerShare:   newShare,
	}
	suite.CheckState(expectedState)
}

/*func (suite *KeeperTestSuite) TestSlash() {
	suite.prepare()
	err := suite.app.OperatorKeeper.OptIn(suite.ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
	optInHeight := suite.ctx.BlockHeight()
	suite.NextBlock()
}*/
