package types

import (
	"testing"

	"github.com/ExocoreNetwork/exocore/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	// sdkerrors "cosmossdk.io/errors""
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePrice_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePrice
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePrice{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreatePrice{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
