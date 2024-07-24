package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

var ErrInvalidGenesisData = errorsmod.Register(
	ModuleName, 2,
	"the genesis data supplied is invalid",
)
