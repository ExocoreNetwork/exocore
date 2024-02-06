package keeper_test

import (
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
	"github.com/stretchr/testify/suite"
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
	operatorAddr sdk.AccAddress
	avsAddr      string
}

func TestOperatorTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "operator module Suite")
}

func (s *KeeperTestSuite) SetupTest() {
	s.DoSetupTest()
}
