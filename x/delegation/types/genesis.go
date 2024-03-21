package types

import (
	errorsmod "cosmossdk.io/errors"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// there is only one item to validate, so we don't sum anything up.
	// we can only validate that the parameters are present, and that
	// the stakerIDs and assetIDs are correctly formatted.
	for _, a := range gs.DelegationsByStakerAssetOperator {
		stakerID := a.StakerID
		if _, _, err := assetstypes.ParseID(stakerID); err != nil {
			return errorsmod.Wrapf(ErrInvalidGenesisData, "stakerID invalid: %s", err)
		}
		for _, b := range a.DelegationsByAssetOperator {
			assetID := b.AssetID
			if _, _, err := assetstypes.ParseID(assetID); err != nil {
				return errorsmod.Wrapf(ErrInvalidGenesisData, "assetID invalid: %s", err)
			}
			for _, c := range b.DelegationsByOperator {
				operatorAddress := c.OperatorAddress
				_, err := sdk.AccAddressFromBech32(operatorAddress)
				if err != nil {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"operatorAddress invalid: %s",
						err,
					)
				}
				amount := c.Amount
				if amount.IsNil() {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"amount is nil for %s: %s: %s",
						stakerID,
						assetID,
						operatorAddress,
					)
				}
				if amount.IsNegative() {
					return errorsmod.Wrapf(
						ErrInvalidGenesisData,
						"amount is negative for %s: %s: %s",
						stakerID,
						assetID,
						operatorAddress,
					)
				}
			}
		}
	}
	return nil
}
