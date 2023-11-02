package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/exoslash module sentinel errors
var (
	ErrInvalidEvmAddressFormat  = errorsmod.Register(ModuleName, 0, "the evm address format is error")
	ErrInvalidLzUaTopicIdLength = errorsmod.Register(ModuleName, 1, "the LZUaTopicId length isn't equal to HashLength")
	ErrNoParamsKey              = errorsmod.Register(ModuleName, 2, "there is no stored key for slash module params")
	ErrSlashAmountIsNegative    = errorsmod.Register(ModuleName, 3, "the slash amount is negative")
	ErrSlashAssetNotExist       = errorsmod.Register(ModuleName, 4, "the slash asset doesn't exist")
)
