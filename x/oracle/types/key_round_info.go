package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// RoundInfoKeyPrefix is the prefix to retrieve all RoundInfo
	RoundInfoKeyPrefix = "RoundInfo/value/"
)

// RoundInfoKey returns the store key to retrieve a RoundInfo from the index fields
func RoundInfoKey(
	tokenId int32,
) []byte {
	var key []byte

	tokenIdBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tokenIdBytes, uint32(tokenId))
	key = append(key, tokenIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
