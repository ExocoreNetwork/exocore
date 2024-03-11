package types

import "cosmossdk.io/math"

// OptedInAssetStateChange This is a struct to describe the desired change that matches with the OptedInAssetState
type OptedInAssetStateChange struct {
	ChangeForAmount math.Int
	ChangeForValue  math.LegacyDec
}
