package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/oracle module sentinel errors
var (
	ErrSample                     = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidMsg                 = sdkerrors.Register(ModuleName, 1, "invalid input create price")
	ErrPriceProposalIgnored       = sdkerrors.Register(ModuleName, 2, "price proposal ignored")
	ErrPriceProposalFormatInvalid = sdkerrors.Register(ModuleName, 3, "price proposal message format invalid")
)
