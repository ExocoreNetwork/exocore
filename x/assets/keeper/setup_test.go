package keeper_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"

	"github.com/stretchr/testify/suite"
)

type StakingAssetsTestSuite struct {
	testutil.BaseTestSuite
}

var s *StakingAssetsTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(StakingAssetsTestSuite)
	suite.Run(t, s)
}

func (suite *StakingAssetsTestSuite) SetupTest() {
	suite.DoSetupTest()
}
