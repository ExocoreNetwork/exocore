package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v14/crypto/ethsecp256k1"
	"github.com/evmos/evmos/v14/testutil"
	utiltx "github.com/evmos/evmos/v14/testutil/tx"
	feemarkettypes "github.com/evmos/evmos/v14/x/feemarket/types"
	"github.com/exocore/app"
	"github.com/exocore/utils"
	exoslashkeeper "github.com/exocore/x/exoslash/keeper"
<<<<<<< HEAD
	types2 "github.com/exocore/x/exoslash/types"
=======
>>>>>>> a715c30 (add some test and fix bugs)
	"github.com/exocore/x/restaking_assets_manage/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx            sdk.Context
	app            *app.ExocoreApp
	address        common.Address
	signer         keyring.Signer
	accAddress     sdk.AccAddress
	exoSlashKeeper exoslashkeeper.Keeper
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	suite.Run(t, s)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest(suite.T())
}
func (suite *KeeperTestSuite) TestSlash() {
	usdtAddress := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	opAccAddr, _ := sdk.AccAddressFromBech32("evmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3h6cprl")
	stakerAddress := common.HexToAddress("0x4675C7e5BaAFBFFbca748158bEcBA61ef3b0a263")
	middlewareContractAddress := common.HexToAddress("0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5")
	slashEvent := &exoslashkeeper.SlashParams{
		ClientChainLzId:           3,
		Action:                    types.Slash,
		AssetsAddress:             usdtAddress.Bytes(),
		OperatorAddress:           opAccAddr,
		StakerAddress:             stakerAddress.Bytes(),
		OpAmount:                  sdkmath.NewInt(200),
		MiddlewareContractAddress: middlewareContractAddress.Bytes(),
		Proportion:                sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(3)),
		Proof:                     nil,
	}
	err := suite.app.ExoslashKeeper.Slash(suite.ctx, slashEvent)

	suite.ErrorContains(err, types2.ErrSlashAssetNotExist.Error())

	assets, err := suite.app.StakingAssetsManageKeeper.GetAllStakingAssetsInfo(suite.ctx)
	suite.NoError(err)
	suite.app.Logger().Info("the assets is:", "assets", assets)

	//test the normal case
	slashEvent.AssetsAddress = stakerAddress[:]
	err = suite.app.ExoslashKeeper.Slash(suite.ctx, slashEvent)
	suite.NoError(err)

	//check state after slash
	stakerId, assetId := types.GetStakeIDAndAssetId(slashEvent.ClientChainLzId, slashEvent.StakerAddress, slashEvent.AssetsAddress)
	info, err := suite.app.StakingAssetsManageKeeper.GetStakerSpecifiedAssetInfo(suite.ctx, stakerId, assetId)
	suite.NoError(err)
	suite.Equal(types.StakerSingleAssetOrChangeInfo{
		TotalDepositAmountOrWantChangeValue:     slashEvent.OpAmount,
		CanWithdrawAmountOrWantChangeValue:      slashEvent.OpAmount,
		WaitUnDelegationAmountOrWantChangeValue: sdkmath.NewInt(0),
	}, *info)

	assetInfo, err := suite.app.StakingAssetsManageKeeper.GetStakingAssetInfo(suite.ctx, assetId)
	suite.NoError(err)
	suite.Equal(slashEvent.OpAmount, assetInfo.StakingTotalAmount)
}

func (suite *KeeperTestSuite) DoSetupTest(t require.TestingT) {
	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = utiltx.NewSigner(priv)

	//accAddress
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}
	rand.Read(pub.Key)
	suite.accAddress = sdk.AccAddress(pub.Address())

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(t, err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())

	chainID := utils.TestnetChainID + "-1"
	suite.app = app.Setup(false, feemarkettypes.DefaultGenesisState(), chainID, false)
	header := testutil.NewHeader(
		1, time.Now().UTC(), chainID, consAddress, nil, nil,
	)
	suite.ctx = suite.app.BaseApp.NewContext(false, header)
}
