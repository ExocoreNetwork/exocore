package types

import (
	errorsmod "cosmossdk.io/errors"
)

// errors
var (
	ErrNoClientChainKey      = errorsmod.Register(ModuleName, 0, "there is no stored key for the input chain index")
	ErrNoClientChainAssetKey = errorsmod.Register(ModuleName, 1, "there is no stored key for the input assetId")

	ErrNoStakerAssetKey = errorsmod.Register(ModuleName, 2, "there is no stored key for the input staker and assetId")

	ErrSubAmountIsMoreThanOrigin = errorsmod.Register(ModuleName, 3, "the amount that want to decrease is more than the original state amount")

	ErrNoOperatorAssetKey = errorsmod.Register(ModuleName, 4, "there is no stored key for the input operator address and assetId")

	ErrParseAssetsStateKey = errorsmod.Register(ModuleName, 5, "assets state key can't be parsed")

	ErrCliCmdInputArg = errorsmod.Register(ModuleName, 6, "there is an error in the input client command args")

	ErrInputPointerIsNil = errorsmod.Register(ModuleName, 7, "the input pointer is nil")

	ErrOperatorAddr = errorsmod.Register(ModuleName, 8, "the operator address isn't a valid acc addr")

	ErrNoKeyInTheStore = errorsmod.Register(ModuleName, 9, "there is not the key for in the store")
)
