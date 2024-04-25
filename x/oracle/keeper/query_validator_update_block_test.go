package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/ExocoreNetwork/exocore/testutil/keeper"
	"github.com/ExocoreNetwork/exocore/testutil/nullify"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

func TestValidatorUpdateBlockQuery(t *testing.T) {
	keeper, ctx := keepertest.OracleKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestValidatorUpdateBlock(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QueryGetValidatorUpdateBlockRequest
		response *types.QueryGetValidatorUpdateBlockResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetValidatorUpdateBlockRequest{},
			response: &types.QueryGetValidatorUpdateBlockResponse{ValidatorUpdateBlock: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ValidatorUpdateBlock(wctx, tc.request)
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
