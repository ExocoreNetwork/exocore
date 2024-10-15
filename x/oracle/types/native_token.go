package types

const maxSize = 100

func NewStakerInfo(stakerAddr, validatorPubkey string) *StakerInfo {
	return &StakerInfo{
		StakerAddr:          stakerAddr,
		StakerIndex:         0,
		ValidatorPubkeyList: []string{validatorPubkey},
		BalanceList:         make([]*BalanceInfo, 0, 1),
	}
}

func (s *StakerInfo) Append(b *BalanceInfo) {
	s.BalanceList = append(s.BalanceList, b)
	if len(s.BalanceList) > maxSize {
		s.BalanceList = s.BalanceList[len(s.BalanceList)-maxSize:]
	}
}
