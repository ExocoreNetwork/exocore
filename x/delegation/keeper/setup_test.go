package keeper_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"

	"github.com/stretchr/testify/suite"
)

type DelegationTestSuite struct {
	testutil.BaseTestSuite
}

var s *DelegationTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(DelegationTestSuite)
	suite.Run(t, s)

}

func (suite *DelegationTestSuite) SetupTest() {
	suite.DoSetupTest()
}
