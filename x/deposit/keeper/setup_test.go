package keeper_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"

	"github.com/stretchr/testify/suite"
)

type DepositTestSuite struct {
	testutil.BaseTestSuite
}

var s *DepositTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(DepositTestSuite)
	suite.Run(t, s)
}

func (suite *DepositTestSuite) SetupTest() {
	suite.DoSetupTest()
}
