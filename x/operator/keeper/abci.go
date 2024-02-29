package keeper

import (
	sdkmath "cosmossdk.io/math"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SharedParameter struct {
	priceChangeAssets     map[string]*operatortypes.PriceChange
	assetsDecimal         map[string]uint32
	assetsOperatorAVSInfo map[string]map[string]string
	stakerShare           map[string]sdkmath.LegacyDec
}

func UpdateShareOfStakerAndOperator(sharedParam *SharedParameter, assetID, stakerID, operatorAddr string, assetAmount sdkmath.Int) {
	priceChange := sharedParam.priceChangeAssets[assetID]
	assetDecimal := sharedParam.assetsDecimal[assetID]
	if avsAddr, ok := sharedParam.assetsOperatorAVSInfo[assetID][operatorAddr]; ok {
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
	if priceChangeAssets == nil || len(priceChangeAssets) == 0 {
		return nil
	}
	avsOperatorShareChange := make(map[string]sdkmath.LegacyDec, 0)
	assetsOperatorAVSInfo := make(map[string]map[string]string, 0)
	assetsDecimal := make(map[string]uint32)
	for assetID, priceChange := range priceChangeAssets {
		//get the decimal of asset
		assetInfo, err := k.restakingStateKeeper.GetStakingAssetInfo(ctx, assetID)
		if err != nil {
			return err
		}
		assetsDecimal[assetID] = assetInfo.AssetBasicInfo.Decimals
		if _, ok := assetsOperatorAVSInfo[assetID]; !ok {
			assetsOperatorAVSInfo[assetID] = make(map[string]string, 0)
		}
		//UpdateStateForAsset
		f := func(assetID string, keys []string, state *operatortypes.AssetOptedInState) error {
			newAssetUSDValue := CalculateShare(state.Amount, priceChange.NewPrice, assetInfo.AssetBasicInfo.Decimals, priceChange.Decimal)
			changeValue := newAssetUSDValue.Sub(state.Value)
			state.Value = newAssetUSDValue

			avsAddr := keys[1]
			avsOperator := string(types.GetJoinedStoreKey(keys[1], keys[2]))
			AddShareInMap(avsOperatorShareChange, avsAddr, changeValue)
			AddShareInMap(avsOperatorShareChange, avsOperator, changeValue)
			assetsOperatorAVSInfo[assetID][keys[2]] = avsAddr
			return nil
		}
		err = k.IterateUpdateAssetState(ctx, assetID, f)
		if err != nil {
			return err
		}
	}
	//BatchUpdateShareForAVSAndOperator
	err = k.BatchUpdateShareForAVSAndOperator(ctx, avsOperatorShareChange)
	if err != nil {
		return err
	}

	//update staker'suite share
	sharedParameter := &SharedParameter{
		priceChangeAssets:     priceChangeAssets,
		assetsDecimal:         assetsDecimal,
		assetsOperatorAVSInfo: assetsOperatorAVSInfo,
		stakerShare:           make(map[string]sdkmath.LegacyDec, 0),
	}
	stakerShareHandleFunc := func(stakerID, assetID, operatorAddr string, state *delegationtype.DelegationAmounts) error {
		UpdateShareOfStakerAndOperator(sharedParameter, assetID, stakerID, operatorAddr, state.CanUndelegationAmount)
		return nil
	}
	err = k.delegationKeeper.IterateDelegationState(ctx, stakerShareHandleFunc)
	if err != nil {
		return err
	}

	operatorShareHandleFunc := func(operatorAddr, assetID string, state *types.OperatorSingleAssetOrChangeInfo) error {
		UpdateShareOfStakerAndOperator(sharedParameter, assetID, "", operatorAddr, state.OperatorOwnAmountOrWantChangeValue)
		return nil
	}
	err = k.restakingStateKeeper.IteratorOperatorAssetState(ctx, operatorShareHandleFunc)
	if err != nil {
		return err
	}
	//BatchSetStakerShare
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
		[]byte{operatortypes.BytePrefixForOperatorAndChainIdToPrevConsKey},
	)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
	return nil
}

// EndBlock : update the assets' share when their prices change
func (k *Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	//todo: need to consider the calling order
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
