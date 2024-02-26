package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PricesList:    []Prices{},
		RoundInfoList: []RoundInfo{},
		RoundDataList: []RoundData{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated index in prices
	pricesIndexMap := make(map[string]struct{})

	for _, elem := range gs.PricesList {
		index := string(PricesKey(elem.TokenId))
		if _, ok := pricesIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for prices")
		}
		pricesIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in roundInfo
	roundInfoIndexMap := make(map[string]struct{})

	for _, elem := range gs.RoundInfoList {
		index := string(RoundInfoKey(elem.TokenId))
		if _, ok := roundInfoIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for roundInfo")
		}
		roundInfoIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in roundData
	roundDataIndexMap := make(map[string]struct{})

	for _, elem := range gs.RoundDataList {
		index := string(RoundDataKey(elem.TokenId))
		if _, ok := roundDataIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for roundData")
		}
		roundDataIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
