package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var (
	ErrInvalidGenesisData = errorsmod.Register(
		ModuleName, 0,
		"the genesis data supplied is invalid",
	)
)
