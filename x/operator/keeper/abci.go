package keeper

import (
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	delegationtype "github.com/exocore/x/delegation/types"
	operatortypes "github.com/exocore/x/operator/types"
	"github.com/exocore/x/restaking_assets_manage/types"
)

// EndBlock : update the assets' share when their prices change
func (k Keeper) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	priceChangeAssets, err := k.oracleKeeper.GetPriceChangeAssets(ctx)
	if err != nil {
		panic(err)
	}
	avsOperatorShareChange := make(map[string]sdkmath.LegacyDec, 0)
	assetsOperator := make(map[string]map[string]string, 0)
	assetsDecimal := make(map[string]uint32)
	for assetId, priceChange := range priceChangeAssets {
		//get the decimal of asset
		assetInfo, err := k.restakingStateKeeper.GetStakingAssetInfo(ctx, assetId)
		if err != nil {
			panic(err)
		}
		assetsDecimal[assetId] = assetInfo.AssetBasicInfo.Decimals
		if _, ok := assetsOperator[assetId]; !ok {
			assetsOperator[assetId] = make(map[string]string, 0)
		}
		//UpdateOperatorAVSAssetsState
		f := func(assetId string, keys []string, state *operatortypes.AssetOptedInState) error {
			newAssetValue := state.Amount.Mul(priceChange.NewPrice).Mul(sdkmath.NewIntWithDecimal(1, int(operatortypes.USDValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetInfo.AssetBasicInfo.Decimals)+int(priceChange.Decimal)))
			newAssetUSDValue := sdkmath.LegacyNewDecFromBigIntWithPrec(newAssetValue.BigInt(), int64(operatortypes.USDValueDefaultDecimal))
			changeValue := newAssetUSDValue.Sub(state.Value)
			state.Value = newAssetUSDValue

			avsAddr := keys[1]
			avsOperator := string(types.GetJoinedStoreKey(keys[1], keys[2]))
			avsOperatorShareChange[avsAddr] = avsOperatorShareChange[avsAddr].Add(changeValue)
			if value, ok := avsOperatorShareChange[avsAddr]; ok {
				avsOperatorShareChange[avsAddr] = value.Add(changeValue)
			} else {
				avsOperatorShareChange[avsAddr] = changeValue
			}
			if value, ok := avsOperatorShareChange[avsOperator]; ok {
				avsOperatorShareChange[avsOperator] = value.Add(changeValue)
			} else {
				avsOperatorShareChange[avsOperator] = changeValue
			}

			assetsOperator[assetId][keys[2]] = avsAddr
			return nil
		}
		err = k.IterateUpdateOperatorAVSAssets(ctx, assetId, f)
		if err != nil {
			panic(err)
		}
	}
	//BatchUpdateAVSAndOperatorTotalValue
	err = k.BatchUpdateAVSAndOperatorTotalValue(ctx, avsOperatorShareChange)
	if err != nil {
		panic(err)
	}

	//update staker's share
	stakerOperatorNewShare := make(map[string]sdkmath.LegacyDec, 0)
	stakerShareHandleFunc := func(restakerId, assetId, operatorAddr string, state *delegationtype.DelegationAmounts) error {
		priceChange := priceChangeAssets[assetId]
		assetDecimal := assetsDecimal[assetId]
		if avsAddr, ok := assetsOperator[assetId][operatorAddr]; ok {
			newAssetValue := state.CanUndelegationAmount.Mul(priceChange.NewPrice).Mul(sdkmath.NewIntWithDecimal(1, int(operatortypes.USDValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetDecimal)+int(priceChange.Decimal)))
			newAssetUSDValue := sdkmath.LegacyNewDecFromBigIntWithPrec(newAssetValue.BigInt(), int64(operatortypes.USDValueDefaultDecimal))
			key := string(types.GetJoinedStoreKey(avsAddr, restakerId, operatorAddr))
			if value, ok := stakerOperatorNewShare[key]; ok {
				stakerOperatorNewShare[key] = value.Add(newAssetUSDValue)
			} else {
				stakerOperatorNewShare[key] = newAssetUSDValue
			}
		}
		return nil
	}
	err = k.delegationKeeper.IteratorDelegationState(ctx, stakerShareHandleFunc)
	if err != nil {
		panic(err)
	}

	operatorShareHandleFunc := func(operatorAddr, assetId string, state *types.OperatorSingleAssetOrChangeInfo) error {
		priceChange := priceChangeAssets[assetId]
		assetDecimal := assetsDecimal[assetId]
		if avsAddr, ok := assetsOperator[assetId][operatorAddr]; ok {
			newAssetValue := state.OperatorOwnAmountOrWantChangeValue.Mul(priceChange.NewPrice).Mul(sdkmath.NewIntWithDecimal(1, int(operatortypes.USDValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetDecimal)+int(priceChange.Decimal)))
			newAssetUSDValue := sdkmath.LegacyNewDecFromBigIntWithPrec(newAssetValue.BigInt(), int64(operatortypes.USDValueDefaultDecimal))
			key := string(types.GetJoinedStoreKey(avsAddr, "", operatorAddr))
			if value, ok := stakerOperatorNewShare[key]; ok {
				stakerOperatorNewShare[key] = value.Add(newAssetUSDValue)
			} else {
				stakerOperatorNewShare[key] = newAssetUSDValue
			}
		}
		return nil
	}
	err = k.restakingStateKeeper.IteratorOperatorAssetState(ctx, operatorShareHandleFunc)
	if err != nil {
		panic(err)
	}
	//BatchSetAVSOperatorStakerShare
	err = k.BatchSetAVSOperatorStakerShare(ctx, stakerOperatorNewShare)
	if err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}
