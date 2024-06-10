package keeper_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"

	"github.com/stretchr/testify/suite"
)

type WithdrawTestSuite struct {
	testutil.BaseTestSuite
}

var s *WithdrawTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(WithdrawTestSuite)
	suite.Run(t, s)
}

// SetupTest setup test environment, it uses`require.TestingT` to support both `testing.T` and `testing.B`.
func (suite *WithdrawTestSuite) SetupTest() {
	suite.DoSetupTest()
}
