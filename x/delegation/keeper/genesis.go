package keeper

import (
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
			for operator, wrappedAmount := range level2.PerOperatorAmounts {
				amount := wrappedAmount.Amount
				// #nosec G703 // already validated
				accAddress, _ := sdk.AccAddressFromBech32(operator)
				delegationParams := &delegationtype.DelegationOrUndelegationParams{
					ClientChainLzID: lzID,
					Action:          assetstype.DelegateTo,
					AssetsAddress:   assetAddressBytes.Bytes(),
					OperatorAddress: accAddress,
					StakerAddress:   stakerAddressBytes.Bytes(),
					OpAmount:        amount,
					// the uninitialized members are not used in this context
					// they are the LzNonce and TxHash
				}
				if err := k.delegateTo(ctx, delegationParams, false); err != nil {
					panic(err)
				}
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
