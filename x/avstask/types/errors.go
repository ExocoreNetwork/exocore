package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/avs_task module sentinel errors
var (
	ErrSample           = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrNoKeyInTheStore  = errorsmod.Register(ModuleName, 0, "there is no such key in the store")
	ErrNotYetRegistered = errorsmod.Register(ModuleName, 1101, "this AVS has not been registered yet")
)
