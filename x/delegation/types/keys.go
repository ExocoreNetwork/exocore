// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strings"
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
	//key-value:
	//reStakerId +'/'+assetId -> totalDelegationAmount
	//reStakerId +'/'+assetId+'/'+operatorAddr -> delegationAmounts

	KeyPrefixRestakerDelegationInfo = []byte{prefixRestakerDelegationInfo}
	// KeyPrefixDelegationUsedSalt key->value: operatorApproveAddr->map[salt]{}
	KeyPrefixDelegationUsedSalt = []byte{prefixDelegationUsedSalt}
	// KeyPrefixOperatorApprovedInfo key-value: operatorApproveAddr->map[reStakerId]{}
	KeyPrefixOperatorApprovedInfo = []byte{prefixOperatorApprovedInfo}

	//KeyPrefixUndelegationInfo singleRecordKey = lzNonce+'/'+txHash+'/'+operatorAddr
	// singleRecordKey -> UndelegateReqRecord
	KeyPrefixUndelegationInfo = []byte{prefixUndelegationInfo}
	//KeyPrefixStakerUndelegationInfo reStakerId+'/'+assetId+'/'+lzNonce -> singleRecordKey
	KeyPrefixStakerUndelegationInfo = []byte{prefixStakerUndelegationInfo}
	//KeyPrefixWaitCompleteUndelegations completeHeight +'/'+lzNonce -> singleRecordKey
	KeyPrefixWaitCompleteUndelegations = []byte{prefixWaitCompleteUndelegations}
)

func GetDelegationStateKey(stakerId, assetId, operatorAddr string) []byte {
	return []byte(strings.Join([]string{stakerId, assetId, operatorAddr}, "/"))
}

func GetDelegationStateIteratorPrefix(stakerId, assetId string) []byte {
	tmp := []byte(strings.Join([]string{stakerId, assetId}, "/"))
	tmp = append(tmp, '/')
	return tmp
}

func ParseStakerAssetIdAndOperatorAddrFromKey(key []byte) (keys *SingleDelegationInfoReq, err error) {
	stringList := strings.Split(string(key), "/")
	if len(stringList) != 3 {
		return nil, errorsmod.Wrap(ErrParseDelegationKey, fmt.Sprintf("the stringList is:%v", stringList))
	}
	return &SingleDelegationInfoReq{StakerId: stringList[0], AssetId: stringList[1], OperatorAddr: stringList[2]}, nil
}

func GetUndelegationRecordKey(lzNonce uint64, txHash string, operatorAddr string) []byte {
	return []byte(strings.Join([]string{hexutil.EncodeUint64(lzNonce), txHash, operatorAddr}, "/"))
}

func GetStakerUndelegationRecordKey(stakerId, assetId string, lzNonce uint64) []byte {
	return []byte(strings.Join([]string{stakerId, assetId, hexutil.EncodeUint64(lzNonce)}, "/"))
}

func GetWaitCompleteRecordKey(height, lzNonce uint64) []byte {
	return []byte(strings.Join([]string{hexutil.EncodeUint64(height), hexutil.EncodeUint64(lzNonce)}, "/"))
}
