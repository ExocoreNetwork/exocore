package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"

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
	prefixOperatorInfo = iota + 1
	prefixRestakerDelegationInfo
	prefixDelegationUsedSalt
	prefixOperatorApprovedInfo

	prefixUndelegationInfo

	prefixStakerUndelegationInfo

	prefixWaitCompleteUndelegations
)

var (
	// KeyPrefixOperatorInfo key-value: operatorAddr->operatorInfo
	KeyPrefixOperatorInfo = []byte{prefixOperatorInfo}
	// KeyPrefixRestakerDelegationInfo reStakerId = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixRestakerDelegationInfo
	// key-value:
	// reStakerId +'/'+assetID -> totalDelegationAmount
	// reStakerId +'/'+assetID+'/'+operatorAddr -> delegationAmounts

	KeyPrefixRestakerDelegationInfo = []byte{prefixRestakerDelegationInfo}
	// KeyPrefixDelegationUsedSalt key->value: operatorApproveAddr->map[salt]{}
	KeyPrefixDelegationUsedSalt = []byte{prefixDelegationUsedSalt}
	// KeyPrefixOperatorApprovedInfo key-value: operatorApproveAddr->map[reStakerId]{}
	KeyPrefixOperatorApprovedInfo = []byte{prefixOperatorApprovedInfo}

	// KeyPrefixUndelegationInfo singleRecordKey = lzNonce+'/'+txHash+'/'+operatorAddr
	// singleRecordKey -> UndelegateReqRecord
	KeyPrefixUndelegationInfo = []byte{prefixUndelegationInfo}
	// KeyPrefixStakerUndelegationInfo reStakerId+'/'+assetID+'/'+lzNonce -> singleRecordKey
	KeyPrefixStakerUndelegationInfo = []byte{prefixStakerUndelegationInfo}
	// KeyPrefixWaitCompleteUndelegations completeHeight +'/'+lzNonce -> singleRecordKey
	KeyPrefixWaitCompleteUndelegations = []byte{prefixWaitCompleteUndelegations}
)

func GetDelegationStateKey(stakerID, assetID, operatorAddr string) []byte {
	return []byte(strings.Join([]string{stakerID, assetID, operatorAddr}, "/"))
}

func GetDelegationStateIteratorPrefix(stakerID, assetID string) []byte {
	tmp := []byte(strings.Join([]string{stakerID, assetID}, "/"))
	tmp = append(tmp, '/')
	return tmp
}

func ParseStakerAssetIDAndOperatorAddrFromKey(key []byte) (keys *SingleDelegationInfoReq, err error) {
	stringList := strings.Split(string(key), "/")
	if len(stringList) != 3 {
		return nil, errorsmod.Wrap(ErrParseDelegationKey, fmt.Sprintf("the stringList is:%v", stringList))
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
