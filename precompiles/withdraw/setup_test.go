package withdraw_test

import (
	"github.com/ExocoreNetwork/exocore/testutil"
	"testing"

	"github.com/ExocoreNetwork/exocore/precompiles/withdraw"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

var s *WithdrawPrecompileTestSuite

type WithdrawPrecompileTestSuite struct {
	testutil.BaseTestSuite
	precompile *withdraw.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(WithdrawPrecompileTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Withdraw Precompile Suite")
}

func (s *WithdrawPrecompileTestSuite) SetupTest() {
	s.DoSetupTest()
	precompile, err := withdraw.NewPrecompile(s.App.StakingAssetsManageKeeper, s.App.WithdrawKeeper, s.App.AuthzKeeper)
	s.Require().NoError(err)
	s.precompile = precompile
}
