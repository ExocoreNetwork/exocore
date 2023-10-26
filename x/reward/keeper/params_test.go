package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "github.com/exocore/testutil/keeper"
	"github.com/exocore/x/reward/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.RewardKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
