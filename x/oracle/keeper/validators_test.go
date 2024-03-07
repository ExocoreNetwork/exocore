package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/ExocoreNetwork/exocore/testutil/keeper"
	"github.com/ExocoreNetwork/exocore/testutil/nullify"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

func createTestValidators(keeper *keeper.Keeper, ctx sdk.Context) types.Validators {
	item := types.Validators{}
	keeper.SetValidators(ctx, item)
	return item
}

func TestValidatorsGet(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	item := createTestValidators(keeper, ctx)
	rst, found := keeper.GetValidators(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestValidatorsRemove(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	createTestValidators(keeper, ctx)
	keeper.RemoveValidators(ctx)
	_, found := keeper.GetValidators(ctx)
	require.False(t, found)
}
