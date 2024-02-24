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

func createNPrices(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Prices {
	items := make([]types.Prices, n)
	for i := range items {
		items[i].TokenId = int32(i)

		keeper.SetPrices(ctx, items[i])
	}
	return items
}

func TestPricesGet(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNPrices(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetPrices(ctx,
			item.TokenId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestPricesRemove(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNPrices(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePrices(ctx,
			item.TokenId,
		)
		_, found := keeper.GetPrices(ctx,
			item.TokenId,
		)
		require.False(t, found)
	}
}

func TestPricesGetAll(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNPrices(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPrices(ctx)),
	)
}
