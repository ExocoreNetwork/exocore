package keeper_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

type StakingAssetsTestSuite struct {
	testutil.BaseTestSuite
}

var s *StakingAssetsTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(StakingAssetsTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *StakingAssetsTestSuite) SetupTest() {
	suite.DoSetupTest()
}
