package types

import (
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// add for dogfood

// AppendMany appends a variable number of byte slices together
func AppendMany(byteses ...[]byte) (out []byte) {
	for _, bytes := range byteses {
		out = append(out, bytes...)
	}
	return out
}

// ChainIdWithLenKey returns the key with the following format:
// bytePrefix | len(chainId) | chainId
// This is similar to Solidity's ABI encoding.
func ChainIdWithLenKey(chainId string) []byte {
	chainIdL := len(chainId)
	return AppendMany(
		// Append the chainId length
		sdk.Uint64ToBigEndian(uint64(chainIdL)),
		// Append the chainId
		[]byte(chainId),
	)
}

// TMCryptoPublicKeyToConsAddr converts a TM public key to an SDK public key
// and returns the associated consensus address
func TMCryptoPublicKeyToConsAddr(k tmprotocrypto.PublicKey) (sdk.ConsAddress, error) {
	sdkK, err := cryptocodec.FromTmProtoPublicKey(k)
	if err != nil {
		return nil, err
	}
	return sdk.GetConsAddress(sdkK), nil
}
