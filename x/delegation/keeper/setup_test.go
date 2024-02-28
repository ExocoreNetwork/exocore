package keeper_test

import (
	"github.com/ExocoreNetwork/exocore/testutil"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

type DelegationTestSuite struct {
	testutil.BaseTestSuite
}

var s *DelegationTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(DelegationTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *DelegationTestSuite) SetupTest() {
	suite.DoSetupTest()
}
