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
	stakers := make(map[string]struct{}, len(gs.Delegations))
	for _, level1 := range gs.Delegations {
		stakerID := level1.StakerID
		// validate staker ID
		if _, _, err := assetstypes.ValidateID(stakerID, true); err != nil {
			return errorsmod.Wrapf(err, "invalid staker ID %s", stakerID)
		}
		// check for duplicate stakers
		if _, ok := stakers[stakerID]; ok {
			return errorsmod.Wrapf(ErrInvalidGenesisData, "duplicate staker ID %s", stakerID)
		}
		stakers[stakerID] = struct{}{}
		assets := make(map[string]struct{}, len(level1.Delegations))
		for _, level2 := range level1.Delegations {
			assetID := level2.AssetID
			// validate asset ID
			if _, _, err := assetstypes.ValidateID(assetID, true); err != nil {
				return errorsmod.Wrapf(err, "invalid asset ID %s", assetID)
			}
			// check for duplicate assets
			if _, ok := assets[assetID]; ok {
				return errorsmod.Wrapf(ErrInvalidGenesisData, "duplicate asset ID %s", assetID)
			}
			assets[assetID] = struct{}{}
			givenTotal := level2.TotalDelegatedAmount
			if givenTotal.IsNegative() || givenTotal.IsNil() {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData, "invalid total delegated amount %d", givenTotal,
				)
			}
			calculatedTotal := sdk.ZeroInt()
			operators := make(map[string]struct{}, len(level2.PerOperatorAmounts))
			for operator, wrappedAmount := range level2.PerOperatorAmounts {
				// check supplied amount
				if wrappedAmount == nil {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData, "nil operator amount for operator %s", operator,
					)
				}
				amount := wrappedAmount.Amount
				if amount.IsNegative() || amount.IsNil() {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"invalid operator amount %d for operator %s", amount, operator,
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
				calculatedTotal = calculatedTotal.Add(amount)
			}
			if !givenTotal.Equal(calculatedTotal) {
				return errorsmod.Wrapf(
					ErrInvalidGenesisData,
					"total delegated amount %d does not match calculated total %d for asset %s",
					givenTotal, calculatedTotal, assetID,
				)
			}
		}
	}
	return nil
}
