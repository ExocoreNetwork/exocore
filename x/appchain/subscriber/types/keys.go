package types

import (
	time "time"

	"github.com/ExocoreNetwork/exocore/utils"
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
	// SubscriberValidatorBytePrefix is the prefix for the subscriber validator key
	SubscriberValidatorBytePrefix
	// CoordinatorChannelBytePrefix is the prefix for the coordinator channel key
	CoordinatorChannelBytePrefix
	// PendingChangesBytePrefix is the prefix for the pending changes key
	PendingChangesBytePrefix
	// PacketMaturityTimeBytePrefix is the prefix for the packet maturity time key
	PacketMaturityTimeBytePrefix
	// OutstandingDowntimeBytePrefix is the prefix for the outstanding downtime key
	OutstandingDowntimeBytePrefix
	// PendingPacketsIndexBytePrefix is the prefix for the pending packets index key
	PendingPacketsIndexBytePrefix
	// PendingDataPacketsBytePrefix is the prefix for the pending data packets key
	PendingDataPacketsBytePrefix
	// HistoricalInfoBytePrefix is the prefix for the historical info key
	HistoricalInfoBytePrefix
	// LastRewardTransmissionHeightBytePrefix is the prefix for the last reward transmission height key
	LastRewardTransmissionHeightBytePrefix
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

// SubscriberValidatorKey returns the key for the subscriber chain validator
// against the provided address.
func SubscriberValidatorKey(address sdk.ConsAddress) []byte {
	return append([]byte{SubscriberValidatorBytePrefix}, address...)
}

// CoordinatorChannelKey returns the key for which the ibc channel id to the coordinator chain
// is stored.
func CoordinatorChannelKey() []byte {
	return []byte{CoordinatorChannelBytePrefix}
}

// PendingChangesKey returns the key for the pending changes
func PendingChangesKey() []byte {
	return []byte{PendingChangesBytePrefix}
}

// PacketMaturityTimeKey returns the key for the packet maturity time
func PacketMaturityTimeKey(maturityTime time.Time, vscID uint64) []byte {
	return utils.AppendMany(
		[]byte{PacketMaturityTimeBytePrefix},
		sdk.FormatTimeBytes(maturityTime),
		sdk.Uint64ToBigEndian(vscID),
	)
}

// OutstandingDowntimeKey returns the key for the outstanding downtime
func OutstandingDowntimeKey(consAddress sdk.ConsAddress) []byte {
	return append([]byte{OutstandingDowntimeBytePrefix}, consAddress.Bytes()...)
}

// PendingPacketsIndexKey returns the key for the pending packets index
func PendingPacketsIndexKey() []byte {
	return []byte{PendingPacketsIndexBytePrefix}
}

// PendingDataPacketsKey returns the key for the pending data packets
func PendingDataPacketsKey(idx uint64) []byte {
	return append([]byte{PendingDataPacketsBytePrefix}, sdk.Uint64ToBigEndian(idx)...)
}

// HistoricalInfoKey returns the key for the historical info
func HistoricalInfoKey(height int64) []byte {
	return append([]byte{HistoricalInfoBytePrefix}, sdk.Uint64ToBigEndian(uint64(height))...)
}

// LastRewardTransmissionHeightKey is the key for the last reward transmission height
func LastRewardTransmissionHeightKey() []byte {
	return []byte{LastRewardTransmissionHeightBytePrefix}
}
