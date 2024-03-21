package keeper_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/stretchr/testify/suite"
)

type DepositTestSuite struct {
	testutil.BaseTestSuite
}

var s *DepositTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(DepositTestSuite)
	suite.Run(t, s)

	// Run Ginkgo integration tests
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func (suite *DepositTestSuite) SetupTest() {
	suite.DoSetupTest()
}
