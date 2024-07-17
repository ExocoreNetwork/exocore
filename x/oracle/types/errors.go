package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

const (
	invalidMsg = iota + 2
	priceProposalIgnored
	priceProposalFormatInvalid
	invalidParams
	getPriceFailedAssetNotFound
	getPriceFailedRoundNotFound
)

// x/oracle module sentinel errors
var (
	ErrInvalidMsg                 = sdkerrors.Register(ModuleName, invalidMsg, "invalid input create price")
	ErrPriceProposalIgnored       = sdkerrors.Register(ModuleName, priceProposalIgnored, "price proposal ignored")
	ErrPriceProposalFormatInvalid = sdkerrors.Register(ModuleName, priceProposalFormatInvalid, "price proposal message format invalid")
	ErrInvalidParams              = sdkerrors.Register(ModuleName, invalidParams, "invalid params")
	ErrGetPriceAssetNotFound      = sdkerrors.Register(ModuleName, getPriceFailedAssetNotFound, "get price failed for asset not found")
	ErrGetPriceRoundNotFound      = sdkerrors.Register(ModuleName, getPriceFailedRoundNotFound, "get price failed for round not found")
)
