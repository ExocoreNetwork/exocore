package keeper_test

import (
	"github.com/ExocoreNetwork/exocore/testutil"
	"testing"

	//nolint:revive // dot imports are fine for Ginkgo
	. "github.com/onsi/ginkgo/v2"
	//nolint:revive // dot imports are fine for Ginkgo
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

type SlashTestSuite struct {
	testutil.BaseTestSuite
}

var s *SlashTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(SlashTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

// SetupTest setup test environment, it uses`require.TestingT` to support both `testing.T` and `testing.B`.
func (suite *SlashTestSuite) SetupTest() {
	suite.DoSetupTest()
}
