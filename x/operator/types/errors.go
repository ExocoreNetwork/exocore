package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrNoKeyInTheStore = errorsmod.Register(ModuleName, 0, "there is not the key for in the store")

	ErrCliCmdInputArg = errorsmod.Register(ModuleName, 1, "there is an error in the input client command args")

	ErrSlashInfo = errorsmod.Register(ModuleName, 2, "there is an error in the field of slash info")

	ErrSlashInfoExist = errorsmod.Register(ModuleName, 3, "the slash info exists")

	ErrParameterInvalid = errorsmod.Register(ModuleName, 4, "the input parameter is invalid")

	ErrAlreadyOptedIn = errorsmod.Register(ModuleName, 5, "the operator has already opted in the avs")

	ErrNotOptedIn = errorsmod.Register(ModuleName, 6, "the operator hasn't opted in the avs")
)
