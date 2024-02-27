package types

import (
	sdkmath "cosmossdk.io/math"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExpectDelegationInterface interface {
	DelegationStateByOperatorAssets(ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{}) (map[string]map[string]delegationtype.DelegationAmounts, error)
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

// add for dogfood

type SlashKeeper interface {
	IsOperatorFrozen(ctx sdk.Context, addr sdk.AccAddress) bool
}

type RedelegationKeeper interface {
	AppChainInfoIsExist(ctx sdk.Context, chainId string) bool
}

type OperatorConsentHooks interface {
	// This hook is called when an operator opts in to a chain.
	AfterOperatorOptIn(
		ctx sdk.Context,
		addr sdk.AccAddress,
		chainId string,
		pubKey tmprotocrypto.PublicKey,
	)
	// This hook is called when an operator's consensus key is replaced for
	// a chain.
	AfterOperatorKeyReplacement(
		ctx sdk.Context,
		addr sdk.AccAddress,
		oldKey tmprotocrypto.PublicKey,
		newKey tmprotocrypto.PublicKey,
		chainId string,
	)
	// This hook is called when an operator opts out of a chain.
	AfterOperatorOptOutInitiated(
		ctx sdk.Context,
		addr sdk.AccAddress,
		chainId string,
		key tmprotocrypto.PublicKey,
	)
}
