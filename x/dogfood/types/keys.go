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
