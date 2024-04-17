package keeper_test

import (
	"strconv"
	"testing"

	keepertest "github.com/ExocoreNetwork/exocore/testutil/keeper"
	"github.com/ExocoreNetwork/exocore/testutil/nullify"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/testdata"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNPrices(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Prices {
	items := make([]types.Prices, n)
	for i := range items {
		items[i].TokenID = uint64(i + 1)
		items[i] = types.Prices{
			TokenID:     uint64(i + 1),
			NextRoundID: 2,
			PriceList: []*types.PriceTimeRound{
				testdata.PTR1,
				testdata.PTR2,
				testdata.PTR3,
				testdata.PTR4,
				testdata.PTR5,
			},
		}
		keeper.SetPrices(ctx, items[i])
	}
	return items
}

func TestPricesGet(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	keeper.SetPrices(ctx, testdata.P1)
	rst, found := keeper.GetPrices(ctx, 1)
	require.True(t, found)
	pRes := testdata.P1
	pRes.PriceList = append([]*types.PriceTimeRound{{}}, testdata.P1.PriceList...)
	require.Equal(t, pRes, rst)
	// items := createNPrices(keeper, ctx, 10)
	//
	//	for _, item := range items {
	//		rst, found := keeper.GetPrices(ctx,
	//			item.TokenId,
	//		)
	//		require.True(t, found)
	//		require.Equal(t,
	//			nullify.Fill(&item),
	//			nullify.Fill(&rst),
	//		)
	//	}
}

func TestPricesRemove(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	items := createNPrices(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePrices(ctx,
			item.TokenID,
		)
		_, found := keeper.GetPrices(ctx,
			item.TokenID,
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
