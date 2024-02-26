package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ValidatorsKeyPrefix is the prefix to retrieve all Validators
	ValidatorsKeyPrefix = "Validators/value/"
)

// ValidatorsKey returns the store key to retrieve a Validators from the index fields
func ValidatorsKey(
	block uint64,
) []byte {
	var key []byte

	blockBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(blockBytes, block)
	key = append(key, blockBytes...)
	key = append(key, []byte("/")...)

	return key
}
