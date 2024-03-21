package assets

import (
	"encoding/json"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"

	"cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState - Create a new genesis state
func NewGenesisState(chain []*assetstype.ClientChainInfo, token []*assetstype.AssetInfo) *assetstype.GenesisState {
	return &assetstype.GenesisState{
		DefaultSupportedClientChains:      chain,
		DefaultSupportedClientChainTokens: token,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *assetstype.GenesisState {
	// todo: set eth as client chain and usdt as asset in the genesis state
	ethClientChain := &assetstype.ClientChainInfo{
		Name:               "ethereum",
		MetaInfo:           "ethereum blockchain",
		ChainId:            1,
		FinalizationBlocks: 10,
		LayerZeroChainID:   101,
		AddressLength:      20,
	}
	usdtClientChainAsset := &assetstype.AssetInfo{
		Name:             "Tether USD",
		Symbol:           "USDT",
		Address:          "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		Decimals:         6,
		LayerZeroChainID: ethClientChain.LayerZeroChainID,
		MetaInfo:         "Tether USD token",
	}
	totalSupply, _ := sdk.NewIntFromString("40022689732746729")
	usdtClientChainAsset.TotalSupply = totalSupply
	return NewGenesisState([]*assetstype.ClientChainInfo{ethClientChain}, []*assetstype.AssetInfo{usdtClientChainAsset})
}

// GetGenesisStateFromAppState returns x/assets GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) assetstype.GenesisState {
	var genesisState assetstype.GenesisState

	if appState[assetstype.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[assetstype.ModuleName], &genesisState)
	}

	return genesisState
}

// ValidateGenesis performs basic validation of assets genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(assetstype.GenesisState) error {
	// todo: check the validation of client chain and token info
	return nil
}

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data assetstype.GenesisState,
) {
	// todo: might need to sort the clientChains and tokens before handling.

	c := sdk.UnwrapSDKContext(ctx)
	var err error
	// save default supported client chain
	for _, chain := range data.DefaultSupportedClientChains {
		err = k.SetClientChainInfo(c, chain)
		if err != nil {
			panic(err)
		}
	}
	// save default supported client chain assets
	for _, asset := range data.DefaultSupportedClientChainTokens {
		err = k.SetStakingAssetInfo(c, &assetstype.StakingAssetInfo{
			AssetBasicInfo:     asset,
			StakingTotalAmount: math.NewInt(0),
		})
		if err != nil {
			panic(err)
		}
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *assetstype.GenesisState {
	clientChainList := make([]*assetstype.ClientChainInfo, 0)
	c := sdk.UnwrapSDKContext(ctx)
	clientChainInfo, _ := k.GetAllClientChainInfo(c)
	for _, v := range clientChainInfo {
		clientChainList = append(clientChainList, v)
	}

	clientChainAssetsList := make([]*assetstype.AssetInfo, 0)
	clientChainAssets, _ := k.GetAllStakingAssetsInfo(c)
	for _, v := range clientChainAssets {
		clientChainAssetsList = append(clientChainAssetsList, v.AssetBasicInfo)
	}
	return &assetstype.GenesisState{
		DefaultSupportedClientChains:      clientChainList,
		DefaultSupportedClientChainTokens: clientChainAssetsList,
	}
}
