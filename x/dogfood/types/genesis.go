package types

import (
	errorsmod "cosmossdk.io/errors"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
)

// NewGenesis creates a new genesis state with the provided parameters and
// validators.
func NewGenesis(params Params, vals []GenesisValidator) *GenesisState {
	return &GenesisState{
		Params:        params,
		InitialValSet: vals,
	}
}

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	// no initial validators intentionally, so that the caller must set them.
	return NewGenesis(DefaultParams(), []GenesisValidator{})
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// #nosec G701 // ok on 64-bit systems.
	maxValidators := int(gs.Params.MaxValidators)
	if len(gs.InitialValSet) > maxValidators {
		return errorsmod.Wrapf(
			ErrInvalidGenesisData,
			"too many validators %d",
			len(gs.InitialValSet),
		)
	}
	// do not complain about 0 validators, let Tendermint do that.
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
		// and since its specific type (ed25519) is already set, it converts
		// easily to the sdk Key format as well.
		if _, err := operatortypes.HexStringToPubKey(
			val.PublicKey,
		); err != nil {
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
