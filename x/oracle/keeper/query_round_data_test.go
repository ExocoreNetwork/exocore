package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/ExocoreNetwork/exocore/testutil/keeper"
	"github.com/ExocoreNetwork/exocore/testutil/nullify"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestRoundDataQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNRoundData(keeper, ctx, 2)
	tests := []struct {
		desc     string
		request  *types.QueryGetRoundDataRequest
		response *types.QueryGetRoundDataResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRoundDataRequest{
				TokenId: msgs[0].TokenId,
			},
			response: &types.QueryGetRoundDataResponse{RoundData: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetRoundDataRequest{
				TokenId: msgs[1].TokenId,
			},
			response: &types.QueryGetRoundDataResponse{RoundData: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRoundDataRequest{
				TokenId: 100000,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.RoundData(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestRoundDataQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNRoundData(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllRoundDataRequest {
		return &types.QueryAllRoundDataRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.RoundDataAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RoundData), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.RoundData),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.RoundDataAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RoundData), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.RoundData),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.RoundDataAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.RoundData),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.RoundDataAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
