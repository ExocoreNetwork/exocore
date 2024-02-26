package restaking_assets_manage

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/keeper"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// NewStakingAssetsProposalHandler creates a governance handler to manage new
// proposal types.
func NewStakingAssetsProposalHandler(k *keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.RegisterClientChainProposal:
			return handleRegisterClientChainProposal(ctx, k, c)
		case *types.DeregisterClientChainProposal:
			return handleDeregisterClientChainProposal(ctx, k, c)
		case *types.RegisterAssetProposal:
			return handleRegisterAssetProposal(ctx, k, c)
		case *types.DeregisterAssetProposal:
			return handleDeregisterAssetProposal(ctx, k, c)
		default:
			return errorsmod.Wrapf(
				errortypes.ErrUnknownRequest,
				"unrecognized %s proposal content type: %T", types.ModuleName, c,
			)
		}
	}
}

func handleRegisterClientChainProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterClientChainProposal) error {
	err := k.RegisterClientChain(ctx, p.ClientChain)
	if err != nil {
		return err
	}
	// todo: emit related event
	return nil
}

func handleDeregisterClientChainProposal(ctx sdk.Context, k *keeper.Keeper, p *types.DeregisterClientChainProposal) error {
	err := k.DeregisterClientChain(ctx, p.ClientChainID)
	if err != nil {
		return err
	}
	// todo: emit related event
	return nil
}

func handleRegisterAssetProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterAssetProposal) error {
	err := k.RegisterAsset(ctx, p.Asset)
	if err != nil {
		return err
	}
	// todo: emit related event
	return nil
}

func handleDeregisterAssetProposal(ctx sdk.Context, k *keeper.Keeper, p *types.DeregisterAssetProposal) error {
	err := k.DeregisterAsset(ctx, p.AssetID)
	if err != nil {
		return err
	}
	// todo: emit related event
	return nil
}
