package keeper

import (
	errorsmod "cosmossdk.io/errors"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// InitGenesis initializes the module's state from a provided genesis state.
// Since this action typically occurs on chain starts, this function is allowed to panic.
func (k Keeper) InitGenesis(
	ctx sdk.Context,
	gs delegationtype.GenesisState,
) []abci.ValidatorUpdate {
	// TODO(mm): is it possible to parallelize these without using goroutines?
	for _, level1 := range gs.Delegations {
		stakerID := level1.StakerID
		// #nosec G703 // already validated
		stakerAddress, lzID, _ := assetstype.ParseID(stakerID)
		// we have checked IsHexAddress already
		stakerAddressBytes := common.HexToAddress(stakerAddress)
		for _, level2 := range level1.Delegations {
			assetID := level2.AssetID
			// #nosec G703 // already validated
			assetAddress, _, _ := assetstype.ParseID(assetID)
			// we have checked IsHexAddress already
			assetAddressBytes := common.HexToAddress(assetAddress)
			for _, level3 := range level2.PerOperatorAmounts {
				operator := level3.Key
				wrappedAmount := level3.Value
				amount := wrappedAmount.Amount
				// #nosec G703 // already validated
				accAddress, _ := sdk.AccAddressFromBech32(operator)
				delegationParams := &delegationtype.DelegationOrUndelegationParams{
					ClientChainID:   lzID,
					Action:          assetstype.DelegateTo,
					AssetsAddress:   assetAddressBytes.Bytes(),
					OperatorAddress: accAddress,
					StakerAddress:   stakerAddressBytes.Bytes(),
					OpAmount:        amount,
					// the uninitialized members are not used in this context
					// they are the LzNonce and TxHash
				}
				if err := k.delegateTo(ctx, delegationParams, false); err != nil {
					panic(errorsmod.Wrap(err, "failed to delegate to operator"))
				}
			}
		}
	}
	for _, association := range gs.Associations {
		stakerID := association.StakerID
		operatorAddress := association.Operator
		// #nosec G703 // already validated
		stakerAddress, clientChainID, _ := assetstype.ParseID(stakerID)
		// we have checked IsHexAddress already
		stakerAddressBytes := common.FromHex(stakerAddress)
		// #nosec G703 // already validated
		accAddress, _ := sdk.AccAddressFromBech32(operatorAddress)
		// this can only fail if the operator is not registered
		if err := k.AssociateOperatorWithStaker(
			ctx, clientChainID, accAddress, stakerAddressBytes,
		); err != nil {
			panic(errorsmod.Wrap(err, "failed to associate operator with staker"))
		}
	}

	// init the state from the general exporting genesis file
	err := k.SetAllDelegationStates(ctx, gs.DelegationStates)
	if err != nil {
		panic(err)
	}
	err = k.SetAllStakerList(ctx, gs.StakersByOperator)
	if err != nil {
		panic(err)
	}
	err = k.SetUndelegationRecords(ctx, gs.Undelegations)
	if err != nil {
		panic(err)
	}
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *delegationtype.GenesisState {
	res := delegationtype.GenesisState{
		Delegations: []delegationtype.DelegationsByStaker{},
	}
	delegationStates, err := k.AllDelegationStates(ctx)
	if err != nil {
		panic(err)
	}
	res.DelegationStates = delegationStates

	stakerList, err := k.AllStakerList(ctx)
	if err != nil {
		panic(err)
	}
	res.StakersByOperator = stakerList

	undelegations, err := k.AllUndelegations(ctx)
	if err != nil {
		panic(err)
	}
	res.Undelegations = undelegations
	// mark the exported genesis as general import
	res.IsGeneralInit = true
	return &res
}
