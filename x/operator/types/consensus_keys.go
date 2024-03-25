package types

import (
	"encoding/base64"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bytes32ToPubKey converts a 32-byte public key (from the Ethereum side of things)
// to a tendermint public key. It is, in effect, a reverse of the command
// `exocored keys consensus-pubkey-to-bytes`
func Bytes32ToPubKey(pubKeyBytes []byte) (*tmprotocrypto.PublicKey, error) {
	if len(pubKeyBytes) != 32 {
		return nil, errorsmod.Wrapf(
			ErrInvalidPubKeyBytes,
			"expected 32 bytes, got %d",
			len(pubKeyBytes),
		)
	}
	return &tmprotocrypto.PublicKey{
		Sum: &tmprotocrypto.PublicKey_Ed25519{
			Ed25519: pubKeyBytes,
		},
	}, nil
}

// StringToPubKey converts a base64-encoded public key to a tendermint public key.
// Typically, this function is fed an input from ParseConsumerKeyFromJson.
func StringToPubKey(pubKey string) (key tmprotocrypto.PublicKey, err error) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return
	}
	subscriberTMConsKey := tmprotocrypto.PublicKey{
		Sum: &tmprotocrypto.PublicKey_Ed25519{
			Ed25519: pubKeyBytes,
		},
	}
	return subscriberTMConsKey, nil
}

// ParseConsumerKeyFromJson parses the consumer key from a JSON string.
// This function replaces deserializing a protobuf any.
func ParseConsumerKeyFromJson(jsonStr string) (pkType, key string, err error) {
	type PubKey struct {
		Type string `json:"type"`
		Key  string `json:"value"`
	}
	var pubKey PubKey
	err = json.Unmarshal([]byte(jsonStr), &pubKey)
	if err != nil {
		return "", "", err
	}
	return pubKey.Type, pubKey.Key, nil
}

// TMCryptoPublicKeyToConsAddr converts a TM public key to an SDK public key
// and returns the associated consensus address
func TMCryptoPublicKeyToConsAddr(k *tmprotocrypto.PublicKey) (sdk.ConsAddress, error) {
	sdkK, err := cryptocodec.FromTmProtoPublicKey(*k)
	if err != nil {
		return nil, err
	}
	return sdk.GetConsAddress(sdkK), nil
}
