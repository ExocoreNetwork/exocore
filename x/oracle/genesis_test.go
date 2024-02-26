package oracle_test

import (
	"testing"

	keepertest "github.com/ExocoreNetwork/exocore/testutil/keeper"
	"github.com/ExocoreNetwork/exocore/testutil/nullify"
	"github.com/ExocoreNetwork/exocore/x/oracle"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PricesList: []types.Prices{
			{
				TokenId: 0,
			},
			{
				TokenId: 1,
			},
		},
		RoundInfoList: []types.RoundInfo{
			{
				TokenId: 0,
			},
			{
				TokenId: 1,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.OracleKeeper(t)
	oracle.InitGenesis(ctx, *k, genesisState)
	got := oracle.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.PricesList, got.PricesList)
	require.ElementsMatch(t, genesisState.RoundInfoList, got.RoundInfoList)
	// this line is used by starport scaffolding # genesis/test/assert
}
