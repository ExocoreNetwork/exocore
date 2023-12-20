package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgSetExoCoreAddr{}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgSetExoCoreAddr) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.FromAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *MsgSetExoCoreAddr) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.FromAddress); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	if _, err := sdk.AccAddressFromBech32(m.SetAddress); err != nil {
		return errorsmod.Wrap(err, "invalid set address")
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m *MsgSetExoCoreAddr) GetSignBytes() []byte {
	return nil
}
