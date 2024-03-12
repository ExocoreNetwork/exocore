package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// constants
const (
	// ModuleName module name
	ModuleName = "assets"

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
	prefixAppChainInfo
	prefixRestakingAssetInfo
	prefixRestakerAssetInfo
	prefixOperatorAssetInfo
	prefixOperatorOptedInMiddlewareAssetInfo

	prefixRestakerExocoreAddr

	prefixRestakerExocoreAddrReverse

	// prefixReStakingAssetList
	// prefixReStakerAssetList
	// prefixOperatorAssetList

	// add for dogfood
	prefixOperatorSnapshot
	prefixOperatorLastSnapshotHeight
	prefixParams
)

// KVStore key prefixes
var (
	/*
		exoCore stored info:

		//stored info in assets module
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

	KeyPrefixAppChainInfo = []byte{prefixAppChainInfo}

	// KeyPrefixReStakingAssetInfo AssetID = AssetAddr+'_'+chainIndex
	// KeyPrefixReStakingAssetInfo key->value: AssetID->ReStakingAssetInfo
	KeyPrefixReStakingAssetInfo = []byte{prefixRestakingAssetInfo}

	// KeyPrefixReStakerAssetInfos restakerID = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixReStakerAssetInfos key->value: restakerID+'_'+AssetID->ReStakerAssetInfo
	// or restakerID->mapping(AssetID->ReStakerAssetInfo)?
	KeyPrefixReStakerAssetInfos = []byte{prefixRestakerAssetInfo}

	// KeyPrefixOperatorAssetInfos key->value: operatorAddr+'_'+AssetID->OperatorAssetInfo
	// or operatorAddr->mapping(AssetID->OperatorAssetInfo) ?
	KeyPrefixOperatorAssetInfos = []byte{prefixOperatorAssetInfo}

	// KeyPrefixOperatorOptedInMiddleWareAssetInfos key->value:
	// operatorAddr+'_'+AssetID->mapping(middleWareAddr->struct{})
	// or operatorAddr->mapping(AssetID->mapping(middleWareAddr->struct{})) ?
	KeyPrefixOperatorOptedInMiddleWareAssetInfos = []byte{
		prefixOperatorOptedInMiddlewareAssetInfo,
	}

	// KeyPrefixReStakerExoCoreAddr restakerID = clientChainAddr+'_'+ExoCoreChainIndex
	// KeyPrefixReStakerExoCoreAddr key-value: restakerID->exoCoreAddr
	KeyPrefixReStakerExoCoreAddr = []byte{prefixRestakerExocoreAddr}
	// KeyPrefixReStakerExoCoreAddrReverse k->v: exocoreAddress ->
	// map[clientChainIndex]clientChainAddress
	// used to retrieve all user assets based on their exoCore address
	KeyPrefixReStakerExoCoreAddrReverse = []byte{prefixRestakerExocoreAddrReverse}

	// KeyPrefixParams This is a key prefix for module parameter
	KeyPrefixParams = []byte{prefixParams}
	ParamsKey       = []byte("Params")
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

// add for dogfood
func OperatorSnapshotKey(operatorAddr sdk.AccAddress, height uint64) []byte {
	base := []byte{prefixOperatorSnapshot}
	base = append(base, operatorAddr.Bytes()...)
	base = append(base, sdk.Uint64ToBigEndian(height)...)
	return base
}

func OperatorLastSnapshotHeightKey(operatorAddr sdk.AccAddress) []byte {
	base := []byte{prefixOperatorLastSnapshotHeight}
	base = append(base, operatorAddr.Bytes()...)
	return base
}
