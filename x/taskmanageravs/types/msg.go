package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &RegisterAVSTaskReq{}
)

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *RegisterAVSTaskReq) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.AVSAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *RegisterAVSTaskReq) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.AVSAddress); err != nil {
		return errorsmod.Wrap(err, "invalid  address")
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m *RegisterAVSTaskReq) GetSignBytes() []byte {
	return nil
}
