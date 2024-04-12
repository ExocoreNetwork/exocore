package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ExocoreNetwork/exocore/testutil/network"
	"github.com/ExocoreNetwork/exocore/testutil/nullify"
	"github.com/ExocoreNetwork/exocore/x/oracle/client/cli"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithRecentMsgObjects(t *testing.T, n int) (*network.Network, []types.RecentMsg) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		recentMsg := types.RecentMsg{
			Block: uint64(i),
		}
		nullify.Fill(&recentMsg)
		state.RecentMsgList = append(state.RecentMsgList, recentMsg)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	network, err := network.New(t, t.TempDir(), cfg)
	require.NoError(t, err)
	return network, state.RecentMsgList
}

func TestShowRecentMsg(t *testing.T) {
	net, objs := networkWithRecentMsgObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc    string
		idBlock uint64

		args []string
		err  error
		obj  types.RecentMsg
	}{
		{
			desc:    "found",
			idBlock: objs[0].Block,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			idBlock: 100000,

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				strconv.Itoa(int(tc.idBlock)),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowRecentMsg(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetRecentMsgResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.RecentMsg)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.RecentMsg),
				)
			}
		})
	}
}

func TestListRecentMsg(t *testing.T) {
	net, objs := networkWithRecentMsgObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListRecentMsg(), args)
			require.NoError(t, err)
			var resp types.QueryAllRecentMsgResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.RecentMsg), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.RecentMsg),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListRecentMsg(), args)
			require.NoError(t, err)
			var resp types.QueryAllRecentMsgResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.RecentMsg), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.RecentMsg),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListRecentMsg(), args)
		require.NoError(t, err)
		var resp types.QueryAllRecentMsgResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.RecentMsg),
		)
	})
}
