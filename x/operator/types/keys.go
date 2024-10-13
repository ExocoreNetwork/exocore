package types

import (
	"math"

	"github.com/ExocoreNetwork/exocore/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/xerrors"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
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

	SlashVetoDuration = int64(1000)

	UnbondingExpiration = 10

	// AccAddressLength is used to parse the key, because the length isn't padded in the key
	// This might be removed if the address length is padded in the key
	AccAddressLength = 20

	// ByteLengthForUint64 the type of chainID length is uint64, uint64 has 8 bytes.
	ByteLengthForUint64 = 8
)

const (
	prefixOperatorInfo = iota + 1

	prefixOperatorOptedAVSInfo

	prefixUSDValueForAVS
	prefixUSDValueForOperator

	prefixOperatorSlashInfo

	prefixSlashAssetsState

	// BytePrefixForOperatorAndChainIDToConsKey is the prefix to store the consensus key for
	// an operator for a chainID.
	BytePrefixForOperatorAndChainIDToConsKey

	// BytePrefixForOperatorAndChainIDToPrevConsKey is the prefix to store the previous
	// consensus key for an operator for a chainID.
	BytePrefixForOperatorAndChainIDToPrevConsKey

	// BytePrefixForChainIDAndOperatorToConsKey is the prefix to store the reverse lookup for
	// a chainID + operator address to the consensus key.
	BytePrefixForChainIDAndOperatorToConsKey

	// BytePrefixForChainIDAndConsKeyToOperator is the prefix to store the reverse lookup for
	// a chainID + consensus key to the operator address.
	BytePrefixForChainIDAndConsKeyToOperator

	// BytePrefixForOperatorKeyRemovalForChainID is the prefix to store that the operator with
	// the given address is in the process of unbonding their key for the given chainID.
	BytePrefixForOperatorKeyRemovalForChainID
)

var (
	// KeyPrefixOperatorInfo key-value: operatorAddr->types.OperatorInfo
	KeyPrefixOperatorInfo = []byte{prefixOperatorInfo}

	// KeyPrefixOperatorOptedAVSInfo key-value:
	// operatorAddr + '/' + AVSAddr -> OptedInfo
	KeyPrefixOperatorOptedAVSInfo = []byte{prefixOperatorOptedAVSInfo}

	// KeyPrefixUSDValueForAVS key-value:
	// AVSAddr -> types.DecValueField（the USD value of specified Avs）
	KeyPrefixUSDValueForAVS = []byte{prefixUSDValueForAVS}

	// KeyPrefixUSDValueForOperator key-value:
	// AVSAddr + '/' + operatorAddr -> types.OperatorOptedUSDValue (the voting power of specified operator and Avs)
	KeyPrefixUSDValueForOperator = []byte{prefixUSDValueForOperator}

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
	partialKey := utils.ChainIDWithLenKey(chainID)
	return utils.AppendMany(
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
	partialKey := utils.ChainIDWithLenKey(chainID)
	return utils.AppendMany(
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

func ParseKeyForOperatorAndChainIDToConsKey(key []byte) (addr sdk.AccAddress, chainID string, err error) {
	if len(key) < AccAddressLength+ByteLengthForUint64 {
		return nil, "", xerrors.New("key length is too short to contain address and chainID length")
	}
	// Extract the address
	addr = key[0:AccAddressLength]
	if len(addr) == 0 {
		return nil, "", xerrors.New("missing address")
	}

	// Extract the chainID length
	chainIDLen := sdk.BigEndianToUint64(key[AccAddressLength : AccAddressLength+ByteLengthForUint64])
	if len(key) != int(AccAddressLength+ByteLengthForUint64+chainIDLen) {
		return nil, "", xerrors.Errorf("invalid key length,expected:%d,got:%d", AccAddressLength+ByteLengthForUint64+chainIDLen, len(key))
	}

	// Extract the chainID
	chainIDBytes := key[AccAddressLength+ByteLengthForUint64:]
	chainID = string(chainIDBytes)

	return addr, chainID, nil
}

func KeyForChainIDAndOperatorToPrevConsKey(chainID string, addr sdk.AccAddress) []byte {
	return ChainIDAndAddrKey(
		BytePrefixForOperatorAndChainIDToPrevConsKey,
		chainID, addr,
	)
}

func ParsePrevConsKey(key []byte) (chainID string, addr sdk.AccAddress, err error) {
	// Check if the key has at least eight byte for the chainID length
	if len(key) < ByteLengthForUint64 {
		return "", nil, xerrors.New("key length is too short to contain chainID length")
	}

	// Extract the chainID length
	chainIDLen := sdk.BigEndianToUint64(key[0:ByteLengthForUint64])
	if len(key) < int(ByteLengthForUint64+chainIDLen) {
		return "", nil, xerrors.New("key too short for chainID length")
	}

	// Extract the chainID
	chainIDBytes := key[ByteLengthForUint64 : ByteLengthForUint64+chainIDLen]
	chainID = string(chainIDBytes)

	// Extract the address
	addr = key[ByteLengthForUint64+chainIDLen:]
	if len(addr) == 0 {
		return "", nil, xerrors.New("missing address")
	}

	return chainID, addr, nil
}

func KeyForChainIDAndOperatorToConsKey(chainID string, addr sdk.AccAddress) []byte {
	return ChainIDAndAddrKey(
		BytePrefixForChainIDAndOperatorToConsKey,
		chainID, addr,
	)
}

func KeyForChainIDAndConsKeyToOperator(chainID string, addr sdk.ConsAddress) []byte {
	return utils.AppendMany(
		[]byte{BytePrefixForChainIDAndConsKeyToOperator},
		utils.ChainIDWithLenKey(chainID),
		addr,
	)
}

func KeyForOperatorKeyRemovalForChainID(addr sdk.AccAddress, chainID string) []byte {
	return utils.AppendMany(
		[]byte{BytePrefixForOperatorKeyRemovalForChainID},
		addr,
		// TODO: it may be possible to just use the chainID here without the length.
		// This is because the chainID is at the end of the key and we can just iterate
		// over all keys with the same operator address.
		utils.ChainIDWithLenKey(chainID),
	)
}

func ParseKeyForOperatorKeyRemoval(key []byte) (addr sdk.AccAddress, chainID string, err error) {
	// Check if the key has at least 20 byte for the operator and eight byte for the chainID length
	if len(key) < AccAddressLength+ByteLengthForUint64 {
		return nil, "", xerrors.New("key length is too short to contain operator address and chainID length")
	}

	// Extract the address
	addr = key[0:AccAddressLength]
	if len(addr) == 0 {
		return nil, "", xerrors.New("missing address")
	}

	// Extract the chainID length
	chainIDLen := sdk.BigEndianToUint64(key[AccAddressLength : AccAddressLength+ByteLengthForUint64])
	if len(key) != int(AccAddressLength+ByteLengthForUint64+chainIDLen) {
		return nil, "", xerrors.Errorf("invalid key length,expected:%d,got:%d", AccAddressLength+ByteLengthForUint64+chainIDLen, len(key))
	}

	// Extract the chainID
	chainIDBytes := key[AccAddressLength+ByteLengthForUint64:]
	chainID = string(chainIDBytes)

	return addr, chainID, nil
}

func IterateOperatorsForAVSPrefix(avsAddr string) []byte {
	tmp := append([]byte(avsAddr), '/')
	return tmp
}
