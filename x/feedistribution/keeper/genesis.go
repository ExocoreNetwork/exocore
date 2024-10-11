package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	epochID := genState.Params.EpochIdentifier
	_, found := k.epochsKeeper.GetEpochInfo(ctx, epochID)
	if !found {
		// the panic is suitable here because it is being done at genesis, when the node
		// is not running. it means that the genesis file is malformed.
		panic("not found the epoch info")
	}

	// Set fee pool
	k.SetFeePool(ctx, &genState.FeePool)

	// Set all the validatorAccumulatedCommission
	for _, elem := range genState.ValidatorAccumulatedCommissions {
		k.SetValidatorAccumulatedCommission(ctx, sdk.ValAddress(elem.ValAddr), *elem.Commission)
	}

	// Set all the validatorCurrentRewards
	for _, elem := range genState.ValidatorCurrentRewardsList {
		k.SetValidatorCurrentRewards(ctx, sdk.ValAddress(elem.ValAddr), *elem.CurrentRewards)
	}

	// Set all the validatorOutstandingRewards
	for _, elem := range genState.ValidatorOutstandingRewardsList {
		k.SetValidatorOutstandingRewards(ctx, sdk.ValAddress(elem.ValAddr), *elem.OutstandingRewards)
	}

	// Set all the stakerRewards
	for _, elem := range genState.StakerOutstandingRewardsList {
		k.SetStakerRewards(ctx, elem.ValAddr, *elem.StakerOutstandingRewards)
	}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	feelPool := k.GetFeePool(ctx)

	genesis.FeePool = *feelPool
	var err error

	validatorData, err := k.GetAllValidatorData(ctx)
	if err != nil {
		panic(errorsmod.Wrap(err, "Error getting validator data").Error())
	}
	genesis.ValidatorAccumulatedCommissions = validatorData["ValidatorAccumulatedCommissions"].([]types.ValidatorAccumulatedCommissions)
	genesis.ValidatorCurrentRewardsList = validatorData["ValidatorCurrentRewardsList"].([]types.ValidatorCurrentRewardsList)
	genesis.ValidatorOutstandingRewardsList = validatorData["ValidatorOutstandingRewardsList"].([]types.ValidatorOutstandingRewardsList)
	genesis.StakerOutstandingRewardsList = validatorData["StakerOutstandingRewardsList"].([]types.StakerOutstandingRewardsList)
	return genesis
}
