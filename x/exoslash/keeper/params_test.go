package keeper_test

import (
	"testing"

	testkeeper "github.com/exocore/testutil/keeper"
	"github.com/exocore/x/exoslash/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.ExoslashKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
