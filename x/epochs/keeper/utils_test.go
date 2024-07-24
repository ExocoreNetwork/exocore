package keeper_test

import (
	"github.com/ExocoreNetwork/exocore/x/epochs/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
)

// Test helpers
func (suite *KeeperTestSuite) PostSetup() {
	// setup query helpers
	queryHelper := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, suite.App.EpochsKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)
}
