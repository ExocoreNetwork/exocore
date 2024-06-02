package keeper_test

import (
	"testing"

	testkeeper "github.com/ExocoreNetwork/exocore/testutil/keeper"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.OracleKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}

func TestUpdateTokenFeeder(t *testing.T) {
	tests := []struct {
		name        string
		tokenFeeder types.TokenFeeder
		height      uint64
		err         error
	}{
		// fail when add/update fields, before validation
		{
			name: "invalid update, empty fields to update",
			tokenFeeder: types.TokenFeeder{
				TokenID: 1,
			},
			height: 1,
			err:    types.ErrInvalidParams.Wrap("invalid tokenFeeder to update, no valid field set"),
		},
		{
			name: "invalid udpate, for not-start feeder, set StartbaseBlock to history height",
			tokenFeeder: types.TokenFeeder{
				TokenID: 1,
				// set current height to 100 to test fail case
				StartBaseBlock: 10,
			},
			height: 100,
			err:    types.ErrInvalidParams.Wrap("invalid tokenFeeder to update, invalid StartBaseBlock"),
		},
		{
			name: "invalid update, for running feeder, set EndBlock to history height",
			tokenFeeder: types.TokenFeeder{
				TokenID: 1,
				// set current height to 2000000 to test fail case
				EndBlock: 1500000,
			},
			height: 2000000,
			err:    types.ErrInvalidParams.Wrap("invalid tokenFeeder to update, invalid EndBlock"),
		},
		{
			name: "invalid update, for stopped feeder, restart a feeder with wrong StartRoundID",
			tokenFeeder: types.TokenFeeder{
				TokenID: 2,
				RuleID:  1,
				// set current height to 100
				StartBaseBlock: 1000,
				// should be 4
				StartRoundID: 5,
				Interval:     10,
			},
			height: 100,
			err:    types.ErrInvalidParams.Wrap("invalid tokenFeeder to update"),
		},
		// success adding/updating, but fail validation
		{
			name: "invalid update, for new feeder, EndBlock is not set properly",
			tokenFeeder: types.TokenFeeder{
				TokenID:        3,
				StartBaseBlock: 10,
				StartRoundID:   1,
				Interval:       10,
				EndBlock:       51,
			},
			height: 1,
			err:    types.ErrInvalidParams.Wrap("invalid tokenFeeder, invalid EndBlock"),
		},
		{
			name: "invalid update, for new feeder, tokenID not exists",
			tokenFeeder: types.TokenFeeder{
				TokenID:        4,
				StartBaseBlock: 10,
				StartRoundID:   1,
				Interval:       10,
				EndBlock:       58,
			},
			height: 1,
			err:    types.ErrInvalidParams.Wrap("invalid tokenFeeder, non-exist tokenID referred"),
		},
	}
	p := types.DefaultParams()
	p.Tokens = append(p.Tokens, &types.Token{
		Name:            "TEST",
		ChainID:         1,
		ContractAddress: "0x",
		Decimal:         8,
		Active:          true,
		AssetID:         "",
	})
	p.Tokens = append(p.Tokens, &types.Token{
		Name:            "TEST_NEW",
		ChainID:         1,
		ContractAddress: "0x",
		Decimal:         8,
		Active:          true,
		AssetID:         "",
	})
	p.TokenFeeders = append(p.TokenFeeders, &types.TokenFeeder{
		TokenID:        2,
		RuleID:         1,
		StartRoundID:   1,
		StartBaseBlock: 10,
		Interval:       10,
		EndBlock:       38,
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := p.UpdateTokenFeeder(&tt.tokenFeeder, tt.height)
			if err == nil {
				err = p.Validate()
			}
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			}
		})
	}
}
