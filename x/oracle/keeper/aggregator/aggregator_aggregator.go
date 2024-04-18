package aggregator

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

type priceWithTimeAndRound struct {
	price      *big.Int
	decimal    int32
	timestamp  string
	detRoundId string //roundId from source if exists
}

type reportPrice struct {
	validator string
	//final price, set to -1 as initial
	price *big.Int
	//sourceId->priceWithTimeAndRound
	prices map[uint64]*priceWithTimeAndRound
	power  *big.Int
}

func (r *reportPrice) aggregate() *big.Int {
	if r.price != nil {
		return r.price
	}
	tmp := make([]*big.Int, 0, len(r.prices))
	for _, p := range r.prices {
		tmp = append(tmp, p.price)
	}
	r.price = common.BigIntList(tmp).Median()
	return r.price
}

type aggregator struct {
	finalPrice *big.Int
	reports    []*reportPrice
	//total valiadtor power who has submitted pice
	reportPower *big.Int
	totalPower  *big.Int
	//validator set total power
	//	totalPower string
	//sourceId->roundId used to track the confirmed DS roundId
	//updated by calculator, detId use string
	dsPrices map[uint64]string
}

// fill price from validator submitting into aggregator, and calculation the voting power and check with the consensus status of deterministic soruce value to decide when to do the aggregation
// TODO: currently apply mode=1 in V1, add swith modes
func (agg *aggregator) fillPrice(pSources []*types.PriceWithSource, validator string, power *big.Int) {
	report := agg.getReport(validator)
	if report == nil {
		report = &reportPrice{
			validator: validator,
			prices:    make(map[uint64]*priceWithTimeAndRound),
			power:     power,
		}
		agg.reports = append(agg.reports, report)
		agg.reportPower = new(big.Int).Add(agg.reportPower, power)
	}

	for _, pSource := range pSources {
		if len(pSource.Prices[0].DetId) == 0 {
			//this is an NS price report, price will just be updated instead of append
			if pTR := report.prices[pSource.SourceId]; pTR == nil {
				pTmp := pSource.Prices[0]
				priceBigInt, _ := (&big.Int{}).SetString(pTmp.Price, 10)
				pTR = &priceWithTimeAndRound{
					price:     priceBigInt,
					decimal:   pTmp.Decimal,
					timestamp: pTmp.Timestamp,
					//			detRoundId: p.DetId,
				}
				report.prices[pSource.SourceId] = pTR
			} else {
				pTR.price, _ = (&big.Int{}).SetString(pSource.Prices[0].Price, 10)
			}
		} else {
			//this is an DS price report
			if pTR := report.prices[pSource.SourceId]; pTR == nil {
				pTmp := pSource.Prices[0]
				pTR = &priceWithTimeAndRound{
					//price:     nil,
					decimal: pTmp.Decimal,
					//	timestamp: "",
					//detRoundId: "",
				}
				if len(agg.dsPrices[pSource.SourceId]) > 0 {
					for _, reportTmp := range agg.reports {
						if priceTmp := reportTmp.prices[pSource.SourceId]; priceTmp != nil && priceTmp.price != nil {
							pTR.price = new(big.Int).Set(priceTmp.price)
							pTR.detRoundId = priceTmp.detRoundId
							pTR.timestamp = priceTmp.timestamp
						}
					}
				}
				report.prices[pSource.SourceId] = pTR
			}
			//skip if this DS's slot exists, DS's value only updated by calculator
		}
	}
}

// TODO: for v1 use mode=1, which means agg.dsPrices with each key only be updated once, switch modes
func (agg *aggregator) confirmDSPrice(confirmedRounds []*confirmedPrice) {
	for _, priceSourceRound := range confirmedRounds {
		//update the latest round-detId for DS, TODO: in v1 we only update this value once since calculator will just ignore any further value once a detId has reached consensus
		//		agg.dsPrices[priceSourceRound.sourceId] = priceSourceRound.detId
		//this id's comparision need to format id to make sure them be the same length
		if id := agg.dsPrices[priceSourceRound.sourceId]; len(id) == 0 || (len(id) > 0 && id < priceSourceRound.detId) {
			agg.dsPrices[priceSourceRound.sourceId] = priceSourceRound.detId
			for _, report := range agg.reports {
				if report.price != nil {
					//price of IVA has completed
					continue
				}
				if price := report.prices[priceSourceRound.sourceId]; price != nil {
					price.detRoundId = priceSourceRound.detId
					price.timestamp = priceSourceRound.timestamp
					price.price = priceSourceRound.price
				} //else TODO: panice in V1
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

func (agg *aggregator) aggregate() *big.Int {
	if agg.finalPrice != nil {
		return agg.finalPrice
	}
	//TODO: implemetn different MODE for definition of consensus,
	//currently: use rule_1+MODE_1: {rule:specified source:`chainlink`, MODE: asap when power exceeds the threshold}
	//1. check OVA threshold
	//2. check IVA consensus with rule, TODO: for v1 we only implement with mode=1&rule=1
	if common.ExceedsThreshold(agg.reportPower, agg.totalPower) {
		//TODO: this is kind of a mock way to suite V1, need update to check with params.rule
		//check if IVA all reached consensus
		if len(agg.dsPrices) > 0 {
			validatorPrices := make([]*big.Int, 0, len(agg.reports))
			//do the aggregation to find out the 'final price'
			for _, validatorReport := range agg.reports {
				validatorPrices = append(validatorPrices, validatorReport.aggregate())
			}
			//vTmp := bigIntList(validatorPrices)
			agg.finalPrice = common.BigIntList(validatorPrices).Median()
			//clear relative aggregator for this feeder, all the aggregator,calculator, filter can be removed since this round has been sealed
		}
	}
	return agg.finalPrice
}

func newAggregator(validatorSetLength int, totalPower *big.Int) *aggregator {
	return &aggregator{
		reports:     make([]*reportPrice, 0, validatorSetLength),
		reportPower: big.NewInt(0),
		dsPrices:    make(map[uint64]string),
		totalPower:  totalPower,
	}
}
