package keeper

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

type reportPrice struct {
	validator string
	//final price, set to -1 as initial
	price *big.Int
	//sourceId->priceWithTimeAndRound
	prices map[int32]*priceWithTimeAndRound
	power  *big.Int
}

type aggregator struct {
	reports []*reportPrice
	//total valiadtor power who has submitted pice
	reportPower *big.Int
	totalPower  *big.Int
	//validator set total power
	//	totalPower string
	//sourceId->roundId used to track the confirmed DS roundId
	//updated by calculator, detId use string
	dsPrices map[int32]string
}

// fill price from validator submittion into aggregator, and calculation the voting power and check with the consensus status of deterministic soruce value to decide when to do the aggregation
func (agg *aggregator) fillPrice(prices []*types.PriceWithSource, validator string, power *big.Int) {
	report := agg.getReport(validator)
	if report == nil {
		report = &reportPrice{
			validator: validator,
			prices:    make(map[int32]*priceWithTimeAndRound),
			power:     power,
		}
		agg.reports = append(agg.reports, report)
		agg.reportPower = new(big.Int).Add(agg.totalPower, power)
	}

	for _, p := range prices {
		if len(p.Prices[0].DetId) == 0 {
			//this is an NS price report, price will just be updated instead of append
			if pTR := report.prices[p.SourceId]; pTR == nil {
				pTmp := p.Prices[0]
				priceBigInt, _ := (&big.Int{}).SetString(pTmp.Price, 10)
				pTR = &priceWithTimeAndRound{
					price:     priceBigInt,
					decimal:   pTmp.Decimal,
					timestamp: pTmp.Timestamp,
					//			detRoundId: p.DetId,
				}
				report.prices[p.SourceId] = pTR
			} else {
				pTR.price, _ = (&big.Int{}).SetString(p.Prices[0].Price, 10)
			}
		} else {
			//this is an DS price report
			if pTR := report.prices[p.SourceId]; pTR == nil {
				pTmp := p.Prices[0]
				pTR = &priceWithTimeAndRound{
					//price:     nil,
					decimal:   pTmp.Decimal,
					timestamp: "",
					//detRoundId: "",
				}
				report.prices[p.SourceId] = pTR
			}
		}
	}
}

func (agg *aggregator) getReport(validator string) *reportPrice {
	for _, r := range agg.reports {
		if r.validator == validator {
			return r
		}
	}
	return nil
}

func (agg *aggregator) aggregate() {
	//TODO: implemetn different MODE for definition of consensus,
	//currently: use rule_1+MODE_1: {rule:specified source:`chainlink`, MODE: asap when power exceeds the threshold}
}

func newAggregator(validatorSetLength int, totalPower *big.Int) *aggregator {
	return &aggregator{
		reports:     make([]*reportPrice, validatorSetLength),
		reportPower: big.NewInt(0),
		dsPrices:    make(map[int32]string),
		totalPower:  totalPower,
	}
}
