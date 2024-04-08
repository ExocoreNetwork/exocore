package avs_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/precompiles/avs"
	"github.com/ExocoreNetwork/exocore/testutil"
	"github.com/stretchr/testify/suite"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var s *AVSManagerPrecompileSuite

type AVSManagerPrecompileSuite struct {
	testutil.BaseTestSuite
	precompile *avs.Precompile
}

func TestPrecompileTestSuite(t *testing.T) {
	s = new(AVSManagerPrecompileSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "AVSManager Precompile Suite")
}

func (s *AVSManagerPrecompileSuite) SetupTest() {
	s.DoSetupTest()
	precompile, err := avs.NewPrecompile(s.App.AVSManagerKeeper, s.App.AuthzKeeper)
	s.Require().NoError(err)
	s.precompile = precompile
}
