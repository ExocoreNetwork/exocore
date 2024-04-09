package types

import (
	errorsmod "cosmossdk.io/errors"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesis returns a new genesis state with the given inputs.
func NewGenesis(
	delegations []DelegationsByStaker,
) *GenesisState {
	return &GenesisState{
		Delegations: delegations,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesis([]DelegationsByStaker{})
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// TODO(mm): this can be a very big hash table and impact system performance.
	// This is likely to be the biggest one amongst the three, and the others
	// are garbage collected within the loop anyway. Maybe reordering the genesis
	// structure could potentially help with this.
	stakers := make(map[string]struct{}, len(gs.Delegations))
	for _, level1 := range gs.Delegations {
		stakerID := level1.StakerID
		// validate staker ID
		var stakerClientChainID uint64
		var err error
		if _, stakerClientChainID, err = assetstypes.ValidateID(stakerID, true); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData, "invalid staker ID %s: %s", stakerID, err,
			)
		}
		// check for duplicate stakers
		if _, ok := stakers[stakerID]; ok {
			return errorsmod.Wrapf(ErrInvalidGenesisData, "duplicate staker ID %s", stakerID)
		}
		stakers[stakerID] = struct{}{}
		assets := make(map[string]struct{}, len(level1.Delegations))
		for _, level2 := range level1.Delegations {
			assetID := level2.AssetID
			// check for duplicate assets
			if _, ok := assets[assetID]; ok {
				return errorsmod.Wrapf(ErrInvalidGenesisData, "duplicate asset ID %s", assetID)
			}
			assets[assetID] = struct{}{}
			// validate asset ID
			var assetClientChainID uint64
			if _, assetClientChainID, err = assetstypes.ValidateID(assetID, true); err != nil {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData, "invalid asset ID %s: %s", assetID, err,
				)
			}
			if assetClientChainID != stakerClientChainID {
				// a staker from chain A is delegating an asset on chain B, which is not
				// something we support right now.
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"asset %s client chain ID %d does not match staker %s client chain ID %d",
					assetID, assetClientChainID, stakerID, stakerClientChainID,
				)
			}
			operators := make(map[string]struct{}, len(level2.PerOperatorAmounts))
			for _, level3 := range level2.PerOperatorAmounts {
				operator := level3.Key
				wrappedAmount := level3.Value
				// check supplied amount
				if wrappedAmount == nil {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData, "nil operator amount for operator %s", operator,
					)
				}
				amount := wrappedAmount.Amount
				if amount.IsNil() || amount.IsNegative() {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"invalid operator amount %s for operator %s", amount, operator,
					)
				}
				// check operator address
				if _, err := sdk.AccAddressFromBech32(operator); err != nil {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"invalid operator address for operator %s", operator,
					)
				}
				// check for duplicate operators
				if _, ok := operators[operator]; ok {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"duplicate operator %s for asset %s", operator, assetID,
					)
				}
				operators[operator] = struct{}{}
			}
		}
	}
	return nil
}
