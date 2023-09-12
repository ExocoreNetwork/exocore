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
	prefixReStakingAssetList

	prefixReStakerAssetInfos
	prefixReStakerAssetList
	prefixReStakerExocoreAddr
	prefixOperatorAssetInfos
	prefixOperatorAssetList
	prefixOperatorOptedInMiddleWareAssetInfos
)

// KVStore key prefixes
var (
	// KeyPrefixClientChainInfo key->value: chainIndex->ClientChainInfo
	KeyPrefixClientChainInfo = []byte{prefixClientChainInfo}

	// KeyPrefixReStakingAssetInfo AssetId = AssetAddr+'_'+chainIndex
	// KeyPrefixReStakingAssetInfo key->value: AssetId->ReStakingAssetInfo
	KeyPrefixReStakingAssetInfo = []byte{prefixReStakingAssetInfo}

	// KeyPrefixReStakingAssetList list: ReStakingAssetList
	KeyPrefixReStakingAssetList = []byte{prefixReStakingAssetList}

	// KeyPrefixReStakerExoCoreAddr key-value: clientChainAddr+'_'+ExoCoreChainIndex : exoCoreAddr
	KeyPrefixReStakerExoCoreAddr = []byte{prefixReStakerExocoreAddr}

	// KeyPrefixReStakerAssetInfos key->value: clientChainAddr+'_'+ExoCoreChainIndex+'_'+AssetId->amount
	// or reStakerAddr+'_'+tokenIndex->amount ?
	KeyPrefixReStakerAssetInfos = []byte{prefixReStakerAssetInfos}

	// KeyPrefixReStakerAssetList key->value: reStakerAddr->ReStakingAssetList
	KeyPrefixReStakerAssetList = []byte{prefixReStakerAssetList}

	// KeyPrefixOperatorAssetInfos key->value: operatorAddr+'_'+AssetId->amount
	KeyPrefixOperatorAssetInfos = []byte{prefixOperatorAssetInfos}

	// KeyPrefixOperatorAssetList key->value: operatorAddr ->ReStakingAssetList
	KeyPrefixOperatorAssetList = []byte{prefixOperatorAssetList}

	// KeyPrefixOperatorOptedInMiddleWareAssetInfos key->value: operatorAddr+'_'+AssetId->middleWareAddr
	KeyPrefixOperatorOptedInMiddleWareAssetInfos = []byte{prefixOperatorOptedInMiddleWareAssetInfos}
)
