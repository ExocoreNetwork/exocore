package keeper_test

import (
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	operatorTypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (suite *AVSTestSuite) TestAVS() {
	avsName, avsAddres, slashAddress := "avsTest", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutash"
	avsOwnerAddress := []string{"exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj1", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj2"}
	assetID := []string{"11", "22", "33"}
	avs := &types.AVSInfo{
		Name:                avsName,
		AvsAddress:          avsAddres,
		SlashAddr:           slashAddress,
		AvsOwnerAddress:     avsOwnerAddress,
		AssetIDs:            assetID,
		AvsUnbondingPeriod:  7,
		MinSelfDelegation:   10,
		EpochIdentifier:     epochstypes.DayEpochID,
		StartingEpoch:       1,
		MinOptInOperators:   100,
		MinTotalStakeAmount: 1000,
		AvsSlash:            sdk.MustNewDecFromStr("0.001"),
		AvsReward:           sdk.MustNewDecFromStr("0.002"),
		TaskAddr:            "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr",
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
	suite.CommitAfter(48*time.Hour + time.Nanosecond)
	// commit will run the EndBlockers for the current block, call app.Commit
	// and then run the BeginBlockers for the next block with the new time.
	// during the BeginBlocker, the epoch will be incremented.
	epoch, found := suite.App.EpochsKeeper.GetEpochInfo(suite.Ctx, epochstypes.DayEpochID)
	suite.Equal(found, true)
	suite.Equal(epoch.CurrentEpoch, int64(2))
	suite.CommitAfter(48*time.Hour + time.Nanosecond)

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

	opAccAddr, err := sdk.AccAddressFromBech32(OperatorAddress)
	suite.NoError(err)
	operatorParams := &avstypes.OperatorOptParams{
		AvsAddress:      avsAddres,
		Action:          avstypes.RegisterAction,
		OperatorAddress: OperatorAddress,
	}
	//  operator Not Exist
	err = suite.App.AVSManagerKeeper.OperatorOptAction(suite.Ctx, operatorParams)
	suite.Error(err)
	suite.Contains(err.Error(), delegationtypes.ErrOperatorNotExist.Error())

	// register operator but avs not register
	// register operator
	registerReq := &operatorTypes.RegisterOperatorReq{
		FromAddress: opAccAddr.String(),
		Info: &operatorTypes.OperatorInfo{
			EarningsAddr: opAccAddr.String(),
		},
	}
	_, err = suite.App.OperatorKeeper.RegisterOperator(suite.Ctx, registerReq)
	suite.NoError(err)
	suite.TestAVS()

	err = suite.App.AVSManagerKeeper.OperatorOptAction(suite.Ctx, operatorParams)
	suite.NoError(err)

}
