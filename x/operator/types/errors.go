package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrNoKeyInTheStore = errorsmod.Register(
		ModuleName, 0,
		"there is no such key in the store",
	)

	ErrCliCmdInputArg = errorsmod.Register(
		ModuleName, 1,
		"there is an error in the input client command args",
	)

	ErrSlashInfo = errorsmod.Register(
		ModuleName, 2,
		"there is an error in the field of slash info",
	)

	ErrSlashInfoExist = errorsmod.Register(
		ModuleName, 3,
		"the slash info exists",
	)

	ErrParameterInvalid = errorsmod.Register(
		ModuleName, 4,
		"the input parameter is invalid",
	)

	ErrAlreadyOptedIn = errorsmod.Register(
		ModuleName, 5,
		"the operator has already opted in the avs",
	)

	ErrNotOptedIn = errorsmod.Register(
		ModuleName, 6,
		"the operator hasn't opted in the avs",
	)

	ErrTheValueIsNegative = errorsmod.Register(
		ModuleName, 7,
		"the value is negative",
	)

	ErrSlashContractNotMatch = errorsmod.Register(
		ModuleName, 8,
		"the slash contract isn't the slash contract address saved in the opted-in info",
	)

	ErrSlashOccurredHeight = errorsmod.Register(
		ModuleName, 9,
		"the occurred height of slash event is error",
	)

	ErrConsKeyAlreadyInUse = errorsmod.Register(
		ModuleName, 10,
		"consensus key already in use by another operator",
	)

	ErrAlreadyOptingOut = errorsmod.Register(
		ModuleName, 11, "operator already opting out",
	)

	ErrInvalidPubKey = errorsmod.Register(
		ModuleName, 12,
		"invalid public key",
	)

	ErrInvalidGenesisData = errorsmod.Register(
		ModuleName, 13,
		"the genesis data supplied is invalid",
	)

	ErrInvalidAvsAddr = errorsmod.Register(
		ModuleName, 14,
		"avs address should be a hex evm contract address",
	)

	ErrInvalidSlashPower = errorsmod.Register(
		ModuleName, 15,
		"the slash power is invalid",
	)

	ErrKeyAlreadyExist = errorsmod.Register(
		ModuleName, 16,
		"the key already exists",
	)

	ErrValueIsNilOrZero = errorsmod.Register(
		ModuleName, 17,
		"the value is nil or zero",
	)
)
