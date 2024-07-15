package keeper

import (
	"fmt"
	"hash"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	"github.com/ExocoreNetwork/exocore/x/avs/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetAVSSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("GetAVSSupportedAssets: key is %s", avsAddr))
	}
	assetIDList := avsInfo.Info.AssetId
	ret := make(map[string]interface{})

	for _, assetID := range assetIDList {
		asset, err := k.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
		if err != nil {
			return nil, errorsmod.Wrap(err, fmt.Sprintf("GetStakingAssetInfo: key is %s", assetID))
		}
		ret[assetID] = asset
	}

	return ret, nil
}

func (k *Keeper) GetAVSSlashContract(ctx sdk.Context, avsAddr string) (string, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return "", errorsmod.Wrap(err, fmt.Sprintf("GetAVSSlashContract: key is %s", avsAddr))
	}

	return avsInfo.Info.SlashAddr, nil
}

// GetAVSMinimumSelfDelegation returns the USD value of minimum self delegation, which
// is set for operator
func (k *Keeper) GetAVSMinimumSelfDelegation(ctx sdk.Context, avsAddr string) (sdkmath.LegacyDec, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return sdkmath.LegacyNewDec(0), errorsmod.Wrap(err, fmt.Sprintf("GetAVSMinimumSelfDelegation: key is %s", avsAddr))
	}

	return sdkmath.LegacyNewDec(avsInfo.Info.MinSelfDelegation.Int64()), nil
}

// GetEpochEndAVSs returns the AVS list where the current block marks the end of their epoch.
func (k *Keeper) GetEpochEndAVSs(ctx sdk.Context, epochIdentifier string, epochNumber int64) ([]string, error) {
	var avsList []types.AVSInfo
	k.IterateAVSInfo(ctx, func(_ int64, avsInfo types.AVSInfo) (stop bool) {
		if epochIdentifier == avsInfo.EpochIdentifier && epochNumber > avsInfo.StartingEpoch {
			avsList = append(avsList, avsInfo)
		}
		return false
	})

	if len(avsList) == 0 {
		return []string{}, nil
	}

	avsAddrList := make([]string, len(avsList))
	for i := range avsList {
		avsAddrList[i] = avsList[i].AvsAddress
	}

	return avsAddrList, nil
}

func (k *Keeper) GetAVSAddrByChainID(ctx sdk.Context, chainID string) (string, error) {
	chainID = ProcessingStr(chainID)
	if len(chainID) == 0 {
		return "", errorsmod.Wrap(types.ErrNotNull, "SetAVSAddrByChainID: chainID is null")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfoByChainID)
	if !store.Has([]byte(chainID)) {
		return "", errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSAddrByChainID: key is %s", chainID))
	}
	avsAddr := store.Get([]byte(chainID))

	return string(avsAddr), nil
}

// SetAVSAddrByChainID creates an avs address given the chainID
func (k Keeper) SetAVSAddrByChainID(ctx sdk.Context, chainID string) (err error) {
	chainID = ProcessingStr(chainID)
	if len(chainID) == 0 {
		return errorsmod.Wrap(err, "SetAVSAddrByChainID: chainID is null")
	}
	avsAddr := common.BytesToAddress(Keccak256([]byte(chainID))).String()

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfoByChainID)

	store.Set([]byte(chainID), []byte(avsAddr))
	return nil
}

// KeccakState wraps sha3.state. In addition to the usual hash methods, it also supports
// Read to get a variable amount of data from the hash state. Read is faster than Sum
// because it doesn't copy the internal state, but also modifies the internal state.
type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

// NewKeccakState creates a new KeccakState
func NewKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := NewKeccakState()
	for _, b := range data {
		d.Write(b)
	}
	_, err := d.Read(b)
	if err != nil {
		return nil
	}
	return b
}

func ProcessingStr(str string) string {
	index := strings.Index(str, "-")
	if index != -1 {
		return str[:index]
	}
	return ""
}
