package keeper

import (
	"errors"
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

type priceWithTimeAndRound struct {
	price      *big.Int
	decimal    int32
	timestamp  string
	detRoundId string //roundId from source if exists
}

type params types.Params

func (p *params) isDeterministicSource(sourceId int32) bool {
	return p.Sources[int(sourceId)].Deterministic
}

func (p *params) isValidSource(sourceId int32) bool {
	if sourceId == 0 {
		//custom defined source
		return true
	}
	return p.Sources[int(sourceId)].Valid
}

func (p *params) getTokenFeeder(feederId int32) *types.TokenFeeder {
	for k, v := range p.TokenFeeders {
		if int32(k) == feederId {
			return v
		}
	}
	return nil
}

func (p *params) checkRules(feederId int32, prices []*types.PriceWithSource) (bool, error) {
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

type set[T comparable] struct {
	size  int
	slice []T
}

func (s *set[T]) Add(value T) bool {
	if len(s.slice) == s.size {
		return false
	}
	for _, v := range s.slice {
		if v == value {
			return false
		}
	}
	s.slice = append(s.slice, value)
	s.size++
	return true
}

func (s *set[T]) Has(value T) bool {
	for _, v := range s.slice {
		if v == value {
			return true
		}
	}
	return false
}

func (s *set[T]) length() int {
	return s.size
}

func newSet[T comparable](length int) *set[T] {
	return &set[T]{
		size:  length,
		slice: make([]T, 0, length),
	}
}

func exceedsThreshold(power *big.Int, totalPower *big.Int) bool {
	return new(big.Int).Mul(power, big.NewInt(threshold_b)).Cmp(new(big.Int).Mul(totalPower, big.NewInt(threshold_a))) > 0
}
