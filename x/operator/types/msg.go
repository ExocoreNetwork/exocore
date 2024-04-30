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
	_ sdk.Msg = &OptIntoAVSReq{}
	_ sdk.Msg = &OptOutOfAVSReq{}
	_ sdk.Msg = &SetConsKeyReq{}
	_ sdk.Msg = &InitConsKeyRemovalReq{}
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

// Route returns the transaction route. This must be specified for successful signing.
func (m *RegisterOperatorReq) Route() string {
	return RouterKey
}

// GetSigners returns the expected signers for the message.
func (m *SetConsKeyReq) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.Address)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *SetConsKeyReq) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	if !types.IsValidChainID(m.ChainID) {
		return errorsmod.Wrap(ErrParameterInvalid, "invalid chain id")
	}
	if keyType, keyString, err := ParseConsensusKeyFromJSON(m.PublicKey); err != nil {
		return errorsmod.Wrap(err, "invalid public key")
	} else if keyType != "/cosmos.crypto.ed25519.PubKey" {
		return errorsmod.Wrap(ErrParameterInvalid, "invalid public key type")
	} else if _, err := StringToPubKey(keyString); err != nil {
		return errorsmod.Wrap(err, "invalid public key")
	}
	return nil
}

// Route returns the transaction route. This must be specified for successful signing.
func (m *SetConsKeyReq) Route() string {
	return RouterKey
}

// GetSigners returns the expected signers for the message.
func (m *InitConsKeyRemovalReq) GetSigners() []sdk.AccAddress {
	addr := sdk.MustAccAddressFromBech32(m.Address)
	return []sdk.AccAddress{addr}
}

// ValidateBasic does a sanity check of the provided data
func (m *InitConsKeyRemovalReq) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Address); err != nil {
		return errorsmod.Wrap(err, "invalid from address")
	}
	if !types.IsValidChainID(m.ChainID) {
		return errorsmod.Wrap(ErrParameterInvalid, "invalid chain id")
	}
	return nil
}

// Route returns the transaction route. This must be specified for successful signing.
func (m *InitConsKeyRemovalReq) Route() string {
	return RouterKey
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

// Route returns the transaction route. This must be specified for successful signing.
func (m *OptIntoAVSReq) Route() string {
	return RouterKey
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

// Route returns the transaction route. This must be specified for successful signing.
func (m *OptOutOfAVSReq) Route() string {
	return RouterKey
}
