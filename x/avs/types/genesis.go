package types

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

func NewGenesisState(
	avsInfos []AVSInfo,
	taskInfos []TaskInfo,
	blsPubKeys []BlsPubKeyInfo,
	taskResultInfos []TaskResultInfo,
	challengeInfos []ChallengeInfo,
	taskNums []TaskID,
	chainIDInfos []ChainIDInfo,
) *GenesisState {
	return &GenesisState{
		AvsInfos:        avsInfos,
		TaskInfos:       taskInfos,
		BlsPubKeys:      blsPubKeys,
		TaskResultInfos: taskResultInfos,
		ChallengeInfos:  challengeInfos,
		TaskNums:        taskNums,
		ChainIdInfos:    chainIDInfos,
	}
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(nil, nil, nil, nil, nil, nil, nil)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	return nil
}
