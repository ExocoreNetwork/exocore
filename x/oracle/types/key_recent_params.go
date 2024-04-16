package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// RecentParamsKeyPrefix is the prefix to retrieve all RecentParams
	RecentParamsKeyPrefix = "RecentParams/value/"
)

// RecentParamsKey returns the store key to retrieve a RecentParams from the index fields
func RecentParamsKey(
	block uint64,
) []byte {
	var key []byte

	blockBytes := Uint64Bytes(block)
	key = append(key, blockBytes...)
	key = append(key, []byte("/")...)

	return key
}
