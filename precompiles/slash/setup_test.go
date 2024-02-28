package slash_test

import (
	"github.com/ExocoreNetwork/exocore/testutil"
	"testing"

	"github.com/ExocoreNetwork/exocore/precompiles/slash"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

var s *SlashPrecompileTestSuite

type SlashPrecompileTestSuite struct {
	testutil.BaseTestSuite
	precompile *slash.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(SlashPrecompileTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Slash Precompile Suite")
}

func (s *SlashPrecompileTestSuite) SetupTest() {
	s.DoSetupTest()
	precompile, err := slash.NewPrecompile(s.App.StakingAssetsManageKeeper, s.App.ExoSlashKeeper, s.App.AuthzKeeper)
	s.Require().NoError(err)
	s.precompile = precompile
}
