package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	delegationtype "github.com/exocore/x/delegation/types"
)

type ExpectDelegationInterface interface {
	GetDelegationStateByOperatorAndAssets(ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{}) (map[string]map[string]delegationtype.DelegationAmounts, error)
	IterateDelegationState(ctx sdk.Context, f func(restakerId, assetId, operatorAddr string, state *delegationtype.DelegationAmounts) error) error
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
	//use USDT as the mock asset
	ret := make(map[string]*PriceChange, 0)
	usdtAssetId := "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"
	ret[usdtAssetId] = &PriceChange{
		NewPrice:      sdkmath.NewInt(1),
		OriginalPrice: sdkmath.NewInt(1),
		Decimal:       0,
	}
	return nil, nil
}

type MockAVS struct{}

func (MockAVS) GetAvsSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error) {
	//set USDT as the default asset supported by AVS
	ret := make(map[string]interface{})
	usdtAssetId := "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"
	ret[usdtAssetId] = nil
	return ret, nil
}

func (MockAVS) GetAvsSlashContract(ctx sdk.Context, avsAddr string) (string, error) {
	return "", nil
}

type ExpectAvsInterface interface {
	GetAvsSupportedAssets(ctx sdk.Context, avsAddr string) (map[string]interface{}, error)
	GetAvsSlashContract(ctx sdk.Context, avsAddr string) (string, error)
}
