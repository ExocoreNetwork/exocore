package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"math"
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

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

const (
	prefixOperatorInfo = iota + 1

	prefixOperatorOptedAVSInfo

	prefixAVSOperatorAssetsTotalValue

	prefixOperatorAVSSingleAssetState

	prefixOperatorAVSStakerShareState

	prefixOperatorSlashInfo

	prefixSlashAssetsState

	//add keys for dogfood
	BytePrefixForOperatorAndChainIdToConsKey = iota
	BytePrefixForOperatorAndChainIdToPrevConsKey
	BytePrefixForChainIdAndOperatorToConsKey
	BytePrefixForChainIdAndConsKeyToOperator
	BytePrefixForOperatorOptOutFromChainId
)

var (
	// KeyPrefixOperatorInfo key-value: operatorAddr->operatorInfo
	KeyPrefixOperatorInfo = []byte{prefixOperatorInfo}

	// KeyPrefixOperatorOptedAVSInfo key-value:
	// operatorAddr + '/' + AVSAddr -> OptedInfo
	KeyPrefixOperatorOptedAVSInfo = []byte{prefixOperatorOptedAVSInfo}

	// KeyPrefixAVSOperatorAssetsTotalValue key-value:
	// AVSAddr -> AVSTotalValue
	// AVSAddr + '/' + operatorAddr -> AVSOperatorTotalValue
	KeyPrefixAVSOperatorAssetsTotalValue = []byte{prefixAVSOperatorAssetsTotalValue}

	// KeyPrefixOperatorAVSSingleAssetState key-value:
	// assetId + '/' + AVSAddr + '/' + operatorAddr -> AssetOptedInState
	KeyPrefixOperatorAVSSingleAssetState = []byte{prefixOperatorAVSSingleAssetState}

	// KeyPrefixAVSOperatorStakerShareState key-value:
	// AVSAddr + '/' + '' + '/' +  operatorAddr -> ownAssetsOptedInValue
	// AVSAddr + '/' + stakerId + '/' + operatorAddr -> assetsOptedInValue
	KeyPrefixAVSOperatorStakerShareState = []byte{prefixOperatorAVSStakerShareState}

	// KeyPrefixOperatorSlashInfo key-value:
	// operator + '/' + AVSAddr + '/' + slashId -> OperatorSlashInfo
	KeyPrefixOperatorSlashInfo = []byte{prefixOperatorSlashInfo}

	// KeyPrefixSlashAssetsState key-value:
	// completeSlashHeight + '/' + assetId -> SlashAmount
	// completeSlashHeight + '/' + assetId + '/' + stakerId -> SlashAmount
	// completeSlashHeight + '/' + assetId + '/' + operatorAddr -> SlashAmount
	KeyPrefixSlashAssetsState = []byte{prefixSlashAssetsState}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func AddrAndChainIdKey(prefix byte, addr sdk.AccAddress, chainId string) []byte {
	partialKey := ChainIdWithLenKey(chainId)
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

func ChainIdAndAddrKey(prefix byte, chainId string, addr sdk.AccAddress) []byte {
	partialKey := ChainIdWithLenKey(chainId)
	return AppendMany(
		// Append the prefix
		[]byte{prefix},
		// Append the partialKey so that we can look for any operator keys
		// corresponding to this chainId easily.
		partialKey,
		addr,
	)
}

func KeyForOperatorAndChainIdToConsKey(addr sdk.AccAddress, chainId string) []byte {
	return AddrAndChainIdKey(
		BytePrefixForOperatorAndChainIdToConsKey,
		addr, chainId,
	)
}

func KeyForOperatorAndChainIdToPrevConsKey(addr sdk.AccAddress, chainId string) []byte {
	return AddrAndChainIdKey(
		BytePrefixForOperatorAndChainIdToPrevConsKey,
		addr, chainId,
	)
}

func KeyForChainIdAndOperatorToConsKey(chainId string, addr sdk.AccAddress) []byte {
	return ChainIdAndAddrKey(
		BytePrefixForChainIdAndOperatorToConsKey,
		chainId, addr,
	)
}

func KeyForChainIdAndConsKeyToOperator(chainId string, addr sdk.ConsAddress) []byte {
	return AppendMany(
		[]byte{BytePrefixForChainIdAndConsKeyToOperator},
		ChainIdWithLenKey(chainId),
		addr,
	)
}

func KeyForOperatorOptOutFromChainId(addr sdk.AccAddress, chainId string) []byte {
	return AppendMany(
		[]byte{BytePrefixForOperatorOptOutFromChainId}, addr,
		ChainIdWithLenKey(chainId),
	)
}
