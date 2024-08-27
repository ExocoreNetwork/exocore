package types

import (
	errorsmod "cosmossdk.io/errors"
)

const (
	errCodeInvalidParams = iota + 2
	errCodeNilRequest
	errCodeDuplicateSubChain
	errCodeNoOperators
)

var (
	// ErrInvalidRegistrationParams is the error returned when the subscriber chain registration params are invalid
	ErrInvalidRegistrationParams = errorsmod.Register(ModuleName, errCodeInvalidParams, "invalid registration params")
	// ErrNilRequest is the error returned when the request is nil
	ErrNilRequest = errorsmod.Register(ModuleName, errCodeNilRequest, "nil request")
	// ErrDuplicateSubChain is the error returned when a client for the chain already exists
	ErrDuplicateSubChain = errorsmod.Register(ModuleName, errCodeDuplicateSubChain, "subscriber chain already exists")
	// ErrNoOperators is the error returned when no qualified operators are available
	ErrNoOperators = errorsmod.Register(ModuleName, errCodeNoOperators, "no operators available")
)
