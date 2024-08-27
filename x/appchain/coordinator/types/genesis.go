package types

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return NewGenesis(DefaultParams())
}

// NewGenesis creates a new genesis state with the provided parameters and
// data.
func NewGenesis(params Params) *GenesisState {
	return &GenesisState{Params: params}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	// TODO: validate anything else added here
	return nil
}
