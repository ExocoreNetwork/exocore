package keeper

import (
	sdkmath "cosmossdk.io/math"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
)

type LegacyDecMap map[string]sdkmath.LegacyDec

func AddShareInMap(shareMap map[string]sdkmath.LegacyDec, key string, addValue sdkmath.LegacyDec) {
	if value, ok := shareMap[key]; ok {
		shareMap[key] = value.Add(addValue)
	} else {
		shareMap[key] = addValue
	}
}

// CalculateUSDValue assetUSDValue = (assetAmount*price*10^USDValueDefaultDecimal)/(10^(asset.decimal+priceDecimal))
func CalculateUSDValue(assetAmount sdkmath.Int, price sdkmath.Int, assetDecimal uint32, priceDecimal uint8) sdkmath.LegacyDec {
	// #nosec G701
	assetValue := assetAmount.Mul(price).Mul(sdkmath.NewIntWithDecimal(1, int(operatortypes.USDValueDefaultDecimal))).Quo(sdkmath.NewIntWithDecimal(1, int(assetDecimal)+int(priceDecimal)))
	// #nosec G701
	assetUSDValue := sdkmath.LegacyNewDecFromBigIntWithPrec(assetValue.BigInt(), int64(operatortypes.USDValueDefaultDecimal))
	return assetUSDValue
}
