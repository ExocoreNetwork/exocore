package types

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/ExocoreNetwork/exocore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
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
	return m.Info.ValidateBasic()
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

// GetSigners returns the expected signers for the message.
func (m *OptIntoAVSReq) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.FromAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *OptIntoAVSReq) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.FromAddress); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	if m.AvsAddress == "" {
		return errorsmod.Wrap(ErrParameterInvalid, "AVS address is empty")
	}
	// TODO: there is no specific restriction applied on ChainID in Cosmos, so this check
	// could potentially be removed.
	if !common.IsHexAddress(m.AvsAddress) && !types.IsValidChainID(m.AvsAddress) {
		return errorsmod.Wrap(
			ErrParameterInvalid, "AVS address is not a valid hex address or chain id",
		)
	}
	return nil
}

// GetSigners returns the expected signers for the message.
func (m *OptOutOfAVSReq) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.FromAddress)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *OptOutOfAVSReq) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.FromAddress); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	if m.AvsAddress == "" {
		return errorsmod.Wrap(ErrParameterInvalid, "AVS address is empty")
	}
	if !common.IsHexAddress(m.AvsAddress) && !types.IsValidChainID(m.AvsAddress) {
		return errorsmod.Wrap(
			ErrParameterInvalid, "AVS address is not a valid hex address or chain id",
		)
	}
	return nil
}
