package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SharedParameter is a shared parameter used to record and update the related
// USD share of staker and operators when the prices of assets change
type SharedParameter struct {
	// priceChangeAssets is the price change information of assets,
	// which is gotten by the expected Oracle interface.
	priceChangeAssets map[string]*operatortypes.PriceChange
	// assetsDecimal is a map to record the decimals of the related assets
	// It will be used when calculate the USD share of the assets.
	assetsDecimal map[string]uint32
	// optedInAssetsInfo : assetID->operator->Avs
	// For staker and operator, only the USD share of opted-in assets needs to be updated
	// when the prices of assets change. But in the delegation and assets module, the opted-in
	// information of assets haven't been stored, so we need this map as a filter when iterate
	// the assets state of delegation and operator
	// It will be set when calling `IterateUpdateAssetState`, because the information of opted-in assets
	// has been stored in types.KeyPrefixOperatorAVSSingleAssetState
	optedInAssetsInfo map[string]map[string]string
	// stakerShare records the latest share for staker and operator after updating
	stakerShare map[string]sdkmath.LegacyDec
}

func UpdateShareOfStakerAndOperator(sharedParam *SharedParameter, assetID, stakerID, operatorAddr string, assetAmount sdkmath.Int) {
	priceChange := sharedParam.priceChangeAssets[assetID]
	assetDecimal := sharedParam.assetsDecimal[assetID]
	if avsAddr, ok := sharedParam.optedInAssetsInfo[assetID][operatorAddr]; ok {
		newAssetUSDValue := CalculateShare(assetAmount, priceChange.NewPrice, assetDecimal, priceChange.Decimal)
		key := string(types.GetJoinedStoreKey(avsAddr, stakerID, operatorAddr))
		AddShareInMap(sharedParam.stakerShare, key, newAssetUSDValue)
	}
}

// PriceChangeHandle update the assets' share when their prices change
func (k *Keeper) PriceChangeHandle(ctx sdk.Context) error {
	priceChangeAssets, err := k.oracleKeeper.GetPriceChangeAssets(ctx)
	if err != nil {
		return err
	}
	if len(priceChangeAssets) == 0 {
		return nil
	}
	shareChangeForAvsOperator := make(map[string]sdkmath.LegacyDec, 0)
	optedInAssetsInfo := make(map[string]map[string]string, 0)
	assetsDecimal := make(map[string]uint32)
	for assetID, priceChange := range priceChangeAssets {
		// get the decimal of asset
		assetInfo, err := k.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
		if err != nil {
			return err
		}
		assetsDecimal[assetID] = assetInfo.AssetBasicInfo.Decimals
		if _, ok := optedInAssetsInfo[assetID]; !ok {
			optedInAssetsInfo[assetID] = make(map[string]string, 0)
		}
		// UpdateStateForAsset
		f := func(assetID string, keys []string, state *operatortypes.OptedInAssetState) error {
			newAssetUSDValue := CalculateShare(state.Amount, priceChange.NewPrice, assetInfo.AssetBasicInfo.Decimals, priceChange.Decimal)
			changeValue := newAssetUSDValue.Sub(state.Value)
			state.Value = newAssetUSDValue

			avsAddr := keys[1]
			avsOperator := string(types.GetJoinedStoreKey(keys[1], keys[2]))
			AddShareInMap(shareChangeForAvsOperator, avsAddr, changeValue)
			AddShareInMap(shareChangeForAvsOperator, avsOperator, changeValue)
			optedInAssetsInfo[assetID][keys[2]] = avsAddr
			return nil
		}
		err = k.IterateUpdateAssetState(ctx, assetID, f)
		if err != nil {
			return err
		}
	}
	// BatchUpdateShareForAVSAndOperator
	err = k.BatchUpdateShareForAVSAndOperator(ctx, shareChangeForAvsOperator)
	if err != nil {
		return err
	}

	// update the USD share for staker and operator
	sharedParameter := &SharedParameter{
		priceChangeAssets: priceChangeAssets,
		assetsDecimal:     assetsDecimal,
		optedInAssetsInfo: optedInAssetsInfo,
		stakerShare:       make(map[string]sdkmath.LegacyDec, 0),
	}
	stakerShareHandleFunc := func(stakerID, assetID, operatorAddr string, state *delegationtype.DelegationAmounts) error {
		UpdateShareOfStakerAndOperator(sharedParameter, assetID, stakerID, operatorAddr, state.UndelegatableAmount)
		return nil
	}
	err = k.delegationKeeper.IterateDelegationState(ctx, stakerShareHandleFunc)
	if err != nil {
		return err
	}

	operatorShareHandleFunc := func(operatorAddr, assetID string, state *types.OperatorAssetInfo) error {
		UpdateShareOfStakerAndOperator(sharedParameter, assetID, "", operatorAddr, state.OperatorAmount)
		return nil
	}
	err = k.assetsKeeper.IteratorOperatorAssetState(ctx, operatorShareHandleFunc)
	if err != nil {
		return err
	}
	// BatchSetStakerShare
	err = k.BatchSetStakerShare(ctx, sharedParameter.stakerShare)
	if err != nil {
		return err
	}
	return nil
}

// ClearPreConsensusPK clears the previous consensus public key for all operators
func (k *Keeper) ClearPreConsensusPK(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(
		store,
		[]byte{operatortypes.BytePrefixForOperatorAndChainIDToPrevConsKey},
	)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
	return nil
}

// EndBlock : update the assets' share when their prices change
func (k *Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	// todo: need to consider the calling order
	err := k.PriceChangeHandle(ctx)
	if err != nil {
		panic(err)
	}

	err = k.ClearPreConsensusPK(ctx)
	if err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}
