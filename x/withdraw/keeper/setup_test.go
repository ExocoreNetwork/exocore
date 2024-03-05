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

type WithdrawTestSuite struct {
	testutil.BaseTestSuite
}

var s *WithdrawTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(WithdrawTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

// SetupTest setup test environment, it uses`require.TestingT` to support both `testing.T` and `testing.B`.
func (suite *WithdrawTestSuite) SetupTest() {
	suite.DoSetupTest()
}
