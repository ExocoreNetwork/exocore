package keeper_test

import (
	"fmt"
	"time"

	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"
)

func (suite *AVSTestSuite) TestEpochIdentifierAfterEpochEnd() {
	testCases := []struct {
		name            string
		epochIdentifier string
		expDistribution bool
	}{
		{
			"correct epoch identifier",
			epochstypes.DayEpochID,
			true,
		},
		{
			"incorrect epoch identifier",
			epochstypes.WeekEpochID,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest()

			futureCtx := suite.ctx.WithBlockTime(time.Now().Add(time.Hour))
			newHeight := suite.app.LastBlockHeight() + 1

			suite.app.EpochsKeeper.BeforeEpochStart(futureCtx, tc.epochIdentifier, newHeight)
			suite.app.EpochsKeeper.AfterEpochEnd(futureCtx, tc.epochIdentifier, newHeight)

			suite.app.EpochsKeeper.AfterEpochEnd(futureCtx, tc.epochIdentifier, newHeight)

		})
	}
}
