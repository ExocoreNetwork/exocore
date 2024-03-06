package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &RegisterOperatorReq{}

	// add for dogfood
	_ sdk.Msg = &OptInToChainIDRequest{}
	_ sdk.Msg = &InitiateOptOutFromChainIDRequest{}
)

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (m *RegisterOperatorReq) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.FromAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *RegisterOperatorReq) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.FromAddress); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	return nil
}

// GetSignBytes implements the LegacyMsg interface.
func (m *RegisterOperatorReq) GetSignBytes() []byte {
	return nil
}

func (m *OptInToChainIDRequest) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.Address)
	return []sdk.AccAddress{addr}
}

func (m *OptInToChainIDRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	return nil
}

func (m *InitiateOptOutFromChainIDRequest) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.Address)
	return []sdk.AccAddress{addr}
}

func (m *InitiateOptOutFromChainIDRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	return nil
}
