package keeper

import (
	"fmt"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
// Since this action typically occurs on chain starts, this function is allowed to panic.
func (k Keeper) InitGenesis(
	ctx sdk.Context,
	genState delegationtype.GenesisState,
) []abci.ValidatorUpdate {
	for _, a := range genState.DelegationsByStakerAssetOperator {
		stakerID := a.StakerID
		for _, b := range a.DelegationsByAssetOperator {
			assetID := b.AssetID
			// check that the asset is registered
			if !k.assetsKeeper.IsStakingAsset(ctx, assetID) {
				panic(
					fmt.Sprintf(
						"%s: %s", assetstype.ErrNoClientChainAssetKey,
						fmt.Sprintf("input assetID is:%s", assetID),
					),
				)
			}
			for _, c := range b.DelegationsByOperator {
				operatorAddress := c.OperatorAddress
				amount := c.Amount // delegation amount

				// the check that the operator is registered cannot be made here, since
				// the genesis of the operator module runs after the delegation genesis,
				// although this order can be changed.
				// if there is an operator within this genesis that is not in the
				// operator module genesis, the operator module will reject that genesis.
				// this is because the operator module checks for consistency between
				// itself and this module at the time of its genesis.

				// at genesis, the operator cannot be frozen so skip that check.
				// validate that enough deposits exist before delegation.
				// note that these deposits are by stakerID for assetID, and not per operator.
				info, err := k.assetsKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetID)
				if err != nil {
					panic(err)
				}
				if amount.GT(info.WithdrawableAmount) {
					panic(
						fmt.Sprintf(
							"delegated amount %s is greater than the staker's asset amount %s",
							amount.String(), info.WithdrawableAmount.String(),
						),
					)
				}
				delegatorAndAmount := make(map[string]*delegationtype.DelegationAmounts, 1)
				delegatorAndAmount[operatorAddress] = &delegationtype.DelegationAmounts{
					UndelegatableAmount: amount,
				}
				if err := k.UpdateDelegationState(
					ctx, stakerID, assetID, delegatorAndAmount,
				); err != nil {
					panic(err)
				}
				if err := k.UpdateStakerDelegationTotalAmount(
					ctx, stakerID, assetID, amount,
				); err != nil {
					panic(err)
				}
				// reduce the amount of available deposit for the next operator iteration.
				if err := k.assetsKeeper.UpdateStakerAssetState(
					ctx, stakerID, assetID, assetstype.StakerSingleAssetChangeInfo{
						WithdrawableAmount: amount.Neg(),
					}); err != nil {
					panic(err)
				}
				// we have checked that delegation amount > deposit amount for each asset.
				// we don't need to check for the total amount, since this genesis only handles
				// delegation amounts (others are implicitly zero).
			}
		}
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the module's exported genesis
func (Keeper) ExportGenesis(sdk.Context) *delegationtype.GenesisState {
	genesis := delegationtype.DefaultGenesis()
	// TODO
	return genesis
}
