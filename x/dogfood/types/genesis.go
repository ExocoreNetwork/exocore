package types

import (
	errorsmod "cosmossdk.io/errors"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// #nosec G701 // ok on 64-bit systems.
	maxValidators := int(gs.Params.MaxValidators)
	if len(gs.InitialValSet) == 0 || len(gs.InitialValSet) > maxValidators {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"invalid number of validators %d",
			len(gs.InitialValSet),
		)
	}
	for _, val := range gs.InitialValSet {
		if _, err := operatortypes.Bytes32ToPubKey(val.PublicKey); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid public key %x: %s",
				val.PublicKey, err,
			)
		}
		power := val.Power
		if power <= 0 {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid power %d",
				power,
			)
		}
	}
	return gs.Params.Validate()
}
