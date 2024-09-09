package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidChannelFlow = errorsmod.Register(
		ModuleName, 2,
		"invalid message sent to channel end",
	)
	ErrDuplicateChannel = errorsmod.Register(
		ModuleName, 3,
		"channel already exists",
	)
	ErrInvalidVersion = errorsmod.Register(
		ModuleName, 4,
		"invalid version",
	)
	ErrInvalidHandshakeMetadata = errorsmod.Register(
		ModuleName, 5,
		"invalid handshake metadata",
	)
	ErrChannelNotFound = errorsmod.Register(
		ModuleName, 6,
		"channel not found",
	)
	ErrClientNotFound = errorsmod.Register(
		ModuleName, 7,
		"client not found",
	)
	ErrInvalidPacketData = errorsmod.Register(
		ModuleName, 8,
		"invalid packet data (but successfully unmarshalled)",
	)
)
