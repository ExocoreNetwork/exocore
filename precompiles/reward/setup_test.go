package reward_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/precompiles/reward"

	"github.com/evmos/evmos/v14/x/evm/statedb"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	evmosapp "github.com/ExocoreNetwork/exocore/app"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
	"github.com/stretchr/testify/suite"
)

var s *PrecompileTestSuite

type PrecompileTestSuite struct {
	suite.Suite

	ctx        sdk.Context
	app        *evmosapp.ExocoreApp
	address    common.Address
	validators []stakingtypes.Validator
	valSet     *tmtypes.ValidatorSet
	ethSigner  ethtypes.Signer
	privKey    cryptotypes.PrivKey
	signer     keyring.Signer
	bondDenom  string

	precompile *reward.Precompile
	stateDB    *statedb.StateDB

	queryClientEVM evmtypes.QueryClient
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(PrecompileTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reward Precompile Suite")
}

func (s *PrecompileTestSuite) SetupTest() {
	s.DoSetupTest()
}
