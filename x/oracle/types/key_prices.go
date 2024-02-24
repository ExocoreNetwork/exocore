package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PricesKeyPrefix is the prefix to retrieve all Prices
	PricesKeyPrefix = "Prices/value/"
)

// PricesNextRoundIdKey is the key set for each tokenId storeKV to store the next round id
var PricesNextRountIdKey = []byte("tokenId/")

// PricesKey returns the store key to retrieve a Prices from the index fields
// this key is actually used as the prefix for kvsotre, TODO: refactor to PriceTokenPrefix
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

// PricesRoundKey returns the store key to retrieve a PriceWithTimeAndRound from the index fields
func PricesRoundKey(
	roundId uint64,
) []byte {
	return append(Uint64Bytes(roundId), []byte("/")...)
}
