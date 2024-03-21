package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	k.SetParams(ctx, &data.Params)
	// client_chain.go
	for _, info := range data.ClientChains {
		k.SetClientChainInfo(ctx, &info)
	}
	// client_chain_asset.go
	for _, info := range data.Tokens {
		k.SetStakingAssetInfo(ctx, &info)
	}
	// operator_asset.go
	for _, level1 := range data.OperatorAssetInfos {
		// we have validated previously that the address is
		// the bech32 encoded address of sdk.AccAddress
		addr := level1.OperatorAddress
		for _, info := range level1.AssetIdAndInfos {
			k.SetOperatorAssetInfo(
				ctx, addr, info.AssetID, info.Info,
			)
		}
	}
	// staker_asset.go
	for _, level1 := range data.StakerAssetInfos {
		staker := level1.StakerID
		for _, info := range level1.AssetIdAndInfos {
			k.SetStakerAssetState(ctx, staker, info.AssetID, info.Info)
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	res := types.GenesisState{}
	// TODO
	return &res
}
