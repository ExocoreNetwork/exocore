package keeper

import (
	"fmt"
	"math/big"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/evmos/v14/x/evm/statedb"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetAVSSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("GetAVSSupportedAssets: key is %s", avsAddr))
	}
	assetIDList := avsInfo.Info.AssetIDs
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

	return sdkmath.LegacyNewDec(int64(avsInfo.Info.MinSelfDelegation)), nil
}

// GetEpochEndAVSs returns the AVS list where the current block marks the end of their epoch.
func (k *Keeper) GetEpochEndAVSs(ctx sdk.Context, epochIdentifier string, epochNumber int64) []string {
	var avsList []string
	k.IterateAVSInfo(ctx, func(_ int64, avsInfo types.AVSInfo) (stop bool) {
		if epochIdentifier == avsInfo.EpochIdentifier && epochNumber > int64(avsInfo.StartingEpoch) {
			avsList = append(avsList, avsInfo.AvsAddress)
		}
		return false
	})

	return avsList
}

// GetTaskChallengeEpochEndAVSs returns the AVS list where the current block marks the end of their epoch.
func (k *Keeper) GetTaskChallengeEpochEndAVSs(ctx sdk.Context, epochIdentifier string, epochNumber int64) []string {
	var avsList []string
	k.IterateAVSInfo(ctx, func(_ int64, avsInfo types.AVSInfo) (stop bool) {
		if epochIdentifier == avsInfo.EpochIdentifier && epochNumber > int64(avsInfo.StartingEpoch) {
			avsList = append(avsList, avsInfo.AvsAddress)
		}
		return false
	})
	return avsList
}

func (k *Keeper) GetAVSAddrByChainID(ctx sdk.Context, chainID string) (string, error) {
	chainID = ChainIDWithoutRevision(chainID)
	if len(chainID) == 0 {
		return "", errorsmod.Wrap(types.ErrNotNull, "RegisterAVSWithChainID: chainID is null")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfoByChainID)
	if !store.Has([]byte(chainID)) {
		return "", errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSAddrByChainID: key is %s", chainID))
	}
	avsAddr := store.Get([]byte(chainID))

	return string(avsAddr), nil
}

func (k *Keeper) GetChainIDByAVSAddr(ctx sdk.Context, avsAddr string) (string, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfoByChainID)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		chainID := string(iterator.Key())
		if string(iterator.Value()) == avsAddr {
			return chainID, nil
		}
	}

	return "", errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetChainIDByAVSAddr: key is %s", avsAddr))
}

// RegisterAVSWithChainID creates an avs address given the chainID
func (k Keeper) RegisterAVSWithChainID(ctx sdk.Context, chainID string) (err error) {
	chainID = ChainIDWithoutRevision(chainID)
	if len(chainID) == 0 {
		return errorsmod.Wrap(types.ErrNotNull, "RegisterAVSWithChainID: chainID is null")
	}
	avsAddr := common.BytesToAddress(crypto.Keccak256([]byte(chainID))).String()

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfoByChainID)

	err = k.evmKeeper.SetAccount(ctx, common.HexToAddress(avsAddr), statedb.Account{
		Balance:  big.NewInt(0),
		CodeHash: crypto.Keccak256Hash([]byte(types.ChainID)).Bytes(),
		Nonce:    1,
	})
	if err != nil {
		return err
	}
	store.Set([]byte(chainID), []byte(avsAddr))
	return nil
}

func ChainIDWithoutRevision(chainID string) string {
	if !ibcclienttypes.IsRevisionFormat(chainID) {
		return chainID
	}
	splitStr := strings.Split(chainID, "-")
	return splitStr[0]
}
