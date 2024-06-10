package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgDelegation{}
	_ sdk.Msg = &MsgUndelegation{}
)

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgDelegation) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.BaseInfo.FromAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *MsgDelegation) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.BaseInfo.FromAddress); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m *MsgDelegation) GetSignBytes() []byte {
	return nil
}

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *MsgUndelegation) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.BaseInfo.FromAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *MsgUndelegation) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.BaseInfo.FromAddress); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m *MsgUndelegation) GetSignBytes() []byte {
	return nil
}
