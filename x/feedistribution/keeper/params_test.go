package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/ExocoreNetwork/exocore/testutil/keeper"
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.FeedistributeKeeper(t)
	params := types.DefaultParams()

	// k.SetParams(ctx, params)
	require.EqualValues(t, params, k.GetParams(ctx))
}
