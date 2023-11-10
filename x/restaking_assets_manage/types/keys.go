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
	// ModuleName module name
	ModuleName = "restaking_assets_manage"

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

// prefix bytes for the reStaking assets manage store
const (
	prefixClientChainInfo = iota + 1
	prefixReStakingAssetInfo
	prefixReStakerAssetInfos
	prefixOperatorAssetInfos
	prefixOperatorOptedInMiddleWareAssetInfos

	prefixReStakerExocoreAddr

	prefixReStakerExocoreAddrReverse

	prefixReStakingAssetList
	prefixReStakerAssetList
	prefixOperatorAssetList
)

// KVStore key prefixes
var (
	/*
		exoCore stored info:

		//stored info in restaking_assets_manage module
		//used to record supported client chain and reStaking token info
		chainIndex->ChainInfo
		tokenIndex->tokenInfo
		chainList ?
		tokenList ?

		//record restaker reStaking info
		restaker->mapping(tokenIndex->amount)
		restaker->ReStakingTokenList ?
		restakerList?

		//record operator reStaking info
		operator->mapping(tokenIndex->amount)
		operator->ReStakingTokenList ?
		operator->mapping(tokenIndex->middleWareAddress) ?


		//stored info in delegation module
		//record the operator info which restaker delegate to
		restaker->mapping(operator->mapping(tokenIndex->amount))
		restaker->operatorList
		operator->operatorInfo

		//stored info in middleWare module
		middleWareAddr->middleWareInfo
		middleWareAddr->OptedInOperatorInfo
	*/
	// KeyPrefixClientChainInfo key->value: chainIndex->ClientChainInfo
	KeyPrefixClientChainInfo = []byte{prefixClientChainInfo}

	// KeyPrefixReStakingAssetInfo AssetId = AssetAddr+'_'+chainIndex
	// KeyPrefixReStakingAssetInfo key->value: AssetId->ReStakingAssetInfo
	KeyPrefixReStakingAssetInfo = []byte{prefixReStakingAssetInfo}

	// KeyPrefixReStakerAssetInfos reStakerId = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixReStakerAssetInfos key->value: reStakerId+'_'+AssetId->ReStakerSingleAssetInfo
	// or reStakerId->mapping(AssetId->ReStakerSingleAssetInfo)?
	KeyPrefixReStakerAssetInfos = []byte{prefixReStakerAssetInfos}

	// KeyPrefixOperatorAssetInfos key->value: operatorAddr+'_'+AssetId->OperatorSingleAssetInfo
	// or operatorAddr->mapping(AssetId->OperatorSingleAssetInfo) ?
	KeyPrefixOperatorAssetInfos = []byte{prefixOperatorAssetInfos}

	// KeyPrefixOperatorOptedInMiddleWareAssetInfos key->value: operatorAddr+'_'+AssetId->mapping(middleWareAddr->struct{})
	//or operatorAddr->mapping(AssetId->mapping(middleWareAddr->struct{})) ?
	KeyPrefixOperatorOptedInMiddleWareAssetInfos = []byte{prefixOperatorOptedInMiddleWareAssetInfos}

	// KeyPrefixReStakerExoCoreAddr reStakerId = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixReStakerExoCoreAddr key-value: reStakerId->exoCoreAddr
	KeyPrefixReStakerExoCoreAddr = []byte{prefixReStakerExocoreAddr}
	//KeyPrefixReStakerExoCoreAddrReverse k->v: exocoreAddress -> map[clientChainIndex]clientChainAddress
	// used to retrieve all user assets based on their exoCore address
	KeyPrefixReStakerExoCoreAddrReverse = []byte{prefixReStakerExocoreAddrReverse}
)

func GetAssetStateKey(stakerId, assetId string) []byte {
	return []byte(strings.Join([]string{stakerId, assetId}, "/"))
}

func ParseStakerAndAssetIdFromKey(key []byte) (stakerId string, assetId string, err error) {
	stringList := strings.Split(string(key), "/")
	if len(stringList) != 2 {
		return "", "", errorsmod.Wrap(ErrParseAssetsStateKey, fmt.Sprintf("the stringList is:%v", stringList))
	}
	return stringList[0], stringList[1], nil
}
