package restaking_assets_manage

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/exocore/x/restaking_assets_manage/keeper"
	types2 "github.com/exocore/x/restaking_assets_manage/types"
)

// NewGenesisState - Create a new genesis state
func NewGenesisState(chain []*types2.ClientChainInfo, token []*types2.ClientChainTokenInfo) *types2.GenesisState {
	return &types2.GenesisState{
		DefaultSupportedClientChains:      chain,
		DefaultSupportedClientChainTokens: token,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *types2.GenesisState {
	return NewGenesisState([]*types2.ClientChainInfo{}, []*types2.ClientChainTokenInfo{})
}

// GetGenesisStateFromAppState returns x/auth GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Codec, appState map[string]json.RawMessage) types2.GenesisState {
	var genesisState types2.GenesisState

	if appState[types2.ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[types2.ModuleName], &genesisState)
	}

	return genesisState
}

// ValidateGenesis performs basic validation of restaking_assets_manage genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data types2.GenesisState) error {
	return nil
}

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data types2.GenesisState,
) {
}

// ExportGenesis export module status
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types2.GenesisState {
	return &types2.GenesisState{}
}
