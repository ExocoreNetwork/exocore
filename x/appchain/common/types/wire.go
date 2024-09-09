package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// PacketAckResult is the acknowledgment result of a packet.
type PacketAckResult []byte

// SlashPacketHandledResult is the success acknowledgment result of a slash packet.
var SlashPacketHandledResult = PacketAckResult([]byte{byte(2)})

func (vdt SlashPacketData) Validate() error {
	// vdt.Validator.Address must be a consensus address
	if err := sdk.VerifyAddressFormat(vdt.Validator.Address); err != nil {
		return ErrInvalidPacketData.Wrapf("invalid validator: %s", err.Error())
	}
	// vdt.Validator.Power must be positive
	if vdt.Validator.Power == 0 {
		return ErrInvalidPacketData.Wrap("validator power cannot be zero")
	}
	// ValsetUpdateId can be zero for the first validator set, so we don't validate it here.
	if vdt.Infraction != stakingtypes.Infraction_INFRACTION_DOWNTIME {
		// only downtime infractions are supported at this time
		return ErrInvalidPacketData.Wrapf("invalid infraction type: %s", vdt.Infraction.String())
	}

	return nil
}
