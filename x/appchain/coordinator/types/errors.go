package types

import (
	errorsmod "cosmossdk.io/errors"
)

const (
	errCodeInvalidParams = iota + 2
	errCodeNilRequest
	errCodeDuplicateSubChain
	errCodeNoOperators
	errCodeInvalidSubscriberClient
	errCodeUnknownSubscriberChannelID
)

var (
	// ErrInvalidRegistrationParams is the error returned when the subscriber chain
	// registration params are invalid
	ErrInvalidRegistrationParams = errorsmod.Register(
		ModuleName, errCodeInvalidParams, "invalid registration params",
	)
	// ErrNilRequest is the error returned when the request is nil
	ErrNilRequest = errorsmod.Register(
		ModuleName, errCodeNilRequest, "nil request",
	)
	// ErrDuplicateSubChain is the error returned when
	// a client for the chain already exists
	ErrDuplicateSubChain = errorsmod.Register(
		ModuleName, errCodeDuplicateSubChain, "subscriber chain already exists",
	)
	// ErrNoOperators is the error returned when no qualified operators are available
	ErrNoOperators = errorsmod.Register(
		ModuleName, errCodeNoOperators, "no operators available",
	)
	// ErrInvalidSubscriberClient is the error returned when the
	// client for the subscriber chain is invalid
	ErrInvalidSubscriberClient = errorsmod.Register(
		ModuleName, errCodeInvalidSubscriberClient, "invalid subscriber client",
	)
	// ErrUnknownSubscriberChannelID is the error returned when the channel ID
	// corresponding to a message from the subscriber chain is unknown
	ErrUnknownSubscriberChannelID = errorsmod.Register(
		ModuleName, errCodeUnknownSubscriberChannelID, "unknown subscriber channel ID",
	)
)
