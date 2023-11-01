package exoslash_test

import (
	"testing"

	keepertest "github.com/exocore/testutil/keeper"
	"github.com/exocore/testutil/nullify"
	"github.com/exocore/x/exoslash"
	"github.com/exocore/x/exoslash/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ExoslashKeeper(t)
	exoslash.InitGenesis(ctx, *k, genesisState)
	got := exoslash.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
