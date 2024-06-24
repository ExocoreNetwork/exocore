package types

import (
	"encoding/base64"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/ethereum/go-ethereum/common/hexutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HexStringToPubKey converts a 32-byte public key (from the Ethereum side of things),
// which is represented as a 66-byte string (with the 0x prefix) within Golang,
// to a tendermint public key. It is, in effect, a reverse of the command
// `exocored keys consensus-pubkey-to-bytes`
func HexStringToPubKey(key string) (*tmprotocrypto.PublicKey, error) {
	if len(key) != 66 {
		return nil, errorsmod.Wrapf(
			ErrInvalidPubKey,
			"expected 66 length string, got %d",
			len(key),
		)
	}
	keyBytes, err := hexutil.Decode(key)
	if err != nil {
		return nil, errorsmod.Wrapf(
			ErrInvalidPubKey,
			"failed to decode hex string: %s",
			err,
		)
	}
	return &tmprotocrypto.PublicKey{
		Sum: &tmprotocrypto.PublicKey_Ed25519{
			Ed25519: keyBytes,
		},
	}, nil
}

// StringToPubKey converts a base64-encoded public key to a tendermint public key.
// Typically, this function is fed an input from ParseConsensusKeyFromJSON.
func StringToPubKey(pubKey string) (*tmprotocrypto.PublicKey, error) {
	pubKeyBytes, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	return &tmprotocrypto.PublicKey{
		Sum: &tmprotocrypto.PublicKey_Ed25519{
			Ed25519: pubKeyBytes,
		},
	}, nil
}

// ParseConsensusKeyFromJSON parses the consensus key from a JSON string.
// It returns the key type and the key itself.
// This function replaces deserializing a protobuf any.
func ParseConsensusKeyFromJSON(jsonStr string) (string, string, error) {
	type PubKey struct {
		Type string `json:"@type"`
		Key  string `json:"key"`
	}
	var pubKey PubKey
	err := json.Unmarshal([]byte(jsonStr), &pubKey)
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

// ValidateConsensusKey checks that the key is a JSON with `@type` and `key` keys
// with the former bearing exactly value `/cosmos.crypto.ed25519.PubKey`, and the
// latter being a valid base64-encoded public key.
func ValidateConsensusKeyJSON(key string) (res *tmprotocrypto.PublicKey, err error) {
	if keyType, keyString, err := ParseConsensusKeyFromJSON(key); err != nil {
		return nil, errorsmod.Wrap(err, "invalid public key")
	} else if keyType != "/cosmos.crypto.ed25519.PubKey" {
		return nil, errorsmod.Wrap(ErrParameterInvalid, "invalid public key type")
	} else if res, err = StringToPubKey(keyString); err != nil {
		return nil, errorsmod.Wrap(err, "invalid public key")
	}
	return res, nil
}
