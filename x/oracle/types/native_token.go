package types

import sdkmath "cosmossdk.io/math"

const maxSize = 100

// TODO: vlaidatorIndex need bridge data
// func NewStakerInfo(stakerAddr string, validatorIndex uint64) *StakerInfo {
func NewStakerInfo(stakerAddr string, validatorIndex uint64) *StakerInfo {
	return &StakerInfo{
		StakerAddr:  stakerAddr,
		StakerIndex: 0,
		// TODO: need bridge information
		ValidatorIndexs: []uint64{validatorIndex},
		// ValidatorIndexs: make([]uint64, 0, 1),
		BalanceList: make([]*BalanceInfo, 0, 1),
	}
}

func NewOperatorInfo(operatorAddr string) *OperatorInfo {
	return &OperatorInfo{
		OperatorAddr: operatorAddr,
		TotalAmount:  sdkmath.NewInt(0),
		PriceList: []*PriceInfo{
			{
				// default price should be 1
				Price:   sdkmath.LegacyNewDec(1),
				Block:   0,
				RoundID: 0,
			},
		},
	}
}

func (s *StakerInfo) Append(b *BalanceInfo) {
	s.BalanceList = append(s.BalanceList, b)
	if len(s.BalanceList) > maxSize {
		s.BalanceList = s.BalanceList[len(s.BalanceList)-maxSize:]
	}
}
