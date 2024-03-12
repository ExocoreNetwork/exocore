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
		Validators: &types.Validators{
			Block: 42,
		},
		ValidatorUpdateBlock: &types.ValidatorUpdateBlock{},
		IndexRecentParams:    &types.IndexRecentParams{},
		IndexRecentMsg:       &types.IndexRecentMsg{},
		RecentMsgList: []types.RecentMsg{
			{
				Block: 0,
			},
			{
				Block: 1,
			},
		},
		RecentParamsList: []types.RecentParams{
			{
				Block: 0,
			},
			{
				Block: 1,
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
	require.Equal(t, genesisState.Validators, got.Validators)
	require.Equal(t, genesisState.ValidatorUpdateBlock, got.ValidatorUpdateBlock)
	require.Equal(t, genesisState.IndexRecentParams, got.IndexRecentParams)
	require.Equal(t, genesisState.IndexRecentMsg, got.IndexRecentMsg)
	require.ElementsMatch(t, genesisState.RecentMsgList, got.RecentMsgList)
	require.ElementsMatch(t, genesisState.RecentParamsList, got.RecentParamsList)
	// this line is used by starport scaffolding # genesis/test/assert
}
