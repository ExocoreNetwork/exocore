package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/avs module sentinel errors
var (
	ErrNoKeyInTheStore = errorsmod.Register(ModuleName, 0, "there is not the key for in the store")
	ErrSample          = sdkerrors.Register(ModuleName, 1100, "sample error")
)
