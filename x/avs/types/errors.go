package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/avs module sentinel errors
var (
	ErrNoKeyInTheStore = errorsmod.Register(ModuleName, 0, "there is no such key in the store")

	ErrAlreadyRegistered = errorsmod.Register(
		ModuleName, 1,
		"Error: Already registered",
	)
	ErrUnregisterNonExistent = errorsmod.Register(
		ModuleName, 2,
		"Error: No available avs to DeRegisterAction",
	)

	ErrInvalidAction = errorsmod.Register(
		ModuleName, 3,
		"Error: Undefined action",
	)
)
