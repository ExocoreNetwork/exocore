package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"

	"github.com/ExocoreNetwork/exocore/testutil/network"
	"github.com/ExocoreNetwork/exocore/testutil/nullify"
	"github.com/ExocoreNetwork/exocore/x/oracle/client/cli"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

func networkWithIndexRecentParamsObjects(t *testing.T) (*network.Network, types.IndexRecentParams) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	indexRecentParams := &types.IndexRecentParams{}
	nullify.Fill(&indexRecentParams)
	state.IndexRecentParams = indexRecentParams
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	network, err := network.New(t, t.TempDir(), cfg)
	require.NoError(t, err)
	return network, *state.IndexRecentParams
}

func TestShowIndexRecentParams(t *testing.T) {
	net, obj := networkWithIndexRecentParamsObjects(t)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc string
		args []string
		err  error
		obj  types.IndexRecentParams
	}{
		{
			desc: "get",
			args: common,
			obj:  obj,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowIndexRecentParams(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetIndexRecentParamsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.IndexRecentParams)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.IndexRecentParams),
				)
			}
		})
	}
}
