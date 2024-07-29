package keeper

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/evmos/v14/x/evm/statedb"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetAVSSupportedAssets returns a map of assets supported by the AVS. The avsAddr supplied must be hex.
func (k *Keeper) GetAVSSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error) {
	if !common.IsHexAddress(avsAddr) {
		return nil, errorsmod.Wrap(types.ErrInvalidAddr, fmt.Sprintf("GetAVSSupportedAssets: key is %s", avsAddr))
	}
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("GetAVSSupportedAssets: key is %s", avsAddr))
	}
	assetIDList := avsInfo.Info.AssetIDs
	ret := make(map[string]interface{})

	for _, assetID := range assetIDList {
		asset, err := k.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
		if err != nil {
			return nil, errorsmod.Wrap(err, fmt.Sprintf("[GetAVSSupportedAssets] GetStakingAssetInfo: key is %s", assetID))
		}
		ret[assetID] = asset
	}

	return ret, nil
}

// GetAVSSlashContract returns the address of the contract that will be used to slash the AVS.
// The avsAddr supplied must be hex.
func (k *Keeper) GetAVSSlashContract(ctx sdk.Context, avsAddr string) (string, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return "", errorsmod.Wrap(err, fmt.Sprintf("GetAVSSlashContract: key is %s", avsAddr))
	}

	return avsInfo.Info.SlashAddr, nil
}

// GetAVSMinimumSelfDelegation returns the minimum self-delegation required for the AVS, on a per-operator basis.
// The avsAddr supplied must be hex.
func (k *Keeper) GetAVSMinimumSelfDelegation(ctx sdk.Context, avsAddr string) (sdkmath.LegacyDec, error) {
	avsInfo, err := k.GetAVSInfo(ctx, avsAddr)
	if err != nil {
		return sdkmath.LegacyNewDec(0), errorsmod.Wrap(err, fmt.Sprintf("GetAVSMinimumSelfDelegation: key is %s", avsAddr))
	}

	return sdkmath.LegacyNewDec(int64(avsInfo.Info.MinSelfDelegation)), nil
}

// GetEpochEndAVSs returns a list of hex AVS addresses for AVSs which are scheduled to start at the end of the
// current epoch, or the beginning of the next one. The address format returned is hex.
func (k *Keeper) GetEpochEndAVSs(ctx sdk.Context, epochIdentifier string, endingEpochNumber int64) []string {
	var avsList []string
	k.IterateAVSInfo(ctx, func(_ int64, avsInfo types.AVSInfo) (stop bool) {
		// consider the dogfood AVS as an example. it is scheduled to start at epoch 0.
		// the currentEpoch is 1, so we will return it.
		// consider another AVS which will start at epoch 5. the current epoch is 4.
		// it should be returned here, since the operator module should start tracking this.
		if epochIdentifier == avsInfo.EpochIdentifier && endingEpochNumber >= int64(avsInfo.StartingEpoch)-1 {
			avsList = append(avsList, avsInfo.AvsAddress)
		}
		return false
	})

	return avsList
}

// GetAVSInfoByTaskAddress returns the AVS  which containing this task address
func (k *Keeper) GetAVSInfoByTaskAddress(ctx sdk.Context, taskAddr string) types.AVSInfo {
	avs := types.AVSInfo{}
	k.IterateAVSInfo(ctx, func(_ int64, avsInfo types.AVSInfo) (stop bool) {
		if taskAddr == avsInfo.GetTaskAddr() {
			avs = avsInfo
		}
		return false
	})
	return avs
}

// GetTaskChallengeEpochEndAVSs returns the task list where the current block marks the end of their challenge period.
func (k *Keeper) GetTaskChallengeEpochEndAVSs(ctx sdk.Context, epochIdentifier string, epochNumber int64) []types.TaskInfo {
	var taskList []types.TaskInfo
	k.IterateTaskAVSInfo(ctx, func(_ int64, taskInfo types.TaskInfo) (stop bool) {
		avsInfo := k.GetAVSInfoByTaskAddress(ctx, taskInfo.TaskContractAddress)
		// Determine if the challenge period has passed, the range of the challenge period is the num marked (StartingEpoch) add TaskChallengePeriod
		if epochIdentifier == avsInfo.EpochIdentifier && epochNumber > int64(taskInfo.TaskChallengePeriod)+int64(taskInfo.StartingEpoch) {
			taskList = append(taskList, taskInfo)
		}
		return false
	})
	return taskList
}

// RegisterAVSWithChainID registers an AVS by its chainID.
// It is responsible for generating an AVS address based on the chainID.
// The following bare minimum parameters must be supplied:
// AssetIDs, EpochsUntilUnbonded, EpochIdentifier, MinSelfDelegation and StartingEpoch.
// This will ensure compatibility with all of the related AVS functions, like
// GetEpochEndAVSs, GetAVSSupportedAssets, and GetAVSMinimumSelfDelegation.
func (k Keeper) RegisterAVSWithChainID(
	oCtx sdk.Context, params *types.AVSRegisterOrDeregisterParams,
) (avsAddr string, err error) {
	// guard against errors
	ctx, writeFunc := oCtx.CacheContext()
	// remove the version number and validate
	params.ChainID = types.ChainIDWithoutRevision(params.ChainID)
	if len(params.ChainID) == 0 {
		return "", errorsmod.Wrap(types.ErrNotNull, "RegisterAVSWithChainID: chainID is null")
	}
	avsAddr = types.GenerateAVSAddr(params.ChainID)
	defer func() {
		if err == nil {
			// store the reverse lookup from AVSAddress to ChainID
			// (the forward can be generated on the fly by hashing).
			k.SetAVSAddrToChainID(ctx, avsAddr, params.ChainID)
			// write the cache
			writeFunc()
			// TODO: do events need to be handled separately? currently no events emitted so not urgent.
		}
	}()
	// Mark the account as occupied by a contract, so that any transactions that originate
	// from it are rejected in `state_transition.go`. This protects us against the very
	// rare case of address collision.
	if err := k.evmKeeper.SetAccount(
		ctx, common.HexToAddress(avsAddr),
		statedb.Account{
			Balance:  big.NewInt(0),
			CodeHash: types.ChainIDCodeHash[:],
			Nonce:    1,
		},
	); err != nil {
		return "", err
	}
	// SetAVSInfo expects sdk.AccAddress not HexAddress
	params.AvsAddress = avsAddr
	params.Action = RegisterAction

	if err := k.AVSInfoUpdate(ctx, params); err != nil {
		return "", err
	}
	return avsAddr, nil
}

// SetAVSAddressToChainID stores a lookup from the generated AVS address to the chainID.
func (k Keeper) SetAVSAddrToChainID(ctx sdk.Context, avsAddr, chainID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSAddressToChainID)
	store.Set([]byte(avsAddr), []byte(chainID))
}

// GetAVSAddrByChainID returns the AVS address for a given chainID. It is a stateless
// function, even though it is implemented as a method on the keeper.
func (k Keeper) GetAVSAddrByChainID(_ sdk.Context, chainID string) string {
	return types.GenerateAVSAddr(
		types.ChainIDWithoutRevision(chainID),
	)
}

// GetChainIDByAVSAddr returns the chainID for a given AVS address. It is a stateful
// function since it only returns the chainID if an AVS with the chainID was previously
// registered.
func (k Keeper) GetChainIDByAVSAddr(ctx sdk.Context, avsAddr string) (string, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSAddressToChainID)
	bz := store.Get([]byte(avsAddr))
	if bz == nil {
		return "", false
	}
	return string(bz), true
}
