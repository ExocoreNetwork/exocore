package types

import "strings"

const (
	// NativeTokenKeyPrefix is the prefix to retrieve all NativeToken
	NativeTokenKeyPrefix                 = "NativeToken/"
	NativeTokenPriceKeyPrefix            = NativeTokenKeyPrefix + "price/value/"
	NativeTokenStakerInfoKeyPrefix       = NativeTokenKeyPrefix + "stakerInfo/value/"
	NativeTokenOperatorInfoKeyPrefix     = NativeTokenKeyPrefix + "operatorInfo/value/"
	NativeTokenStakerListKeyPrefix       = NativeTokenKeyPrefix + "stakerList/value/"
	NativeTokenStakerDelegationKeyPrefix = NativeTokenKeyPrefix + "stakerDelegation/value/"
)

func NativeTokenStakerDelegationKey(assetID, stakerAddr, operatorAddr string) []byte {
	return append([]byte(NativeTokenStakerDelegationKeyPrefix), []byte(assetID)...)
}

func NativeTokenStakerListKey(assetID string) []byte {
	return append([]byte(NativeTokenStakerListKeyPrefix), []byte(assetID)...)
}

func NativeTokenStakerKey(assetID, stakerAddr string) []byte {
	assetID = strings.Join([]string{assetID, stakerAddr}, "/")
	return append([]byte(NativeTokenStakerInfoKeyPrefix), []byte(assetID)...)
}

func NativeTokenOperatorKey(assetID, operatorAddr string) []byte {
	assetID = strings.Join([]string{assetID, operatorAddr}, "/")
	return append([]byte(NativeTokenOperatorInfoKeyPrefix), []byte(assetID)...)
}
