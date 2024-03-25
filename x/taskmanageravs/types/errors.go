package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/taskmanageravs module sentinel errors
var (
	ErrSample           = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrNoKeyInTheStore  = errorsmod.Register(ModuleName, 0, "there is not the key for in the store")
	ErrNotYetRegistered = errorsmod.Register(ModuleName, 1101, "this AVS has not been registered yet")
)
