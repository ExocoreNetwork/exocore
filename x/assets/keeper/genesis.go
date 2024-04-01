package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	if err := k.SetParams(ctx, &data.Params); err != nil {
		panic(err)
	}
	// client_chain.go
	for _, infoCopy := range data.ClientChains {
		info := infoCopy // prevent implicit memory aliasing
		if err := k.SetClientChainInfo(ctx, &info); err != nil {
			panic(err)
		}
	}
	// client_chain_asset.go
	for _, infoCopy := range data.Tokens {
		info := infoCopy // prevent implicit memory aliasing
		if err := k.SetStakingAssetInfo(ctx, &info); err != nil {
			panic(err)
		}
	}
	// operator_asset.go
	for _, level1 := range data.OperatorAssetInfos {
		// we have validated previously that the address is
		// the bech32 encoded address of sdk.AccAddress
		addr := level1.OperatorAddress
		for _, info := range level1.AssetIdAndInfos {
			if err := k.SetOperatorAssetInfo(
				ctx, addr, info.AssetID, info.Info,
			); err != nil {
				panic(err)
			}
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
