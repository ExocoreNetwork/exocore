package deposit_test

import (
	"github.com/ExocoreNetwork/exocore/testutil"
	"testing"

	"github.com/ExocoreNetwork/exocore/precompiles/deposit"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

var s *DepositPrecompileSuite

type DepositPrecompileSuite struct {
	testutil.BaseTestSuite

	precompile *deposit.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(DepositPrecompileSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Distribution Precompile Suite")
}

func (s *DepositPrecompileSuite) SetupTest() {
	s.DoSetupTest()
	precompile, err := deposit.NewPrecompile(s.App.StakingAssetsManageKeeper, s.App.DepositKeeper, s.App.AuthzKeeper)
	s.Require().NoError(err)
	s.precompile = precompile
}
