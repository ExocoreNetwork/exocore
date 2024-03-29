package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrInvalidEvmAddressFormat  = errorsmod.Register(ModuleName, 0, "the evm address format is error")
	ErrInvalidLzUaTopicIDLength = errorsmod.Register(ModuleName, 1, "the LZUaTopicID length isn't equal to HashLength")
	ErrNoParamsKey              = errorsmod.Register(ModuleName, 2, "there is no stored key for deposit module params")
	ErrDepositAmountIsNegative  = errorsmod.Register(ModuleName, 3, "the deposit amount is negative")
	ErrDepositAssetNotExist     = errorsmod.Register(ModuleName, 4, "the deposit asset doesn't exist")
)
