package types

func NewGenesisState(
	operators []OperatorInfo,
	records []OperatorConsKeyRecord,
) *GenesisState {
	return &GenesisState{
		Operators:       operators,
		OperatorRecords: records,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState([]OperatorInfo{}, []OperatorConsKeyRecord{})
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// TODO
	return nil
}
