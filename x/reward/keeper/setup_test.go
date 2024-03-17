package keeper_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"

	//nolint:revive // dot imports are fine for Ginkgo
	. "github.com/onsi/ginkgo/v2"
	//nolint:revive // dot imports are fine for Ginkgo
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

type RewardTestSuite struct {
	testutil.BaseTestSuite
}

var s *RewardTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(RewardTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

// SetupTest setup test environment, it uses`require.TestingT` to support both `testing.T` and `testing.B`.
func (suite *RewardTestSuite) SetupTest() {
	suite.DoSetupTest()
}
