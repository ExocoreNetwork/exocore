package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

const (
	invalidMsg = iota + 2
	priceProposalIgnored
	priceProposalFormatInvalid
)

// x/oracle module sentinel errors
var (
	ErrInvalidMsg                 = sdkerrors.Register(ModuleName, invalidMsg, "invalid input create price")
	ErrPriceProposalIgnored       = sdkerrors.Register(ModuleName, priceProposalIgnored, "price proposal ignored")
	ErrPriceProposalFormatInvalid = sdkerrors.Register(ModuleName, priceProposalFormatInvalid, "price proposal message format invalid")
)
