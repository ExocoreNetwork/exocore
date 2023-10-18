package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrInvalidEvmAddressFormat = errorsmod.Register(ModuleName, 0, "the evm address format is error")
	ErrNoParamsKey             = errorsmod.Register(ModuleName, 0, "there is no stored key for params")
)
