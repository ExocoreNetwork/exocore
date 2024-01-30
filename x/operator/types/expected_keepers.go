package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	delegationtype "github.com/exocore/x/delegation/types"
)

type ExpectDelegationInterface interface {
	GetDelegationStateByOperatorAndAssetList(ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{}) (map[string]map[string]delegationtype.DelegationAmounts, error)
	IteratorDelegationState(ctx sdk.Context, f func(restakerId, assetId, operatorAddr string, state *delegationtype.DelegationAmounts) error) error
	UpdateDelegationState(ctx sdk.Context, stakerId string, assetId string, delegationAmounts map[string]*delegationtype.DelegationAmounts) (err error)

	UpdateStakerDelegationTotalAmount(ctx sdk.Context, stakerId string, assetId string, opAmount sdkmath.Int) error
}

type PriceChange struct {
	OriginalPrice sdkmath.Int
	NewPrice      sdkmath.Int
	Decimal       uint8
}
type ExpectOracleInterface interface {
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetsId string) (sdkmath.Int, uint8, error)
	GetPriceChangeAssets(ctx sdk.Context) (map[string]*PriceChange, error)
}

type MockOracle struct{}

func (MockOracle) GetSpecifiedAssetsPrice(ctx sdk.Context, assetsId string) (sdkmath.Int, uint8, error) {
	return sdkmath.NewInt(1), 0, nil
}

func (MockOracle) GetPriceChangeAssets(ctx sdk.Context) (map[string]*PriceChange, error) {
	return nil, nil
}

type MockAVS struct{}

func (MockAVS) GetAvsSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error) {
	return nil, nil
}

func (MockAVS) GetAvsSlashContract(ctx sdk.Context, avsAddr string) (string, error) {
	return "", nil
}

type ExpectAvsInterface interface {
	GetAvsSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error)
	GetAvsSlashContract(ctx sdk.Context, avsAddr string) (string, error)
}
