package types

import (
	"strings"

	"github.com/ExocoreNetwork/exocore/x/assets/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// constants
const (
	// ModuleName module name
	ModuleName = "delegation"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

const (
	prefixRestakerDelegationInfo = iota + 1
	prefixDelegationUsedSalt
	prefixOperatorApprovedInfo

	prefixUndelegationInfo

	prefixStakerUndelegationInfo

	prefixWaitCompleteUndelegations

	// add for dogfood
	prefixUndelegationOnHold
)

var (
	// KeyPrefixRestakerDelegationInfo restakerID = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixRestakerDelegationInfo
	// key-value:
	// restakerID +'/'+assetID -> totalDelegationAmount
	// restakerID +'/'+assetID+'/'+operatorAddr -> delegationAmounts
	KeyPrefixRestakerDelegationInfo = []byte{prefixRestakerDelegationInfo}
	// KeyPrefixDelegationUsedSalt key->value: operatorApproveAddr->map[salt]{}
	KeyPrefixDelegationUsedSalt = []byte{prefixDelegationUsedSalt}
	// KeyPrefixOperatorApprovedInfo key-value: operatorApproveAddr->map[restakerID]{}
	KeyPrefixOperatorApprovedInfo = []byte{prefixOperatorApprovedInfo}

	// KeyPrefixUndelegationInfo singleRecordKey = lzNonce+'/'+txHash+'/'+operatorAddr
	// singleRecordKey -> UndelegateReqRecord
	KeyPrefixUndelegationInfo = []byte{prefixUndelegationInfo}
	// KeyPrefixStakerUndelegationInfo restakerID+'/'+assetID+'/'+lzNonce -> singleRecordKey
	KeyPrefixStakerUndelegationInfo = []byte{prefixStakerUndelegationInfo}
	// KeyPrefixWaitCompleteUndelegations completeHeight +'/'+lzNonce -> singleRecordKey
	KeyPrefixWaitCompleteUndelegations = []byte{prefixWaitCompleteUndelegations}
)

func GetDelegationStateIteratorPrefix(stakerID, assetID string) []byte {
	tmp := []byte(strings.Join([]string{stakerID, assetID}, "/"))
	tmp = append(tmp, '/')
	return tmp
}

func ParseStakerAssetIDAndOperatorAddrFromKey(key []byte) (keys *SingleDelegationInfoReq, err error) {
	stringList, err := types.ParseJoinedStoreKey(key, 3)
	if err != nil {
		return nil, err
	}
	return &SingleDelegationInfoReq{StakerID: stringList[0], AssetID: stringList[1], OperatorAddr: stringList[2]}, nil
}

func GetUndelegationRecordKey(lzNonce uint64, txHash string, operatorAddr string) []byte {
	return []byte(strings.Join([]string{hexutil.EncodeUint64(lzNonce), txHash, operatorAddr}, "/"))
}

func GetStakerUndelegationRecordKey(stakerID, assetID string, lzNonce uint64) []byte {
	return []byte(strings.Join([]string{stakerID, assetID, hexutil.EncodeUint64(lzNonce)}, "/"))
}

func GetWaitCompleteRecordKey(height, lzNonce uint64) []byte {
	return []byte(strings.Join([]string{hexutil.EncodeUint64(height), hexutil.EncodeUint64(lzNonce)}, "/"))
}

// GetUndelegationOnHoldKey add for dogfood
func GetUndelegationOnHoldKey(recordKey []byte) []byte {
	return append([]byte{prefixUndelegationOnHold}, recordKey...)
}
