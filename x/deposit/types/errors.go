package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrNoParamsKey             = errorsmod.Register(ModuleName, 1, "there is no stored key for deposit module params")
	ErrDepositAmountIsNegative = errorsmod.Register(ModuleName, 2, "the deposit amount is negative")
	ErrDepositAssetNotExist    = errorsmod.Register(ModuleName, 3, "the deposit asset doesn't exist")
)
