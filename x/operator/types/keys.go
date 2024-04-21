package types

import (
	"math"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// constants
const (
	// ModuleName module name
	ModuleName = "operator"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName

	DefaultOptedOutHeight = uint64(math.MaxUint64)

	USDValueDefaultDecimal = uint8(8)

	SlashVetoDuration = int64(1000)

	UnbondingExpiration = 10
)

const (
	prefixOperatorInfo = iota + 1

	prefixOperatorOptedAVSInfo

	prefixVotingPowerForAVSOperator

	prefixOperatorSlashInfo

	prefixSlashAssetsState

	// add keys for dogfood
	BytePrefixForOperatorAndChainIDToConsKey = iota
	BytePrefixForOperatorAndChainIDToPrevConsKey
	BytePrefixForChainIDAndOperatorToConsKey
	BytePrefixForChainIDAndConsKeyToOperator
	BytePrefixForOperatorOptOutFromChainID
)

var (
	// KeyPrefixOperatorInfo key-value: operatorAddr->types.OperatorInfo
	KeyPrefixOperatorInfo = []byte{prefixOperatorInfo}

	// KeyPrefixOperatorOptedAVSInfo key-value:
	// operatorAddr + '/' + AVSAddr -> OptedInfo
	KeyPrefixOperatorOptedAVSInfo = []byte{prefixOperatorOptedAVSInfo}

	// KeyPrefixVotingPowerForAVSOperator key-value:
	// AVSAddr -> types.DecValueField（the voting power of specified Avs）
	// AVSAddr + '/' + operatorAddr -> types.DecValueField (the voting power of specified operator and Avs)
	KeyPrefixVotingPowerForAVSOperator = []byte{prefixVotingPowerForAVSOperator}

	// KeyPrefixOperatorSlashInfo key-value:
	// operator + '/' + AVSAddr + '/' + slashId -> OperatorSlashInfo
	KeyPrefixOperatorSlashInfo = []byte{prefixOperatorSlashInfo}

	// KeyPrefixSlashAssetsState key-value:
	// processedSlashHeight + '/' + assetID -> SlashAmount
	// processedSlashHeight + '/' + assetID + '/' + stakerID -> SlashAmount
	// processedSlashHeight + '/' + assetID + '/' + operatorAddr -> SlashAmount
	KeyPrefixSlashAssetsState = []byte{prefixSlashAssetsState}
)

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

func AddrAndChainIDKey(prefix byte, addr sdk.AccAddress, chainID string) []byte {
	partialKey := ChainIDWithLenKey(chainID)
	return AppendMany(
		// Append the prefix
		[]byte{prefix},
		// Append the addr bytes first so we can iterate over all chain ids
		// belonging to an operator easily.
		addr,
		// Append the partialKey
		partialKey,
	)
}

func ChainIDAndAddrKey(prefix byte, chainID string, addr sdk.AccAddress) []byte {
	partialKey := ChainIDWithLenKey(chainID)
	return AppendMany(
		// Append the prefix
		[]byte{prefix},
		// Append the partialKey so that we can look for any operator keys
		// corresponding to this chainID easily.
		partialKey,
		addr,
	)
}

func KeyForOperatorAndChainIDToConsKey(addr sdk.AccAddress, chainID string) []byte {
	return AddrAndChainIDKey(
		BytePrefixForOperatorAndChainIDToConsKey,
		addr, chainID,
	)
}

func KeyForOperatorAndChainIDToPrevConsKey(addr sdk.AccAddress, chainID string) []byte {
	return AddrAndChainIDKey(
		BytePrefixForOperatorAndChainIDToPrevConsKey,
		addr, chainID,
	)
}

func KeyForChainIDAndOperatorToConsKey(chainID string, addr sdk.AccAddress) []byte {
	return ChainIDAndAddrKey(
		BytePrefixForChainIDAndOperatorToConsKey,
		chainID, addr,
	)
}

func KeyForChainIDAndConsKeyToOperator(chainID string, addr sdk.ConsAddress) []byte {
	return AppendMany(
		[]byte{BytePrefixForChainIDAndConsKeyToOperator},
		ChainIDWithLenKey(chainID),
		addr,
	)
}

func KeyForOperatorOptOutFromChainID(addr sdk.AccAddress, chainID string) []byte {
	return AppendMany(
		[]byte{BytePrefixForOperatorOptOutFromChainID}, addr,
		ChainIDWithLenKey(chainID),
	)
}

func IterateOperatorsForAVSPrefix(avsAddr string) []byte {
	tmp := append([]byte(avsAddr), '/')
	return tmp
}
