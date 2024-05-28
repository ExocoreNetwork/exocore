package common

import (
	"errors"
	"math/big"
	"sort"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

const (
	// maxNonce indicates how many messages a validator can submit in a single round to offer price
	// current we use this as a mock distance
	MaxNonce = 3

	// these two threshold value used to set the threshold to tell when the price had come to consensus and was able to get a final price of that round
	ThresholdA = 2
	ThresholdB = 3

	// maxDetId each validator can submit, so the calculator can cache maximum of maxDetId*count(validators) values, this is for resistance of malicious validator submmiting invalid detId
	MaxDetID = 5

	// consensus mode: v1: as soon as possbile
	Mode = 1
)

type Params types.Params

func (p Params) GetTokenFeeders() []*types.TokenFeeder {
	return p.TokenFeeders
}

func (p Params) IsDeterministicSource(sourceID uint64) bool {
	return p.Sources[sourceID].Deterministic
}

func (p Params) IsValidSource(sourceID uint64) bool {
	if sourceID == 0 {
		// custom defined source
		return true
	}
	return p.Sources[sourceID].Valid
}

func (p Params) GetTokenFeeder(feederID uint64) *types.TokenFeeder {
	for k, v := range p.TokenFeeders {
		if uint64(k) == feederID {
			return v
		}
	}
	return nil
}

func (p Params) GetTokenInfo(feederID uint64) *types.Token {
	for k, v := range p.TokenFeeders {
		if uint64(k) == feederID {
			return p.Tokens[v.TokenID]
		}
	}
	return nil
}

func (p Params) CheckRules(feederID uint64, prices []*types.PriceSource) (bool, error) {
	feeder := p.TokenFeeders[feederID]
	rule := p.Rules[feeder.RuleID]
	// specified sources set, v1 use this rule to set `chainlink` as official source
	if rule.SourceIDs != nil && len(rule.SourceIDs) > 0 {
		if len(rule.SourceIDs) != len(prices) {
			return false, errors.New("count prices should match rule")
		}
		notFound := false
		if rule.SourceIDs[0] == 0 {
			// match all sources listed
			for sID, source := range p.Sources {
				if sID == 0 {
					continue
				}
				if source.Valid {
					notFound = true
					for _, p := range prices {
						if p.SourceID == uint64(sID) {
							notFound = false
							break
						}
					}

				}
			}
		} else {
			for _, source := range rule.SourceIDs {
				notFound = true
				for _, p := range prices {
					if p.SourceID == source {
						notFound = false
						break
					}
				}
				// return false, errors.New("price source not match with rule")
			}
		}
		if notFound {
			return false, errors.New("price source not match with rule")
		}
	}

	// TODO: check NOM
	// return true if no rule set, we will accept any source
	return true, nil
}

type Set[T comparable] struct {
	size  int
	slice []T
}

func (s *Set[T]) Copy() *Set[T] {
	ret := NewSet[T](s.Length())
	ret.slice = append(ret.slice, s.slice...)
	return ret
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
	return new(big.Int).Mul(power, big.NewInt(ThresholdB)).Cmp(new(big.Int).Mul(totalPower, big.NewInt(ThresholdA))) > 0
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
