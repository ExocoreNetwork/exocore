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
	vals := make(map[string]struct{}, len(gs.InitialValSet))
	for _, val := range gs.InitialValSet {
		// check for duplicates
		if _, ok := vals[val.PublicKey]; ok {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"duplicate public key %s", val.PublicKey,
			)
		}
		vals[val.PublicKey] = struct{}{}
		// HexStringToPubKey checks the size and returns a tmprotocrypto type.
		// and since its specific type (ed25519) is already set, it converts easily to
		// sdk Key format as well.
		if _, err := operatortypes.HexStringToPubKey(val.PublicKey); err != nil {
			return errorsmod.Wrapf(
				ErrInvalidGenesisData,
				"invalid public key %s: %s",
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
