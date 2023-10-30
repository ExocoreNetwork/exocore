package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrNoOperatorInfoKey = errorsmod.Register(ModuleName, 0, "there is no stored key for the input operator address")
)
