package utils

import (
	"bytes"
	"sort"
	"strings"

	"github.com/evmos/evmos/v14/crypto/ethsecp256k1"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// MainnetChainID defines the Evmos EIP155 chain ID for mainnet
	// TODO: the mainnet chainid is still under consideration and need to be finalized later
	MainnetChainID = "exocore_233"
	// TestnetChainID defines the Evmos EIP155 chain ID for testnet
	// TODO: the testnet chainid is still under consideration and need to be finalized later
	TestnetChainID = "exocoretestnet_233"
	// DefaultChainID is the standard chain id used for testing purposes
	DefaultChainID = MainnetChainID + "-1"
	// BaseDenom defines the Evmos mainnet denomination
	BaseDenom = "aexo"
)

// IsMainnet returns true if the chain-id has the Evmos mainnet EIP155 chain prefix.
func IsMainnet(chainID string) bool {
	return strings.HasPrefix(chainID, MainnetChainID)
}

// IsTestnet returns true if the chain-id has the Evmos testnet EIP155 chain prefix.
func IsTestnet(chainID string) bool {
	return strings.HasPrefix(chainID, TestnetChainID)
}

// IsSupportedKey returns true if the pubkey type is supported by the chain
// (i.e eth_secp256k1, amino multisig, ed25519).
// NOTE: Nested multisigs are not supported.
func IsSupportedKey(pubkey cryptotypes.PubKey) bool {
	switch pubkey := pubkey.(type) {
	case *ethsecp256k1.PubKey, *ed25519.PubKey:
		return true
	case multisig.PubKey:
		if len(pubkey.GetPubKeys()) == 0 {
			return false
		}

		for _, pk := range pubkey.GetPubKeys() {
			switch pk.(type) {
			case *ethsecp256k1.PubKey, *ed25519.PubKey:
				continue
			default:
				// Nested multisigs are unsupported
				return false
			}
		}

		return true
	default:
		return false
	}
}

// // GetExocoreAddressFromBech32 returns the sdk.Account address of given address,
// // while also changing bech32 human readable prefix (HRP) to the value set on
// // the global sdk.Config (eg: `evmos`).
// // The function fails if the provided bech32 address is invalid.
// func GetExocoreAddressFromBech32(address string) (sdk.AccAddress, error) {
// 	bech32Prefix := strings.SplitN(address, "1", 2)[0]
// 	if bech32Prefix == address {
// 		return nil, errorsmod.Wrapf(errortypes.ErrInvalidAddress, "invalid bech32 address: %s", address)
// 	}

// 	addressBz, err := sdk.GetFromBech32(address, bech32Prefix)
// 	if err != nil {
// 		return nil, errorsmod.Wrapf(errortypes.ErrInvalidAddress, "invalid address %s, %s", address, err.Error())
// 	}

// 	// safety check: shouldn't happen
// 	if err := sdk.VerifyAddressFormat(addressBz); err != nil {
// 		return nil, err
// 	}

// 	return sdk.AccAddress(addressBz), nil
// }

// func DecodeHexString(hexString string) ([]byte, error) {
// 	if strings.HasPrefix(hexString, "0x") || strings.HasPrefix(hexString, "0X") {
// 		hexString = hexString[2:]
// 	}
// 	if len(hexString)%2 != 0 {
// 		hexString = "0" + hexString
// 	}
// 	return hex.DecodeString(hexString)
// }

// // ProcessAddress converts a hex address into the bech32 address format.
// // If the input address is already in bech32 format, it returns the same address.
// func ProcessAddress(address string) (string, error) {
// 	switch {
// 	case common.IsHexAddress(address):
// 		b := common.FromHex(address)
// 		encodedAddress, err := bech32.EncodeFromBase256(config.Bech32Prefix, b)
// 		if err != nil {
// 			return "", err
// 		}
// 		return encodedAddress, nil
// 	case strings.HasPrefix(address, config.Bech32Prefix):
// 		return address, nil
// 	default:
// 		return "", errorsmod.Wrapf(errortypes.ErrInvalidAddress, "invalid input address: %s", address)
// 	}
// }

// SortByPower sorts operators, their pubkeys, and their powers by the powers.
// the sorting is descending, so the highest power is first. If the powers are equal,
// the operator address (bytes, not string!) is used as a tiebreaker. The bytes
// are preferred since that is how the operator module stores them, indexed by
// the bytes.
func SortByPower(
	operatorAddrs []sdk.AccAddress,
	pubKeys []operatortypes.WrappedConsKey,
	powers []int64,
) ([]sdk.AccAddress, []operatortypes.WrappedConsKey, []int64) {
	// Create a slice of indices
	indices := make([]int, len(powers))
	for i := range indices {
		indices[i] = i
	}

	// Sort the indices slice based on the powers slice
	sort.SliceStable(indices, func(i, j int) bool {
		if powers[indices[i]] == powers[indices[j]] {
			// ascending order of operator address as tiebreaker
			return bytes.Compare(operatorAddrs[indices[i]], operatorAddrs[indices[j]]) < 0
		}
		// descending order of power
		return powers[indices[i]] > powers[indices[j]]
	})

	// Reorder all slices using the sorted indices
	sortedOperatorAddrs := make([]sdk.AccAddress, len(operatorAddrs))
	sortedPubKeys := make([]operatortypes.WrappedConsKey, len(pubKeys))
	sortedPowers := make([]int64, len(powers))
	for i, idx := range indices {
		sortedOperatorAddrs[i] = operatorAddrs[idx]
		sortedPubKeys[i] = pubKeys[idx]
		sortedPowers[i] = powers[idx]
	}
	return sortedOperatorAddrs, sortedPubKeys, sortedPowers
}
