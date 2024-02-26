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

func createNRoundInfo(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.RoundInfo {
	items := make([]types.RoundInfo, n)
	for i := range items {
		items[i].TokenId = int32(i)

		keeper.SetRoundInfo(ctx, items[i])
	}
	return items
}

func TestRoundInfoGet(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNRoundInfo(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetRoundInfo(ctx,
			item.TokenId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestRoundInfoRemove(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNRoundInfo(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveRoundInfo(ctx,
			item.TokenId,
		)
		_, found := keeper.GetRoundInfo(ctx,
			item.TokenId,
		)
		require.False(t, found)
	}
}

func TestRoundInfoGetAll(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNRoundInfo(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllRoundInfo(ctx)),
	)
}
