package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	delegationtype "github.com/exocore/x/delegation/types"
)

type ExpectDelegationInterface interface {
	GetDelegationStateByOperatorAndAssetList(ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{}) (map[string]map[string]sdkmath.Int, error)
	UpdateDelegationState(ctx sdk.Context, stakerId string, assetId string, delegationAmounts map[string]*delegationtype.DelegationAmounts) (err error)
}

type ExpectOracleInterface interface {
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetsId string) (sdkmath.Int, uint8, error)
}

type ExpectAvsInterface interface {
	GetAvsSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error)
	GetAvsSlashContract(ctx sdk.Context, avsAddr string) (string, error)
}
