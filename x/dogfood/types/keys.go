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

	// ValidatorSetBytePrefix is the prefix for the historical validator set store.
	ValidatorSetBytePrefix

	// ValidatorSetIDBytePrefix is the prefix for the validator set id store.
	ValidatorSetIDBytePrefix

	// HeaderBytePrefix is the prefix for the header store.
	HeaderBytePrefix
)

// ExocoreValidatorKey returns the key for the validator store.
func ExocoreValidatorKey(address sdk.AccAddress) []byte {
	return append([]byte{ExocoreValidatorBytePrefix}, address.Bytes()...)
}

// ValidatorSetKey returns the key for the historical validator set store.
func ValidatorSetKey(id uint64) []byte {
	bz := sdk.Uint64ToBigEndian(id)
	return append([]byte{ValidatorSetBytePrefix}, bz...)
}

// ValidatorSetIDKey returns the key for the validator set id store.
func ValidatorSetIDKey(height int64) []byte {
	bz := sdk.Uint64ToBigEndian(uint64(height))
	return append([]byte{ValidatorSetIDBytePrefix}, bz...)
}

// HeaderKey returns the key for the header store.
func HeaderKey(height int64) []byte {
	bz := sdk.Uint64ToBigEndian(uint64(height))
	return append([]byte{HeaderBytePrefix}, bz...)
}
