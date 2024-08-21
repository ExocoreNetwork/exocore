package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	assetskeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	avskeeper "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	operatorTypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

func (suite *AVSTestSuite) prepareOperator() {
	opAccAddr, err := sdk.AccAddressFromBech32("exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr")
	suite.operatorAddr = opAccAddr
	suite.NoError(err)
	// register operator
	registerReq := &operatorTypes.RegisterOperatorReq{
		FromAddress: suite.operatorAddr.String(),
		Info: &operatorTypes.OperatorInfo{
			EarningsAddr: suite.operatorAddr.String(),
		},
	}
	_, err = s.OperatorMsgServer.RegisterOperator(s.Ctx, registerReq)
	suite.NoError(err)
}

func (suite *AVSTestSuite) prepareDeposit(assetAddr common.Address, amount sdkmath.Int) {
	clientChainLzID := uint64(101)
	suite.avsAddr = common.BytesToAddress([]byte("avsTestAddr")).String()
	suite.assetAddr = assetAddr
	suite.assetDecimal = 6
	suite.clientChainLzID = clientChainLzID
	suite.depositAmount = amount
	suite.updatedAmountForOptIn = sdkmath.NewInt(20)
	suite.stakerID, suite.assetID = assetstypes.GetStakeIDAndAssetID(suite.clientChainLzID, suite.Address[:], suite.assetAddr[:])
	// staking assets
	depositParam := &assetskeeper.DepositWithdrawParams{
		ClientChainLzID: suite.clientChainLzID,
		Action:          assetstypes.Deposit,
		StakerAddress:   suite.Address[:],
		OpAmount:        suite.depositAmount,
		AssetsAddress:   assetAddr[:],
	}
	err := suite.App.AssetsKeeper.PerformDepositOrWithdraw(suite.Ctx, depositParam)
	suite.NoError(err)
}

func (suite *AVSTestSuite) prepareDelegation(isDelegation bool, assetAddr common.Address, amount sdkmath.Int) {
	suite.delegationAmount = amount
	param := &delegationtype.DelegationOrUndelegationParams{
		ClientChainID:   suite.clientChainLzID,
		AssetsAddress:   assetAddr[:],
		OperatorAddress: suite.operatorAddr,
		StakerAddress:   suite.Address[:],
		OpAmount:        amount,
		LzNonce:         0,
		TxHash:          common.HexToHash("0x24c4a315d757249c12a7a1d7b6fb96261d49deee26f06a3e1787d008b445c3ac"),
	}
	var err error
	if isDelegation {
		err = suite.App.DelegationKeeper.DelegateTo(suite.Ctx, param)
	} else {
		err = suite.App.DelegationKeeper.UndelegateFrom(suite.Ctx, param)
	}
	suite.NoError(err)
}
func (suite *AVSTestSuite) prepareAvs(assetIDs []string) {
	err := suite.App.AVSManagerKeeper.AVSInfoUpdate(suite.Ctx, &avstypes.AVSRegisterOrDeregisterParams{
		Action:          avskeeper.RegisterAction,
		EpochIdentifier: epochstypes.HourEpochID,
		AvsAddress:      suite.avsAddr,
		AssetID:         assetIDs,
	})
	suite.NoError(err)
}
func (suite *AVSTestSuite) prepareOptIn() {
	err := suite.App.OperatorKeeper.OptIn(suite.Ctx, suite.operatorAddr, suite.avsAddr)
	suite.NoError(err)
}

func (suite *AVSTestSuite) prepare() {
	usdtAddress := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	depositAmount := sdkmath.NewInt(100)
	delegationAmount := sdkmath.NewInt(50)
	suite.prepareOperator()
	suite.prepareDeposit(usdtAddress, depositAmount)
	suite.prepareDelegation(true, usdtAddress, delegationAmount)
	suite.prepareAvs([]string{"0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"})
	suite.prepareOptIn()
	suite.CommitAfter(time.Hour*1 + time.Nanosecond)
	suite.CommitAfter(time.Hour*1 + time.Nanosecond)
	suite.CommitAfter(time.Hour*1 + time.Nanosecond)
}

func (suite *AVSTestSuite) TestOptInList() {
	suite.prepare()
	operatorList, err := suite.App.OperatorKeeper.GetOptedInOperatorListByAVS(suite.Ctx, suite.avsAddr)
	suite.NoError(err)
	suite.Contains(operatorList, suite.operatorAddr.String())

	avsList, err := suite.App.OperatorKeeper.GetOptedInAVSForOperator(suite.Ctx, suite.operatorAddr.String())
	suite.NoError(err)

	suite.Contains(avsList, suite.avsAddr)
}
func (suite *AVSTestSuite) TestAVSUSDValue() {
	suite.prepare()
	avsUSDValue, err := suite.App.OperatorKeeper.GetAVSUSDValue(suite.Ctx, suite.avsAddr)
	suite.NoError(err)
	optedUSDValues, err := suite.App.OperatorKeeper.GetOperatorOptedUSDValue(suite.Ctx, suite.avsAddr, suite.operatorAddr.String())
	suite.NoError(err)
	fmt.Println(avsUSDValue)
	fmt.Println(optedUSDValues)

}
