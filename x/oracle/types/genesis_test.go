package types_test

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				PricesList: []types.Prices{
					{
						TokenId: 0,
					},
					{
						TokenId: 1,
					},
				},
				Validators: &types.Validators{
					Block: 45,
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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated prices",
			genState: &types.GenesisState{
				PricesList: []types.Prices{
					{
						TokenId: 0,
					},
					{
						TokenId: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated recentMsg",
			genState: &types.GenesisState{
				RecentMsgList: []types.RecentMsg{
					{
						Block: 0,
					},
					{
						Block: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated recentParams",
			genState: &types.GenesisState{
				RecentParamsList: []types.RecentParams{
					{
						Block: 0,
					},
					{
						Block: 0,
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
