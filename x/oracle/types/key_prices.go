package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PricesKeyPrefix is the prefix to retrieve all Prices
	PricesKeyPrefix = "Prices/value/"
)

// PricesKey returns the store key to retrieve a Prices from the index fields
func PricesKey(
	tokenId int32,
) []byte {
	var key []byte

	tokenIdBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tokenIdBytes, uint32(tokenId))
	key = append(key, tokenIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
