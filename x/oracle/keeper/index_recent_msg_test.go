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

func createTestIndexRecentMsg(keeper *keeper.Keeper, ctx sdk.Context) types.IndexRecentMsg {
	item := types.IndexRecentMsg{}
	keeper.SetIndexRecentMsg(ctx, item)
	return item
}

func TestIndexRecentMsgGet(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	item := createTestIndexRecentMsg(keeper, ctx)
	rst, found := keeper.GetIndexRecentMsg(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}

func TestIndexRecentMsgRemove(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	createTestIndexRecentMsg(keeper, ctx)
	keeper.RemoveIndexRecentMsg(ctx)
	_, found := keeper.GetIndexRecentMsg(ctx)
	require.False(t, found)
}
