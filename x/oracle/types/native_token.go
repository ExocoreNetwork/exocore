package types

const maxSize = 100

// TODO: vlaidatorIndex need bridge data
func NewStakerInfo(stakerAddr string, validatorIndex uint64) *StakerInfo {
	return &StakerInfo{
		StakerAddr:  stakerAddr,
		StakerIndex: 0,
		// TODO: need bridge information
		ValidatorIndexes: []uint64{validatorIndex},
		BalanceList:      make([]*BalanceInfo, 0, 1),
	}
}

func (s *StakerInfo) Append(b *BalanceInfo) {
	s.BalanceList = append(s.BalanceList, b)
	if len(s.BalanceList) > maxSize {
		s.BalanceList = s.BalanceList[len(s.BalanceList)-maxSize:]
	}
}
