package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// RoundDataKeyPrefix is the prefix to retrieve all RoundData
	RoundDataKeyPrefix = "RoundData/value/"
)

// RoundDataKey returns the store key to retrieve a RoundData from the index fields
func RoundDataKey(
	tokenId int32,
) []byte {
	var key []byte

	tokenIdBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tokenIdBytes, uint32(tokenId))
	key = append(key, tokenIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
