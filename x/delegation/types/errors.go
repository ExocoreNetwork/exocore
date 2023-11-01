package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrNoKeyInTheStore  = errorsmod.Register(ModuleName, 0, "there is not the key for in the store")
	ErrOperatorIsFrozen = errorsmod.Register(ModuleName, 1, "the operator has been frozen")

	ErrOperatorNotExist = errorsmod.Register(ModuleName, 2, "the operator has not been registered")

	ErrOpAmountIsNegative = errorsmod.Register(ModuleName, 3, "the delegation or unDelegation amount is negative")

	OperatorAddrIsNotAccAddr = errorsmod.Register(ModuleName, 4, "the operator address isn't a valid acc addr")

	ErrSubAmountIsGreaterThanOriginal = errorsmod.Register(ModuleName, 5, "the sub amount is greater than the original amount")

	ErrParseDelegationKey = errorsmod.Register(ModuleName, 6, "delegation state key can't be parsed")
)
