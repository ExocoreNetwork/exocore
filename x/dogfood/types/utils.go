package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// TMCryptoPublicKeyToConsAddr converts a TM public key to an SDK public key
// and returns the associated consensus address.
func TMCryptoPublicKeyToConsAddr(k tmprotocrypto.PublicKey) (sdk.ConsAddress, error) {
	sdkK, err := cryptocodec.FromTmProtoPublicKey(k)
	if err != nil {
		return nil, err
	}
	return sdk.GetConsAddress(sdkK), nil
}

// RemoveFromBytesList removes an address from a list of addresses
// or a byte slice from a list of byte slices.
func RemoveFromBytesList(list [][]byte, addr []byte) [][]byte {
	for i, a := range list {
		if bytes.Equal(a, addr) {
			return append(list[:i], list[i+1:]...)
		}
	}
	panic("address not found in list")
}
