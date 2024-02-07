package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/evmos/evmos/v14/x/evm/statedb"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmosapp "github.com/exocore/app"
)

var s *KeeperTestSuite

type KeeperTestSuite struct {
	suite.Suite

	ctx        sdk.Context
	app        *evmosapp.ExocoreApp
	address    common.Address
	accAddress sdk.AccAddress

	validators []stakingtypes.Validator
	valSet     *tmtypes.ValidatorSet
	ethSigner  ethtypes.Signer
	privKey    cryptotypes.PrivKey
	signer     keyring.Signer
	bondDenom  string
	stateDB    *statedb.StateDB

	//needed by test
	operatorAddr          sdk.AccAddress
	avsAddr               string
	assetId               string
	stakerId              string
	assetAddr             common.Address
	assetDecimal          uint32
	clientChainLzId       uint64
	depositAmount         sdkmath.Int
	delegationAmount      sdkmath.Int
	updatedAmountForOptIn sdkmath.Int
}

func TestOperatorTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "operator module Suite")
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest()
}
