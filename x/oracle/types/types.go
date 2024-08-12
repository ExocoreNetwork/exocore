package types

import (
	"encoding/binary"

	sdkmath "cosmossdk.io/math"
)

type OracleInfo struct {
	Token struct {
		Name  string `json:"name"`
		Chain struct {
			Name string `json:"name"`
			Desc string `json:"desc"`
		} `json:"chain"`
		Decimal  string `json:"decimal"`
		Contract string `json:"contract"`
		AssetID  string `json:"asset_id"`
	} `json:"token"`
	Feeder struct {
		Start    string `json:"start"`
		End      string `json:"end"`
		Interval string `json:"interval"`
		RuleID   string `json:rule_id"`
	} `json:"feeder"`
	AssetID string `json:"asset_id"`
}

type Price struct {
	Value   sdkmath.Int
	Decimal uint8
}

const (
	DefaultPriceValue   = 1
	DefaultPriceDecimal = 0
)

func Uint64Bytes(value uint64) []byte {
	valueBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(valueBytes, value)
	return valueBytes
}
