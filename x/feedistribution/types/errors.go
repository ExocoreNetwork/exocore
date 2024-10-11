package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/feedistribution module sentinel errors
var (
	ErrEpochNotFound = sdkerrors.Register(
		ModuleName, 1102,
		"Error: epoch info not found",
	)
)
