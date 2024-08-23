package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) CalculateActualThreshold(ctx sdk.Context, total sdkmath.LegacyDec, avs string) (t sdkmath.LegacyDec) {
	usd, err := k.operatorKeeper.GetAVSUSDValue(ctx, avs)
	if err != nil {
		return sdkmath.LegacyZeroDec()
	}

	if usd.IsZero() || total.IsZero() {
		return sdkmath.LegacyZeroDec()
	}
	return total.Quo(usd).Mul(sdk.NewDec(100))
}

func Difference(a, b []string) []string {
	var different []string //nolint:prealloc

	diffMap := make(map[string]bool)

	// Add all elements of a to the map
	for _, item := range a {
		diffMap[item] = true
	}

	// Remove elements found in b from the map and collect differences
	for _, item := range b {
		if diffMap[item] {
			delete(diffMap, item)
		} else {
			different = append(different, item)
		}
	}

	// Calculate the final size for the different slice
	finalSize := len(different) + len(diffMap)

	// Pre-allocate the different slice with the final size
	different = make([]string, 0, finalSize)

	// Add remaining elements from the map to different
	for item := range diffMap {
		different = append(different, item)
	}

	return different
}
