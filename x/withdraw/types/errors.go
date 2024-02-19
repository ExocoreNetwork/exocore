package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/withdraw module sentinel errors
var (
	ErrSample                   = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrInvalidEvmAddressFormat  = errorsmod.Register(ModuleName, 0, "the evm address format is error")
	ErrInvalidLzUaTopicIDLength = errorsmod.Register(ModuleName, 1, "the LZUaTopicID length isn't equal to HashLength")
	ErrNoParamsKey              = errorsmod.Register(ModuleName, 2, "there is no stored key for params")
	ErrWithdrawAmountIsNegative = errorsmod.Register(ModuleName, 3, "the withdraw amount is negative")
	ErrWithdrawAssetNotExist    = errorsmod.Register(ModuleName, 4, "the withdraw asset doesn't exist")
)
