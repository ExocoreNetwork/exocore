// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"strings"
)

// constants
const (
	// module name
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
)

var (
	// KeyPrefixOperatorInfo key-value: operatorAddr->operatorInfo
	KeyPrefixOperatorInfo = []byte{prefixOperatorInfo}
	// KeyPrefixRestakerDelegationInfo reStakerId = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixRestakerDelegationInfo key-value: reStakerId +'/'+assetId+'/'+operatorAddr -> delegationAmount
	KeyPrefixRestakerDelegationInfo = []byte{prefixRestakerDelegationInfo}
	// KeyPrefixDelegationUsedSalt key->value: operatorApproveAddr->map[salt]{}
	KeyPrefixDelegationUsedSalt = []byte{prefixDelegationUsedSalt}
	// KeyPrefixOperatorApprovedInfo key-value: operatorApproveAddr->map[reStakerId]{}
	KeyPrefixOperatorApprovedInfo = []byte{prefixOperatorApprovedInfo}

	//KeyPrefixUnDelegationInfo key-value: ReStakerId+'_'+nonce -> UnDelegateReqRecord
	KeyPrefixUnDelegationInfo = []byte{prefixUndelegationInfo}
)

func GetDelegationStateKey(stakerId, assetId, operatorAddr string) []byte {
	return []byte(strings.Join([]string{stakerId, assetId, operatorAddr}, "/"))
}

func ParseStakerAssetIdAndOperatorAddrFromKey(key []byte) (keys *SingleDelegationInfoReq, err error) {
	stringList := strings.Split(string(key), "/")
	if len(stringList) != 3 {
		return nil, errorsmod.Wrap(ErrParseDelegationKey, fmt.Sprintf("the stringList is:%v", stringList))
	}
	return &SingleDelegationInfoReq{StakerId: stringList[0], AssetId: stringList[1], OperatorAddr: stringList[2]}, nil
}
