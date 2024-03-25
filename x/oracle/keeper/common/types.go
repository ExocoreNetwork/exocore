package common

import (
	"errors"
	"math/big"
	"sort"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

const (
	//maxNonce indicates how many messages a validator can submit in a single roudn to offer price
	//current we use this as a mock distance
	MaxNonce = 3
	//these two threshold value used to set the threshold to tell when the price had come to consensus and was able to get a final price of that round
	Threshold_a = 2
	Threshold_b = 3
	//maxDetId each validator can submit, so the calculator can cache maximum of maxDetId*count(validators) values, this is for resistance of malicious validator submmiting invalid detId
	MaxDetId = 5
	//consensus mode: v1: as soon as possbile
	Mode = 1
)

type Params types.Params

func (p Params) GetTokenFeeders() []*types.TokenFeeder {
	return p.TokenFeeders
}

func (p Params) IsDeterministicSource(sourceId int32) bool {
	return p.Sources[int(sourceId)].Deterministic
}

func (p Params) IsValidSource(sourceId int32) bool {
	if sourceId == 0 {
		//custom defined source
		return true
	}
	return p.Sources[int(sourceId)].Valid
}

func (p Params) GetTokenFeeder(feederId int32) *types.TokenFeeder {
	for k, v := range p.TokenFeeders {
		if int32(k) == feederId {
			return v
		}
	}
	return nil
}
func (p Params) GetTokenInfo(feederId int32) *types.Token {
	for k, v := range p.TokenFeeders {
		if int32(k) == feederId {
			return p.Tokens[v.TokenId]
		}
	}
	return nil
}

func (p Params) CheckRules(feederId int32, prices []*types.PriceWithSource) (bool, error) {
	feeder := p.TokenFeeders[feederId]
	rule := p.Rules[feeder.RuleId]
	//specified sources set, v1 use this rule to set `chainlink` as official source
	if rule.SourceIds != nil && len(rule.SourceIds) > 0 {
		if len(rule.SourceIds) != len(prices) {
			return false, errors.New("")
		}
		for _, source := range rule.SourceIds {
			for _, p := range prices {
				if p.SourceId == source {
					continue
				}
			}
			return false, errors.New("")
		}
	}

	//TODO: check NOM
	//return true if no rule set, we will accept any source
	return true, nil
}

type Set[T comparable] struct {
	size  int
	slice []T
}

func (s *Set[T]) Add(value T) bool {
	if len(s.slice) == s.size {
		return false
	}
	for _, v := range s.slice {
		if v == value {
			return false
		}
	}
	s.slice = append(s.slice, value)
	return true
}

func (s *Set[T]) Has(value T) bool {
	for _, v := range s.slice {
		if v == value {
			return true
		}
	}
	return false
}

func (s *Set[T]) Length() int {
	return s.size
}

func NewSet[T comparable](length int) *Set[T] {
	return &Set[T]{
		size:  length,
		slice: make([]T, 0, length),
	}
}

func ExceedsThreshold(power *big.Int, totalPower *big.Int) bool {
	return new(big.Int).Mul(power, big.NewInt(Threshold_b)).Cmp(new(big.Int).Mul(totalPower, big.NewInt(Threshold_a))) > 0
}

type BigIntList []*big.Int

func (b BigIntList) Len() int {
	return len(b)
}
func (b BigIntList) Less(i, j int) bool {
	return b[i].Cmp(b[j]) < 0
}
func (b BigIntList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BigIntList) Median() *big.Int {
	sort.Sort(b)
	l := len(b)
	if l%2 == 1 {
		return b[l/2]
	}
	return new(big.Int).Div(new(big.Int).Add(b[l/2], b[l/2-1]), big.NewInt(2))
}
