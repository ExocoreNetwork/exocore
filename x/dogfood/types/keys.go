package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "dogfood"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

const (
	// ExocoreValidatorBytePrefix is the prefix for the validator store.
	ExocoreValidatorBytePrefix byte = iota + 1

	// HistoricalInfoBytePrefix is the prefix for the historical info store.
	HistoricalInfoBytePrefix

	// QueuedOperationsByte is the byte used to store the queue of operations.
	QueuedOperationsByte

	// OptOutsToFinishBytePrefix is the byte used to store the list of operator addresses whose
	// opt outs are maturing at the provided epoch.
	OptOutsToFinishBytePrefix

	// OperatorOptOutFinishEpochBytePrefix is the byte prefix to store the epoch at which an
	// operator's opt out will mature.
	OperatorOptOutFinishEpochBytePrefix

	// ConsensusAddrsToPruneBytePrefix is the byte prefix to store the list of consensus
	// addresses that can be pruned from the operator module at the provided epoch.
	ConsensusAddrsToPruneBytePrefix

	// UnbondingReleaseMaturityBytePrefix is the byte prefix to store the list of undelegations
	// that will mature at the provided epoch.
	UnbondingReleaseMaturityBytePrefix
)

// ExocoreValidatorKey returns the key for the validator store.
func ExocoreValidatorKey(address sdk.AccAddress) []byte {
	return append([]byte{ExocoreValidatorBytePrefix}, address.Bytes()...)
}

// HistoricalInfoKey returns the key for the historical info store.
func HistoricalInfoKey(height int64) []byte {
	bz := sdk.Uint64ToBigEndian(uint64(height))
	return append([]byte{HistoricalInfoBytePrefix}, bz...)
}

// QueuedOperationsKey returns the key for the queued operations store.
func QueuedOperationsKey() []byte {
	return []byte{QueuedOperationsByte}
}

// OptOutsToFinishKey returns the key for the operator opt out maturity store (epoch -> list of
// addresses).
func OptOutsToFinishKey(epoch int64) []byte {
	return append([]byte{OptOutsToFinishBytePrefix}, sdk.Uint64ToBigEndian(uint64(epoch))...)
}

// OperatorOptOutFinishEpochKey is the key for the operator opt out maturity store
// (sdk.AccAddress -> epoch)
func OperatorOptOutFinishEpochKey(address sdk.AccAddress) []byte {
	return append([]byte{OperatorOptOutFinishEpochBytePrefix}, address.Bytes()...)
}

// ConsensusAddrsToPruneKey is the key to lookup the list of operator consensus addresses that
// can be pruned from the operator module at the provided epoch.
func ConsensusAddrsToPruneKey(epoch int64) []byte {
	return append(
		[]byte{ConsensusAddrsToPruneBytePrefix},
		sdk.Uint64ToBigEndian(uint64(epoch))...)
}

// UnbondingReleaseMaturityKey is the key to lookup the list of undelegations that will mature
// at the provided epoch.
func UnbondingReleaseMaturityKey(epoch int64) []byte {
	return append(
		[]byte{UnbondingReleaseMaturityBytePrefix},
		sdk.Uint64ToBigEndian(uint64(epoch))...)
}
