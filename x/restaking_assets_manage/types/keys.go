// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// constants
const (
	// module name
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
	// KeyPrefixReStakerAssetInfos key->value: reStakerId+'_'+AssetId->amount
	// or reStakerId->mapping(AssetId->amount)
	// or reStakerAddr+'_'+tokenIndex->amount ?
	KeyPrefixReStakerAssetInfos = []byte{prefixReStakerAssetInfos}

	// KeyPrefixOperatorAssetInfos key->value: operatorAddr+'_'+AssetId->amount
	// or operatorAddr->mapping(AssetId->amount) ?
	KeyPrefixOperatorAssetInfos = []byte{prefixOperatorAssetInfos}

	// KeyPrefixOperatorOptedInMiddleWareAssetInfos key->value: operatorAddr+'_'+AssetId->mapping(middleWareAddr->struct{})
	//or operatorAddr->mapping(AssetId->mapping(middleWareAddr->struct{})) ?
	KeyPrefixOperatorOptedInMiddleWareAssetInfos = []byte{prefixOperatorOptedInMiddleWareAssetInfos}
)
