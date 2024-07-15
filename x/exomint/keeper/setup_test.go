package keeper_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil"
	"github.com/stretchr/testify/suite"

	"github.com/ExocoreNetwork/exocore/x/exomint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

var s *KeeperTestSuite

type KeeperTestSuite struct {
	testutil.BaseTestSuite
	queryClient types.QueryClient
}

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	suite.Run(t, s)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.DoSetupTest()
	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.App.ExomintKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}
