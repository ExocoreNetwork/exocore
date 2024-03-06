package types

import (
	sdkmath "cosmossdk.io/math"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ExpectDelegationInterface interface {
	DelegationStateByOperatorAssets(ctx sdk.Context, operatorAddr string, assetsFilter map[string]interface{}) (map[string]map[string]delegationtype.DelegationAmounts, error)
	IterateDelegationState(ctx sdk.Context, f func(restakerID, assetID, operatorAddr string, state *delegationtype.DelegationAmounts) error) error
	UpdateDelegationState(ctx sdk.Context, stakerID string, assetID string, delegationAmounts map[string]*delegationtype.DelegationAmounts) (err error)

	UpdateStakerDelegationTotalAmount(ctx sdk.Context, stakerID string, assetID string, opAmount sdkmath.Int) error
}

type PriceChange struct {
	OriginalPrice sdkmath.Int
	NewPrice      sdkmath.Int
	Decimal       uint8
}
type ExpectOracleInterface interface {
	GetSpecifiedAssetsPrice(ctx sdk.Context, assetsID string) (sdkmath.Int, uint8, error)
	GetPriceChangeAssets(ctx sdk.Context) (map[string]*PriceChange, error)
}

type MockOracle struct{}

func (MockOracle) GetSpecifiedAssetsPrice(_ sdk.Context, _ string) (sdkmath.Int, uint8, error) {
	return sdkmath.NewInt(1), 0, nil
}

func (MockOracle) GetPriceChangeAssets(_ sdk.Context) (map[string]*PriceChange, error) {
	// use USDT as the mock asset
	ret := make(map[string]*PriceChange, 0)
	usdtAssetID := "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"
	ret[usdtAssetID] = &PriceChange{
		NewPrice:      sdkmath.NewInt(1),
		OriginalPrice: sdkmath.NewInt(1),
		Decimal:       0,
	}
	return nil, nil
}

type MockAVS struct{}

func (MockAVS) GetAvsSupportedAssets(_ sdk.Context, _ string) (map[string]interface{}, error) {
	// set USDT as the default asset supported by AVS
	ret := make(map[string]interface{})
	usdtAssetID := "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"
	ret[usdtAssetID] = nil
	return ret, nil
}

func (MockAVS) GetAvsSlashContract(_ sdk.Context, _ string) (string, error) {
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
	AppChainInfoIsExist(ctx sdk.Context, chainID string) bool
}

type OperatorConsentHooks interface {
	// This hook is called when an operator opts in to a chain.
	AfterOperatorOptIn(
		ctx sdk.Context,
		addr sdk.AccAddress,
		chainID string,
		pubKey tmprotocrypto.PublicKey,
	)
	// This hook is called when an operator's consensus key is replaced for
	// a chain.
	AfterOperatorKeyReplacement(
		ctx sdk.Context,
		addr sdk.AccAddress,
		oldKey tmprotocrypto.PublicKey,
		newKey tmprotocrypto.PublicKey,
		chainID string,
	)
	// This hook is called when an operator opts out of a chain.
	AfterOperatorOptOutInitiated(
		ctx sdk.Context,
		addr sdk.AccAddress,
		chainID string,
		key tmprotocrypto.PublicKey,
	)
}
