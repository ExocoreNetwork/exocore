package restaking_assets_manage

import (
	"encoding/json"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/restaking_assets_manage/keeper"
	restakingtype "github.com/exocore/x/restaking_assets_manage/types"
)

// NewGenesisState - Create a new genesis state
func NewGenesisState(chain []*restakingtype.ClientChainInfo, token []*restakingtype.ClientChainTokenInfo) *restakingtype.GenesisState {
	return &restakingtype.GenesisState{
		DefaultSupportedClientChains:      chain,
		DefaultSupportedClientChainTokens: token,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *restakingtype.GenesisState {
	// todo: set eth as client chain and usdt as asset in the genesis state
	ethClientChain := &restakingtype.ClientChainInfo{
		ChainName:              "ethereum",
		ChainMetaInfo:          "ethereum blockchain",
		OriginChainId:          1,
		FinalityNeedBlockDelay: 10,
		LayerZeroChainId:       101,
		AddressLength:          20,
	}
	usdtClientChainAsset := &restakingtype.ClientChainTokenInfo{
		Name:             "Tether USD",
		Symbol:           "USDT",
		Address:          "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		Decimals:         6,
		LayerZeroChainId: ethClientChain.LayerZeroChainId,
		AssetMetaInfo:    "Tether USD token",
	}
	totalSupply, _ := sdk.NewIntFromString("40022689732746729")
	usdtClientChainAsset.TotalSupply = totalSupply
	return NewGenesisState([]*restakingtype.ClientChainInfo{ethClientChain}, []*restakingtype.ClientChainTokenInfo{usdtClientChainAsset})
}

// GetGenesisStateFromAppState returns x/restaking_assets_manage GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) restakingtype.GenesisState {
	var genesisState restakingtype.GenesisState

	if appState[restakingtype.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[restakingtype.ModuleName], &genesisState)
	}

	return genesisState
}

// ValidateGenesis performs basic validation of restaking_assets_manage genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data restakingtype.GenesisState) error {
	// todo: check the validation of client chain and token info
	return nil
}

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data restakingtype.GenesisState,
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
		err = k.SetStakingAssetInfo(c, &restakingtype.StakingAssetInfo{
			AssetBasicInfo:     asset,
			StakingTotalAmount: math.NewInt(0),
		})
		if err != nil {
			panic(err)
		}
	}
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *restakingtype.GenesisState {
	clientChainList := make([]*restakingtype.ClientChainInfo, 0)
	c := sdk.UnwrapSDKContext(ctx)
	clientChainInfo, _ := k.GetAllClientChainInfo(c)
	for _, v := range clientChainInfo {
		clientChainList = append(clientChainList, v)
	}

	clientChainAssetsList := make([]*restakingtype.ClientChainTokenInfo, 0)
	clientChainAssets, _ := k.GetAllStakingAssetsInfo(c)
	for _, v := range clientChainAssets {
		clientChainAssetsList = append(clientChainAssetsList, v.AssetBasicInfo)
	}
	return &restakingtype.GenesisState{
		DefaultSupportedClientChains:      clientChainList,
		DefaultSupportedClientChainTokens: clientChainAssetsList,
	}
}
