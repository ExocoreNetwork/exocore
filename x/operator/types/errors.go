package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrNoKeyInTheStore = errorsmod.Register(
		ModuleName, 2,
		"there is no such key in the store",
	)

	ErrCliCmdInputArg = errorsmod.Register(
		ModuleName, 3,
		"there is an error in the input client command args",
	)

	ErrSlashInfo = errorsmod.Register(
		ModuleName, 4,
		"there is an error in the field of slash info",
	)

	ErrSlashInfoExist = errorsmod.Register(
		ModuleName, 5,
		"the slash info exists",
	)

	ErrParameterInvalid = errorsmod.Register(
		ModuleName, 6,
		"the input parameter is invalid",
	)

	ErrAlreadyOptedIn = errorsmod.Register(
		ModuleName, 7,
		"the operator has already opted in the avs",
	)

	ErrNotOptedIn = errorsmod.Register(
		ModuleName, 8,
		"the operator hasn't opted in the avs",
	)

	ErrSlashOccurredHeight = errorsmod.Register(
		ModuleName, 9,
		"the occurred height of slash event is error",
	)

	ErrConsKeyAlreadyInUse = errorsmod.Register(
		ModuleName, 10,
		"consensus key already in use by another operator",
	)

	ErrAlreadyRemovingKey = errorsmod.Register(
		ModuleName, 11, "operator already removing consensus key",
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

	ErrOperatorAlreadyExists = errorsmod.Register(
		ModuleName, 15,
		"operator already exists",
	)

	ErrUnknownChainID = errorsmod.Register(
		ModuleName, 16,
		"unknown chain id",
	)

	ErrOperatorNotRemovingKey = errorsmod.Register(
		ModuleName, 17,
		"operator not removing key",
	)

	ErrInvalidSlashPower = errorsmod.Register(
		ModuleName, 18,
		"the slash power is invalid",
	)

	ErrKeyAlreadyExist = errorsmod.Register(
		ModuleName, 19,
		"the key already exists",
	)

	ErrValueIsNilOrZero = errorsmod.Register(
		ModuleName, 20,
		"the value is nil or zero",
	)

	ErrKeyNotExistInMap = errorsmod.Register(
		ModuleName, 21,
		"the key doesn't exist in the map",
	)

	ErrNoSuchAvs = errorsmod.Register(
		ModuleName, 22,
		"no such avs",
	)

	ErrInvalidConsKey = errorsmod.Register(
		ModuleName, 23,
		"invalid consensus key",
	)

	ErrMinDelegationNotMet = errorsmod.Register(
		ModuleName, 24,
		"the operator USD value is less than the minimum self delegation",
	)
)
