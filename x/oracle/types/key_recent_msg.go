package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// RecentMsgKeyPrefix is the prefix to retrieve all RecentMsg
	RecentMsgKeyPrefix = "RecentMsg/value/"
)

// RecentMsgKey returns the store key to retrieve a RecentMsg from the index fields
func RecentMsgKey(
	block uint64,
) []byte {
	var key []byte

	blockBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(blockBytes, block)
	key = append(key, blockBytes...)
	key = append(key, []byte("/")...)

	return key
}
