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

func networkWithRoundDataObjects(t *testing.T, n int) (*network.Network, []types.RoundData) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	for i := 0; i < n; i++ {
		roundData := types.RoundData{
			TokenId: int32(i),
		}
		nullify.Fill(&roundData)
		state.RoundDataList = append(state.RoundDataList, roundData)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.RoundDataList
}

func TestShowRoundData(t *testing.T) {
	net, objs := networkWithRoundDataObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc      string
		idTokenId int32

		args []string
		err  error
		obj  types.RoundData
	}{
		{
			desc:      "found",
			idTokenId: objs[0].TokenId,

			args: common,
			obj:  objs[0],
		},
		{
			desc:      "not found",
			idTokenId: 100000,

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				strconv.Itoa(int(tc.idTokenId)),
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowRoundData(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryGetRoundDataResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.RoundData)
				require.Equal(t,
					nullify.Fill(&tc.obj),
					nullify.Fill(&resp.RoundData),
				)
			}
		})
	}
}

func TestListRoundData(t *testing.T) {
	net, objs := networkWithRoundDataObjects(t, 5)

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
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListRoundData(), args)
			require.NoError(t, err)
			var resp types.QueryAllRoundDataResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.RoundData), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.RoundData),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListRoundData(), args)
			require.NoError(t, err)
			var resp types.QueryAllRoundDataResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.RoundData), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.RoundData),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListRoundData(), args)
		require.NoError(t, err)
		var resp types.QueryAllRoundDataResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.RoundData),
		)
	})
}
