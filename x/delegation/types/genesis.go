package types

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
	// TODO
	return nil
}
