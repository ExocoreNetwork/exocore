package types

import (
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
)

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return NewGenesis(commontypes.DefaultSubscriberParams())
}

// NewGenesis creates a new genesis state with the provided parameters and
// data. Since most of the genesis fields are filled by the coordinator,
// the subscriber module only needs to fill the subscriber params.
// Even those will be overwritten.
func NewGenesis(params commontypes.SubscriberParams) *GenesisState {
	return &GenesisState{Params: params}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return nil
}
