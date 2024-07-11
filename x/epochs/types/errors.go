package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

const (
	codeErrInvalidGenesisData = iota + 2
	codeErrDuplicateEpochInfo
)

var (
	ErrInvalidGenesisData = errorsmod.Register(
		ModuleName, codeErrInvalidGenesisData, "the genesis data supplied is invalid",
	)
	ErrDuplicateEpochInfo = errorsmod.Register(
		ModuleName, codeErrDuplicateEpochInfo, "epoch info already exists in the store",
	)
)
