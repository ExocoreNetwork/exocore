package keeper

import (
	"fmt"
	"math/big"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
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

func (k *Keeper) GetAVSAddrByChainID(ctx sdk.Context, chainID string) string {
	avsAddr := ""
	k.IterateAVSInfo(ctx, func(_ int64, avsInfo types.AVSInfo) (stop bool) {
		if chainID == avsInfo.GetChainId() {
			avsAddr = avsInfo.AvsAddress
		}
		return false
	})
	return avsAddr
}

func (k *Keeper) GetChainIDByAVSAddr(ctx sdk.Context, avsAddr string) string {
	chainID := ""
	k.IterateAVSInfo(ctx, func(_ int64, avsInfo types.AVSInfo) (stop bool) {
		if avsAddr == avsInfo.GetAvsAddress() {
			chainID = avsInfo.ChainId
		}
		return false
	})
	return chainID
}

// RegisterAVSForDogFood During this registration, the caller should be able to provide
// AssetIDs, EpochsUntilUnbonded, EpochIdentifier, MinSelfDelegation and StartingEpoch.
// Other values should be empty / blank, and we should call AVsInfoUpdate.
// This will ensure compatibility with all of the related AVS functions, like
// GetEpochEndAVSs, GetAVSSupportedAssets, and GetAVSMinimumSelfDelegation
func (k Keeper) RegisterAVSForDogFood(ctx sdk.Context, params *AVSRegisterOrDeregisterParams) (err error) {
	chainID := ChainIDWithoutRevision(params.ChainID)
	if len(chainID) == 0 {
		return errorsmod.Wrap(types.ErrNotNull, "RegisterAVSWithChainID: chainID is null")
	}
	params.ChainID = chainID
	avsAddr := common.BytesToAddress(crypto.Keccak256([]byte(chainID))).String()

	err = k.evmKeeper.SetAccount(ctx, common.HexToAddress(avsAddr), statedb.Account{
		Balance:  big.NewInt(0),
		CodeHash: crypto.Keccak256Hash([]byte(types.ChainID)).Bytes(),
		Nonce:    1,
	})
	if err != nil {
		return err
	}
	params.AvsAddress = avsAddr
	params.Action = RegisterAction

	return k.AVSInfoUpdate(ctx, params)
}

func ChainIDWithoutRevision(chainID string) string {
	if !ibcclienttypes.IsRevisionFormat(chainID) {
		return chainID
	}
	splitStr := strings.Split(chainID, "-")
	return splitStr[0]
}
