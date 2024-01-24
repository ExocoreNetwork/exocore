package types

import (
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

	UsdValueDefaultDecimal = uint8(8)

	SlashVetoDuration = int64(1000)
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
	KeyPrefixSlashAssetsState = []byte{prefixSlashAssetsState}
)
