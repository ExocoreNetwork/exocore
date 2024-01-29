package keeper

import (
	sdkmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	for assetId, priceChange := range priceChangeAssets {
		//get the decimal of asset
		assetInfo, err := k.restakingStateKeeper.GetStakingAssetInfo(ctx, assetId)
		if err != nil {
			panic(err)
		}
		//UpdateOperatorAVSAssetsState
		f := func(assetId string, keys []string, state *operatortypes.AssetOptedInState) error {
			newAssetValue := state.Amount.Mul(priceChange.NewPrice).Mul(sdkmath.NewIntWithDecimal(1, int(operatortypes.UsdValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetInfo.AssetBasicInfo.Decimals)+int(priceChange.Decimal)))
			newAssetUsdValue := sdkmath.LegacyNewDecFromBigIntWithPrec(newAssetValue.BigInt(), int64(operatortypes.UsdValueDefaultDecimal))
			changeValue := newAssetUsdValue.Sub(state.Value)
			state.Value = newAssetUsdValue

			avsAddr := keys[1]
			avsOperator := string(types.GetJoinedStoreKey(keys[1], keys[2]))
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

	return []abci.ValidatorUpdate{}
}
