package keeper

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

// internal struct
type priceAndPower struct {
	price *big.Int
	power *big.Int
}

type roundPrices struct { //0 means NS
	detId     string
	prices    []*priceAndPower
	confirmed bool
}

type roundPricesList []*roundPrices

func (r roundPricesList) getRound(detId string) *roundPrices {
	for _, round := range r {
		if round.detId == detId {
			return round
		}
	}
	return nil
}

// calculator used to get consensus on deterministic source based data from validator set reports of price
type calculator struct {
	//sourceId->{[]{roundId, []{price,power}, confirmed}}, confirmed value will be set in [0]
	deterministicSource map[int32]roundPricesList
	validatorLength     int
}

// fillPrice called upon new MsgCreatPrice arrived, to trigger the calculation to get to consensus on the same roundID_of_deterministic_source
func (c *calculator) fillPrice(prices []*types.PriceWithSource, validator string, power *big.Int) {
	for _, pSource := range prices {
		rounds := c.deterministicSource[pSource.SourceId]
		if rounds == nil {
			rounds = make([]*roundPrices, 0, maxDetId*c.validatorLength)
			c.deterministicSource[pSource.SourceId] = rounds
		}

		for _, pDetId := range pSource.Prices {

			round := rounds.getRound(pDetId.DetId)
			if round == nil {
				if len(rounds) < cap(rounds) {
					//add a new roundId from source
					round = &roundPrices{
						detId:  pDetId.DetId,
						prices: make([]*priceAndPower, 0, c.validatorLength),
						//confirmed: false,
					}
					roundPrice, _ := new(big.Int).SetString(pDetId.Price, 10)
					round.prices = append(round.prices, &priceAndPower{roundPrice, power})
					//TODO: check if power exceeds the threshold, which means single validator has most of the voting power, bad.
					rounds = append(rounds, round)
					c.deterministicSource[pSource.SourceId] = rounds
				}
				//ignore this source price
				//only accept maxDetId count different roundId
			} else {
				//TODO: do the calculation and trigger the aggregator update if any value's power exceeds the threshold

			}
		}
	}
}

func newCalculator(validatorSetLength int) *calculator {
	return &calculator{
		deterministicSource: make(map[int32]roundPricesList),
		validatorLength:     validatorSetLength,
	}
}
