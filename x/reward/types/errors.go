package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/reward module sentinel errors
var (
	ErrSample                   = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrInvalidEvmAddressFormat  = errorsmod.Register(ModuleName, 0, "the evm address format is error")
	ErrInvalidLzUaTopicIdLength = errorsmod.Register(ModuleName, 1, "the LZUaTopicId length isn't equal to HashLength")
	ErrNoParamsKey              = errorsmod.Register(ModuleName, 2, "there is no stored key for params")
	ErrRewardAmountIsNegative  = errorsmod.Register(ModuleName, 3, "the reward amount is negative")
	ErrRewardAssetNotExist     = errorsmod.Register(ModuleName, 4, "the reward asset doesn't exist")
)
