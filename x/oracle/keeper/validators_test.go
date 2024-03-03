package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/ExocoreNetwork/exocore/testutil/keeper"
	"github.com/ExocoreNetwork/exocore/testutil/nullify"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNValidators(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Validators {
	items := make([]types.Validators, n)
	for i := range items {
		items[i].Block = uint64(i)

		keeper.SetValidators(ctx, items[i])
	}
	return items
}

func TestValidatorsGet(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNValidators(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetValidators(ctx,
			item.Block,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestValidatorsRemove(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNValidators(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveValidators(ctx,
			item.Block,
		)
		_, found := keeper.GetValidators(ctx,
			item.Block,
		)
		require.False(t, found)
	}
}

func TestValidatorsGetAll(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNValidators(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllValidators(ctx)),
	)
}
