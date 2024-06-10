package keeper

import (
	sdkmath "cosmossdk.io/math"
)

// CalculateUSDValue assetUSDValue = (assetAmount*price)/(10^(asset.decimal+priceDecimal))
func CalculateUSDValue(assetAmount sdkmath.Int, price sdkmath.Int, assetDecimal uint32, priceDecimal uint8) sdkmath.LegacyDec {
	assetValue := assetAmount.Mul(price)
	assetValueDec := sdkmath.LegacyNewDecFromBigInt(assetValue.BigInt())
	// #nosec G701
	divisor := sdkmath.NewIntWithDecimal(1, int(assetDecimal)+int(priceDecimal))
	return assetValueDec.QuoInt(divisor)
}
