package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// interface guards
var (
	_ sdk.Msg = &RegisterOperatorReq{}
	_ sdk.Msg = &OptInToCosmosChainRequest{}
	_ sdk.Msg = &InitOptOutFromCosmosChainRequest{}
)

// GetSigners returns the expected signers for the message.
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

// GetSigners returns the expected signers for the message.
func (m *OptInToCosmosChainRequest) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.Address)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *OptInToCosmosChainRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	return nil
}

// GetSigners returns the expected signers for the message.
func (m *InitOptOutFromCosmosChainRequest) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.Address)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *InitOptOutFromCosmosChainRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	return nil
}
