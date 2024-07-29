package types

import (
	"encoding/base64"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"

	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WrappedConsKey is an interface that wraps the different formats of a consensus public key.
// To create an object of this type, use one of the factory NewWrappedConsKeyFrom* functions.
// Since those functions are factory functions, they return the interface type and not the concrete type.
// Note that the address is a compact representation of the public key, and the public key
// cannot be recovered from the address.
type WrappedConsKey interface {
	// ToJSON returns the JSON string representation of the public key. It is used in the CLI.
	ToJSON() string
	// ToHex returns the 32-byte string representation of the public key. It is used in the Bootstrap contract.
	ToHex() string
	// ToTmKey returns the tmprotocrypto (Tendermint format) of the public key.
	ToTmKey() *tmprotocrypto.PublicKey
	// ToSdkKey returns the cryptotypes.PubKey (SDK format) of the public key.
	ToSdkKey() cryptotypes.PubKey
	// ToConsAddr returns the consensus address of the public key.
	ToConsAddr() sdk.ConsAddress
	// Equals returns true if the public key is the same as the other public key.
	Equals(WrappedConsKey) bool
}

// interface guard
var (
	_ WrappedConsKey = &wrappedConsKeyImpl{}
)

const (
	// supportedKeyType is the type of the public key that is supported by the SDK.
	// It must be the value of the `@type` key in the JSON string.
	supportedKeyType = "/cosmos.crypto.ed25519.PubKey"
)

// jsonPubKey is a data structure used to de/serialize consensus public key from/to a JSON string.
type jsonPubKey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

// wrappedConsKeyImpl is a data structure used to store a consensus public key.
// A key may be initiailized from any of the formats below, and it can be
// converted to any of the formats below.
// The storage format is the TmKey, since it supports marshaling and unmarshalling into bytes.
// The format forwarded to Tendermint is TmKey as well, since that is what it uses.
// The format used by our modules (dogfood, appchain) is SdkKey, same as x/staking.
type wrappedConsKeyImpl struct {
	// jsonString is the JSON string representation of the public key, used in the CLI
	// exocored tendermint show-validator
	jsonString string
	// bytes32String is the 32-byte string representation of the public key, used in the Bootstrap contract
	// exocored keys consensus-pubkey-to-bytes
	bytes32String string
	// tmKey is the format used for storage by x/operator and forwarding to Tendermint
	tmKey *tmprotocrypto.PublicKey
	// sdkKey is the format used by all other modules, particularly x/dogfood
	sdkKey cryptotypes.PubKey
	// consAddress cannot be converted back to the public key.
	consAddress sdk.ConsAddress
}

// NewPublicConsKeyFromJSON takes a JSON string and returns a PublicConsKey.
// It validates the jsonStr, and as a side effect, it sets the TmKey field.
// No other fields are set. It only accepts ed25519 keys.
func NewWrappedConsKeyFromJSON(jsonStr string) WrappedConsKey {
	tmKey, err := tmKeyFromJSON(jsonStr)
	if err != nil {
		return nil
	}
	return &wrappedConsKeyImpl{
		tmKey:      tmKey,
		jsonString: jsonStr,
	}
}

// NewPublicConsKeyFromHex takes a hex string and returns a PublicConsKey.
// It validates the key, and as a side effect, it sets the TmKey field.
// No other fields are set. It only accepts ed25519 keys.
func NewWrappedConsKeyFromHex(key string) WrappedConsKey {
	tmKey, err := tmKeyFromHex(key)
	if err != nil {
		return nil
	}
	return &wrappedConsKeyImpl{
		tmKey:         tmKey,
		bytes32String: key,
	}
}

// NewPublicConsKeyFromTmKey takes a tendermint public key and returns a PublicConsKey.
// It only accepts ed25519 keys.
func NewWrappedConsKeyFromTmKey(tmKey *tmprotocrypto.PublicKey) WrappedConsKey {
	if tmKey == nil {
		return nil
	}
	switch tmKey.Sum.(type) {
	case *tmprotocrypto.PublicKey_Ed25519:
		return &wrappedConsKeyImpl{
			tmKey: tmKey,
		}
	default:
		return nil
	}
}

// NewPublicConsKeyFromSdkKey takes an SDK public key and returns a PublicConsKey.
// It validates the key, and as a side effect, it sets the TmKey field.
// No other fields are set. It only accepts ed25519 keys.
func NewWrappedConsKeyFromSdkKey(sdkKey cryptotypes.PubKey) WrappedConsKey {
	// Convert the public key to a tendermint public key.
	tmKey, err := cryptocodec.ToTmProtoPublicKey(sdkKey)
	if err != nil {
		return nil
	}
	// Check if the tmKey so created is an ed25519 key.
	switch tmKey.Sum.(type) {
	case *tmprotocrypto.PublicKey_Ed25519:
		return &wrappedConsKeyImpl{
			tmKey:  &tmKey,
			sdkKey: sdkKey,
		}
	default:
		return nil
	}
}

// ToJSON returns the JSON string representation of the public key.
func (w *wrappedConsKeyImpl) ToJSON() string {
	if w.jsonString == "" {
		w.jsonString = tmKeyToJSON(w.tmKey)
	}
	return w.jsonString
}

// ToHex returns the 32-byte string representation of the public key.
func (w *wrappedConsKeyImpl) ToHex() string {
	if w.bytes32String == "" {
		w.bytes32String = hexutil.Encode(w.tmKey.GetEd25519())
	}
	return w.bytes32String
}

// ToTmKey returns the tendermint public key.
func (w *wrappedConsKeyImpl) ToTmKey() *tmprotocrypto.PublicKey {
	// always initialized, but we return a copy to prevent modification
	cpy := *w.tmKey
	return &cpy
}

// ToSdkKey returns the SDK public key.
func (w *wrappedConsKeyImpl) ToSdkKey() cryptotypes.PubKey {
	if w.sdkKey == nil {
		// #nosec G703 // only errors if key type is unknown, which cannot happen
		sdkKey, _ := cryptocodec.FromTmProtoPublicKey(*w.tmKey)
		w.sdkKey = sdkKey
	}
	return w.sdkKey
}

// ToConsAddr returns the consensus address of the public key.
func (w *wrappedConsKeyImpl) ToConsAddr() sdk.ConsAddress {
	if w.consAddress == nil {
		if w.sdkKey == nil {
			_ = w.ToSdkKey()
		}
		w.consAddress = sdk.GetConsAddress(w.sdkKey)
	}
	return w.consAddress
}

// Equals returns true if the public key is the same as the other public key.
func (w *wrappedConsKeyImpl) Equals(other WrappedConsKey) bool {
	if w == nil || other == nil {
		return false
	}
	// use ToTmKey to compare since it is always initialized
	return w.ToTmKey().Equal(other.ToTmKey())
}

// validateConsensusKey checks that the key is a JSON with `@type` and `key` keys
// with the former bearing exactly value `/cosmos.crypto.ed25519.PubKey`, and the
// latter being a valid base64-encoded public key.
func tmKeyFromJSON(key string) (res *tmprotocrypto.PublicKey, err error) {
	if keyType, keyString, err := base64KeyFromJSON(key); err != nil {
		return nil, errorsmod.Wrap(err, "invalid public key")
	} else if keyType != supportedKeyType {
		return nil, errorsmod.Wrap(ErrParameterInvalid, "invalid public key type")
	} else if res, err = tmKeyFromBase64Key(keyString); err != nil {
		return nil, errorsmod.Wrap(err, "invalid public key")
	}
	return res, nil
}

// tmKeyFromJSON parses the consensus key from a JSON string.
// It returns the key type and the key itself.
// This function replaces deserializing a protobuf any.
func base64KeyFromJSON(jsonStr string) (string, string, error) {
	var pubKey jsonPubKey
	err := json.Unmarshal([]byte(jsonStr), &pubKey)
	if err != nil {
		return "", "", err
	}
	return pubKey.Type, pubKey.Key, nil
}

// tmKeyFromBase64Key converts a base64-encoded public key to a tendermint public key.
// Typically, this function is fed an input from base64KeyFromJSON.
func tmKeyFromBase64Key(pubKey string) (*tmprotocrypto.PublicKey, error) {
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

// tmKeyFromHex converts a 32-byte public key (from the Ethereum side of things),
// which is represented as a 66-byte string (with the 0x prefix) within Golang,
// to a tendermint public key. It is, in effect, a reverse of the command
// `exocored keys consensus-pubkey-to-bytes`
func tmKeyFromHex(key string) (*tmprotocrypto.PublicKey, error) {
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

// consensusKeyToJSON converts a tendermint public key to a JSON string.
// It only supports ed25519 keys. It is the reverse of tmKeyFromJSON
func tmKeyToJSON(key *tmprotocrypto.PublicKey) string {
	pubKey := &jsonPubKey{
		Type: supportedKeyType,
		Key:  base64.StdEncoding.EncodeToString(key.GetEd25519()),
	}
	res, err := json.Marshal(pubKey)
	if err != nil {
		return ""
	}
	return string(res)
}
