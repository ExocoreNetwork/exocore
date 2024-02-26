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

func createNRoundData(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.RoundData {
	items := make([]types.RoundData, n)
	for i := range items {
		items[i].TokenId = int32(i)

		keeper.SetRoundData(ctx, items[i])
	}
	return items
}

func TestRoundDataGet(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNRoundData(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetRoundData(ctx,
			item.TokenId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestRoundDataRemove(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNRoundData(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveRoundData(ctx,
			item.TokenId,
		)
		_, found := keeper.GetRoundData(ctx,
			item.TokenId,
		)
		require.False(t, found)
	}
}

func TestRoundDataGetAll(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNRoundData(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllRoundData(ctx)),
	)
}
