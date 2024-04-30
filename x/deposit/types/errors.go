package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrNoParamsKey          = errorsmod.Register(ModuleName, 1, "there is no stored key for deposit module params")
	ErrInvalidDepositAmount = errorsmod.Register(ModuleName, 2, "the deposit amount is invalid")
)
