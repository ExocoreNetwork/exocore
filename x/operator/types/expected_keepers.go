package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExpectDelegationInterface interface {
	GetDelegationStateByOperator(ctx sdk.Context, operatorAddr string) (map[string]sdkmath.Int, error)
}

type ExpectOracleInterface interface {
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetsId string) (sdkmath.Int, uint8, error)
}
