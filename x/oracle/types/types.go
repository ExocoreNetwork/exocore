package types

import (
	"encoding/binary"

	sdkmath "cosmossdk.io/math"
)

type Price struct {
	Value   sdkmath.Int
	Decimal uint8
}

func Uint64Bytes(value uint64) []byte {
	valueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valueBytes, value)
	return valueBytes
}
