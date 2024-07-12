package keeper

import (
	"fmt"

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
