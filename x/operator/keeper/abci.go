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

func UpdateShareOfStakerAndOperator(sharedParam *SharedParameter, assetId, stakerId, operatorAddr string, assetAmount sdkmath.Int) {
	priceChange := sharedParam.priceChangeAssets[assetId]
	assetDecimal := sharedParam.assetsDecimal[assetId]
	if avsAddr, ok := sharedParam.assetsOperatorAVSInfo[assetId][operatorAddr]; ok {
		newAssetUSDValue := CalculateShare(assetAmount, priceChange.NewPrice, assetDecimal, priceChange.Decimal)
		key := string(types.GetJoinedStoreKey(avsAddr, stakerId, operatorAddr))
		AddShareInMap(sharedParam.stakerShare, key, newAssetUSDValue)
	}
}

// EndBlock : update the assets' share when their prices change
func (k *Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	priceChangeAssets, err := k.oracleKeeper.GetPriceChangeAssets(ctx)
	if err != nil {
		panic(err)
	}
	if priceChangeAssets == nil || len(priceChangeAssets) == 0 {
		return nil
	}
	avsOperatorShareChange := make(map[string]sdkmath.LegacyDec, 0)
	assetsOperatorAVSInfo := make(map[string]map[string]string, 0)
	assetsDecimal := make(map[string]uint32)
	for assetId, priceChange := range priceChangeAssets {
		//get the decimal of asset
		assetInfo, err := k.restakingStateKeeper.GetStakingAssetInfo(ctx, assetId)
		if err != nil {
			panic(err)
		}
		assetsDecimal[assetId] = assetInfo.AssetBasicInfo.Decimals
		if _, ok := assetsOperatorAVSInfo[assetId]; !ok {
			assetsOperatorAVSInfo[assetId] = make(map[string]string, 0)
		}
		//UpdateStateForAsset
		f := func(assetId string, keys []string, state *operatortypes.AssetOptedInState) error {
			newAssetUSDValue := CalculateShare(state.Amount, priceChange.NewPrice, assetInfo.AssetBasicInfo.Decimals, priceChange.Decimal)
			changeValue := newAssetUSDValue.Sub(state.Value)
			state.Value = newAssetUSDValue

			avsAddr := keys[1]
			avsOperator := string(types.GetJoinedStoreKey(keys[1], keys[2]))
			AddShareInMap(avsOperatorShareChange, avsAddr, changeValue)
			AddShareInMap(avsOperatorShareChange, avsOperator, changeValue)
			assetsOperatorAVSInfo[assetId][keys[2]] = avsAddr
			return nil
		}
		err = k.IterateUpdateAssetState(ctx, assetId, f)
		if err != nil {
			panic(err)
		}
	}
	//BatchUpdateShareForAVSAndOperator
	err = k.BatchUpdateShareForAVSAndOperator(ctx, avsOperatorShareChange)
	if err != nil {
		panic(err)
	}

	//update staker'suite share
	sharedParameter := &SharedParameter{
		priceChangeAssets:     priceChangeAssets,
		assetsDecimal:         assetsDecimal,
		assetsOperatorAVSInfo: assetsOperatorAVSInfo,
		stakerShare:           make(map[string]sdkmath.LegacyDec, 0),
	}
	stakerShareHandleFunc := func(stakerId, assetId, operatorAddr string, state *delegationtype.DelegationAmounts) error {
		UpdateShareOfStakerAndOperator(sharedParameter, assetId, stakerId, operatorAddr, state.CanUndelegationAmount)
		return nil
	}
	err = k.delegationKeeper.IterateDelegationState(ctx, stakerShareHandleFunc)
	if err != nil {
		panic(err)
	}

	operatorShareHandleFunc := func(operatorAddr, assetId string, state *types.OperatorSingleAssetOrChangeInfo) error {
		UpdateShareOfStakerAndOperator(sharedParameter, assetId, "", operatorAddr, state.OperatorOwnAmountOrWantChangeValue)
		return nil
	}
	err = k.restakingStateKeeper.IteratorOperatorAssetState(ctx, operatorShareHandleFunc)
	if err != nil {
		panic(err)
	}
	//BatchSetStakerShare
	err = k.BatchSetStakerShare(ctx, sharedParameter.stakerShare)
	if err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}
