package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
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
	prefixRestakingAssetInfo
	prefixRestakerAssetInfo
	prefixOperatorAssetInfo
	prefixOperatorOptedInMiddlewareAssetInfo

	prefixRestakerExocoreAddr

	prefixRestakerExocoreAddrReverse

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
	KeyPrefixReStakingAssetInfo = []byte{prefixRestakingAssetInfo}

	// KeyPrefixReStakerAssetInfos reStakerId = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixReStakerAssetInfos key->value: reStakerId+'_'+AssetId->ReStakerSingleAssetInfo
	// or reStakerId->mapping(AssetId->ReStakerSingleAssetInfo)?
	KeyPrefixReStakerAssetInfos = []byte{prefixRestakerAssetInfo}

	// KeyPrefixOperatorAssetInfos key->value: operatorAddr+'_'+AssetId->OperatorSingleAssetInfo
	// or operatorAddr->mapping(AssetId->OperatorSingleAssetInfo) ?
	KeyPrefixOperatorAssetInfos = []byte{prefixOperatorAssetInfo}

	// KeyPrefixOperatorOptedInMiddleWareAssetInfos key->value: operatorAddr+'_'+AssetId->mapping(middleWareAddr->struct{})
	// or operatorAddr->mapping(AssetId->mapping(middleWareAddr->struct{})) ?
	KeyPrefixOperatorOptedInMiddleWareAssetInfos = []byte{prefixOperatorOptedInMiddlewareAssetInfo}

	// KeyPrefixReStakerExoCoreAddr reStakerId = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixReStakerExoCoreAddr key-value: reStakerId->exoCoreAddr
	KeyPrefixReStakerExoCoreAddr = []byte{prefixRestakerExocoreAddr}
	// KeyPrefixReStakerExoCoreAddrReverse k->v: exocoreAddress -> map[clientChainIndex]clientChainAddress
	// used to retrieve all user assets based on their exoCore address
	KeyPrefixReStakerExoCoreAddrReverse = []byte{prefixRestakerExocoreAddrReverse}
)

func GetJoinedStoreKey(keys ...string) []byte {
	return []byte(strings.Join(keys, "/"))
}

func ParseJoinedKey(key []byte) (keys []string, err error) {
	stringList := strings.Split(string(key), "/")
	return stringList, nil
}

func ParseJoinedStoreKey(key []byte, number int) (keys []string, err error) {
	stringList := strings.Split(string(key), "/")
	if len(stringList) != number {
		return nil, errorsmod.Wrap(ErrParseAssetsStateKey, fmt.Sprintf("expected length:%d,actual length:%d,the stringList is:%v", number, len(stringList), stringList))
	}
	return stringList, nil
}
