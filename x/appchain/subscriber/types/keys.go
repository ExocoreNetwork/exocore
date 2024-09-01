package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "subscriber"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_subscriber"

	// PortID is the default port id that module binds to
	PortID = "subscriber"

	// SubscriberRedistributeName is the name of the fee pool address that stores
	// the tokens that aren't sent to the coordinator
	SubscriberRedistributeName = "subscriber_redistribute"

	// SubscriberToSendToCoordinatorName is the name of the fee pool address that
	// stores the tokens that are sent to the coordinator
	SubscriberToSendToCoordinatorName = "subscriber_to_send_to_coordinator"
)

const (
	// FirstValsetUpdateID is the first update ID for the validator set
	FirstValsetUpdateID uint64 = 0
)

const (
	// ParamsBytePrefix is the prefix for the params key
	ParamsBytePrefix byte = iota + 1
	// PortBytePrefix is the prefix for the port key
	PortBytePrefix
	// CoordinatorClientIDBytePrefix is the prefix for the coordinator client ID key
	CoordinatorClientIDBytePrefix
	// ValsetUpdateIDBytePrefix is the prefix for the valset update ID key
	ValsetUpdateIDBytePrefix
	// OmniChainValidatorBytePrefix is the prefix for the omni chain validator key
	OmniChainValidatorBytePrefix
)

// ParamsKey returns the key for the params
func ParamsKey() []byte {
	return []byte{ParamsBytePrefix}
}

// PortKey returns the key for the port (hello Harry Potter!)
func PortKey() []byte {
	return []byte{PortBytePrefix}
}

// CoordinatorClientIDKey returns the key for the coordinator client ID
func CoordinatorClientIDKey() []byte {
	return []byte{CoordinatorClientIDBytePrefix}
}

// ValsetUpdateIDKey returns the key for the valset update ID against the provided height.
func ValsetUpdateIDKey(height int64) []byte {
	return append(
		[]byte{ValsetUpdateIDBytePrefix},
		sdk.Uint64ToBigEndian(uint64(height))...,
	)
}

// OmniChainValidatorKey returns the key for the omni chain validator against the provided address.
func OmniChainValidatorKey(address sdk.ConsAddress) []byte {
	return append([]byte{OmniChainValidatorBytePrefix}, address...)
}
