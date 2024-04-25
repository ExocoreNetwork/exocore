package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// PricesKeyPrefix is the prefix to retrieve all Prices
	PricesKeyPrefix = "Prices/value/"
)

// PricesNextRoundIDKey is the key set for each tokenId storeKV to store the next round id
var PricesNextRoundIDKey = []byte("nextRoundID/")

// PricesKey returns the store key to retrieve a Prices from the index fields
// this key is actually used as the prefix for kvsotre, TODO: refactor to PriceTokenPrefix
func PricesKey(
	tokenID uint64,
) []byte {
	var key []byte

	tokenIDBytes := Uint64Bytes(tokenID)
	key = append(key, tokenIDBytes...)
	key = append(key, []byte("/")...)

	return key
}

// PricesRoundKey returns the store key to retrieve a PriceTimeRound from the index fields
func PricesRoundKey(
	roundID uint64,
) []byte {
	return append(Uint64Bytes(roundID), []byte("/")...)
}
