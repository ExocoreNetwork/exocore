package keeper_test

import (
	"cosmossdk.io/math"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	operatortype "github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (suite *AVSTestSuite) TestAVS() {
	avsName, avsAddres, slashAddress := "avsTest", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutash"
	avsOwnerAddress := []string{"exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj1", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj2"}
	assetID := []string{"11", "22", "33"}
	avs := &types.AVSInfo{
		Name:               avsName,
		AvsAddress:         avsAddres,
		SlashAddr:          slashAddress,
		AvsOwnerAddress:    avsOwnerAddress,
		AssetId:            assetID,
		AvsUnbondingPeriod: uint32(7),
		MinSelfDelegation:  sdk.NewIntFromUint64(10),
		OperatorAddress:    nil,
		EpochIdentifier:    epochstypes.DayEpochID,
		StartingEpoch:      1,
	}

	err := suite.App.AVSManagerKeeper.SetAVSInfo(suite.Ctx, avs)
	suite.NoError(err)

	info, err := suite.App.AVSManagerKeeper.GetAVSInfo(suite.Ctx, avsAddres)

	suite.NoError(err)
	suite.Equal(avsAddres, info.GetInfo().AvsAddress)

	var avsList []types.AVSInfo
	suite.App.AVSManagerKeeper.IterateAVSInfo(suite.Ctx, func(_ int64, epochEndAVSInfo types.AVSInfo) (stop bool) {
		avsList = append(avsList, epochEndAVSInfo)
		return false
	})
	suite.Equal(len(avsList), 1)
}

func (suite *AVSTestSuite) TestAVSInfoUpdate_Register() {
	avsName, avsAddres, slashAddress, rewardAddress := "avsTest", "exo18cggcpvwspnd5c6ny8wrqxpffj5zmhklprtnph", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"
	avsOwnerAddress := []string{"exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj1", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj2"}
	assetID := []string{"11", "22", "33"}

	avsParams := &avstypes.AVSRegisterOrDeregisterParams{
		AvsName:            avsName,
		AvsAddress:         avsAddres,
		Action:             avstypes.RegisterAction,
		RewardContractAddr: rewardAddress,
		AvsOwnerAddress:    avsOwnerAddress,
		AssetID:            assetID,
		MinSelfDelegation:  uint64(10),
		UnbondingPeriod:    uint64(7),
		SlashContractAddr:  slashAddress,
		EpochIdentifier:    epochstypes.DayEpochID,
		OperatorAddress:    nil,
	}

	err := suite.App.AVSManagerKeeper.AVSInfoUpdate(suite.Ctx, avsParams)
	suite.NoError(err)

	info, err := suite.App.AVSManagerKeeper.GetAVSInfo(suite.Ctx, avsAddres)

	suite.NoError(err)
	suite.Equal(avsAddres, info.GetInfo().AvsAddress)

	err = suite.App.AVSManagerKeeper.AVSInfoUpdate(suite.Ctx, avsParams)
	suite.Error(err)
	suite.Contains(err.Error(), types.ErrAlreadyRegistered.Error())
}

func (suite *AVSTestSuite) TestAVSInfoUpdate_DeRegister() {
	// Test case setup
	avsName, avsAddres, slashAddress := "avsTest", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutash"
	avsOwnerAddress := []string{"exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj1", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj2"}
	assetID := []string{"11", "22", "33", "44", "55"} // Multiple assets

	avsParams := &avstypes.AVSRegisterOrDeregisterParams{
		AvsName:           avsName,
		AvsAddress:        avsAddres,
		Action:            avstypes.DeRegisterAction,
		AvsOwnerAddress:   avsOwnerAddress,
		AssetID:           assetID,
		MinSelfDelegation: uint64(10),
		UnbondingPeriod:   uint64(7),
		SlashContractAddr: slashAddress,
		EpochIdentifier:   epochstypes.DayEpochID,
		OperatorAddress:   nil,
	}

	err := suite.App.AVSManagerKeeper.AVSInfoUpdate(suite.Ctx, avsParams)
	suite.Error(err)
	suite.Contains(err.Error(), types.ErrUnregisterNonExistent.Error())

	avsParams.Action = avstypes.RegisterAction
	err = suite.App.AVSManagerKeeper.AVSInfoUpdate(suite.Ctx, avsParams)
	suite.NoError(err)
	info, err := suite.App.AVSManagerKeeper.GetAVSInfo(suite.Ctx, avsAddres)
	suite.Equal(avsAddres, info.GetInfo().AvsAddress)

	avsParams.Action = avstypes.DeRegisterAction
	avsParams.CallerAddress = "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr"
	err = suite.App.AVSManagerKeeper.AVSInfoUpdate(suite.Ctx, avsParams)
	suite.NoError(err)
	info, err = suite.App.AVSManagerKeeper.GetAVSInfo(suite.Ctx, avsAddres)
	suite.Error(err)
	suite.Contains(err.Error(), types.ErrNoKeyInTheStore.Error())
}

func (suite *AVSTestSuite) TestAVSInfoUpdateWithOperator_Register() {
	avsAddres, OperatorAddress := "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo19get9l6tj7pvn9xt83m7twpmup3usjvydpfscg"

	operatorParams := &avstypes.OperatorOptParams{
		AvsAddress:      avsAddres,
		Action:          avstypes.RegisterAction,
		OperatorAddress: OperatorAddress,
	}
	//  operator Not Exist
	err := suite.App.AVSManagerKeeper.AVSInfoUpdateWithOperator(suite.Ctx, operatorParams)
	suite.Error(err)
	suite.Contains(err.Error(), delegationtypes.ErrOperatorNotExist.Error())

	// register operator but avs not register
	info := &operatortype.OperatorInfo{
		EarningsAddr:     suite.AccAddress.String(),
		ApproveAddr:      "",
		OperatorMetaInfo: "test operator",
		ClientChainEarningsAddr: &operatortype.ClientChainEarningAddrList{
			EarningInfoList: []*operatortype.ClientChainEarningAddrInfo{
				{101, "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},
			},
		},
		Commission: stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
	}
	err = suite.App.OperatorKeeper.SetOperatorInfo(suite.Ctx, suite.AccAddress.String(), info)
	suite.NoError(err)
	operatorParams.OperatorAddress = info.EarningsAddr
	err = suite.App.AVSManagerKeeper.AVSInfoUpdateWithOperator(suite.Ctx, operatorParams)
	suite.Error(err)
	suite.Contains(err.Error(), types.ErrNoKeyInTheStore.Error())

	// register avs
	avsName, avsAddres, slashAddress := "avsTest", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutash"
	avsOwnerAddress := []string{"exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj1", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj2"}
	assetID := []string{"11", "22", "33"}

	avsParams := &avstypes.AVSRegisterOrDeregisterParams{
		AvsName:           avsName,
		AvsAddress:        avsAddres,
		Action:            avstypes.RegisterAction,
		AvsOwnerAddress:   avsOwnerAddress,
		AssetID:           assetID,
		MinSelfDelegation: uint64(10),
		UnbondingPeriod:   uint64(7),
		SlashContractAddr: slashAddress,
		EpochIdentifier:   epochstypes.DayEpochID,
		OperatorAddress:   nil,
	}

	err = suite.App.AVSManagerKeeper.AVSInfoUpdate(suite.Ctx, avsParams)
	suite.NoError(err)

	operatorParams.AvsAddress = avsAddres
	err = suite.App.AVSManagerKeeper.AVSInfoUpdateWithOperator(suite.Ctx, operatorParams)
	suite.NoError(err)
	// duplicate register operator
	err = suite.App.AVSManagerKeeper.AVSInfoUpdateWithOperator(suite.Ctx, operatorParams)
	suite.Error(err)
	suite.Contains(err.Error(), types.ErrAlreadyRegistered.Error())
	// deregister operator
	operatorParams.Action = avstypes.DeRegisterAction
	err = suite.App.AVSManagerKeeper.AVSInfoUpdateWithOperator(suite.Ctx, operatorParams)
	suite.NoError(err)

	// duplicate deregister operator
	operatorParams.Action = avstypes.DeRegisterAction
	err = suite.App.AVSManagerKeeper.AVSInfoUpdateWithOperator(suite.Ctx, operatorParams)
	suite.Error(err)
	suite.Contains(err.Error(), types.ErrUnregisterNonExistent.Error())
}
