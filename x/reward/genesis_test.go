package reward_test

import (
	"testing"

	keepertest "github.com/exocore/testutil/keeper"
	"github.com/exocore/testutil/nullify"
	"github.com/exocore/x/reward"
	"github.com/exocore/x/reward/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:	types.DefaultParams(),
		
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.RewardKeeper(t)
	reward.InitGenesis(ctx, *k, genesisState)
	got := reward.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	

	// this line is used by starport scaffolding # genesis/test/assert
}
