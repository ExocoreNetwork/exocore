package aggregator

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

type confirmedPrice struct {
	sourceId  int32
	detId     string
	price     *big.Int
	timestamp string
}

// internal struct
type priceAndPower struct {
	price *big.Int
	power *big.Int
}

// for a specific DS round, it could have multiple values provided by different validators(should not be true if there's no malicious validator)
type roundPrices struct { //0 means NS
	detId     string
	prices    []*priceAndPower
	price     *big.Int
	timestamp string
	// confirmed bool
}

// udpate priceAndPower for a specific DSRoundID, if the price exists, increase its power with provided data
// return confirmed=true, when detect power exceeds the threshold
func (r *roundPrices) updatePriceAndPower(pw *priceAndPower, totalPower *big.Int) (updated bool, confirmed bool) {
	if r.price != nil {
		confirmed = true
		return
	}
	for _, item := range r.prices {
		if item.price.Cmp(pw.price) == 0 {
			item.power = new(big.Int).Add(item.power, pw.power)
			updated = true
			if common.ExceedsThreshold(item.power, totalPower) {
				r.price = item.price
				confirmed = true
			}
			return
		}
	}
	if len(r.prices) < cap(r.prices) {
		r.prices = append(r.prices, pw)
		updated = true
		if common.ExceedsThreshold(pw.power, totalPower) {
			r.price = pw.price
			//			r.confirmed = true
			confirmed = true
		}
	}
	return
}

// each DS corresponding a roundPriceList to represent its multiple rounds(DS round) in one oracle-round
type roundPricesList struct {
	roundPricesList []*roundPrices
	//each round can have at most roundPricesCount priceAndPower
	roundPricesCount int
}

// to tell if any round of this DS has reached consensus/confirmed
func (r *roundPricesList) hasConfirmedDetId() bool {
	for _, round := range r.roundPricesList {
		if round.price != nil {
			return true
		}
	}
	return false
}

// get the roundPriceList correspond to specifid detID of a DS
// if no required data and the pricesList has not reach its limitation, we will add a new slot for this detId
func (r *roundPricesList) getOrNewRound(detId string, timestamp string) (round *roundPrices) {
	for _, round = range r.roundPricesList {
		if round.detId == detId {
			if round.price != nil {
				round = nil
			}
			return
		}
	}

	if len(r.roundPricesList) < cap(r.roundPricesList) {
		round = &roundPrices{
			detId:     detId,
			prices:    make([]*priceAndPower, 0, r.roundPricesCount),
			timestamp: timestamp,
		}
		r.roundPricesList = append(r.roundPricesList, round)
		return
	}
	return
}

// calculator used to get consensus on deterministic source based data from validator set reports of price
type calculator struct {
	//sourceId->{[]{roundId, []{price,power}, confirmed}}, confirmed value will be set in [0]
	deterministicSource map[int32]*roundPricesList
	validatorLength     int
	totalPower          *big.Int
}

func (c *calculator) newRoundPricesList() *roundPricesList {
	return &roundPricesList{
		roundPricesList: make([]*roundPrices, 0, common.MaxDetId*c.validatorLength),
		//for each DS-roundId, the count of prices provided is the number of validators at most
		roundPricesCount: c.validatorLength,
	}
}

func (c *calculator) getOrNewSourceId(sourceId int32) *roundPricesList {
	rounds := c.deterministicSource[sourceId]
	if rounds == nil {
		rounds = c.newRoundPricesList()
		c.deterministicSource[sourceId] = rounds
	}
	return rounds
}

// fillPrice called upon new MsgCreatPrice arrived, to trigger the calculation to get to consensus on the same roundID_of_deterministic_source
// v1 use mode1, TODO: switch modes
func (c *calculator) fillPrice(pSources []*types.PriceWithSource, validator string, power *big.Int) (confirmedRounds []*confirmedPrice) {
	//	fmt.Println("debug calculator.fillPrice, calculator.ds[1]", c.deterministicSource[1], pSources)
	for _, pSource := range pSources {
		rounds := c.getOrNewSourceId(pSource.SourceId)
		if rounds.hasConfirmedDetId() {
			//TODO: this skip is just for V1 to do fast calculation and release EndBlocker pressure, may lead to 'not latest detId' be chosen
			break
		}
		for _, pDetId := range pSource.Prices {

			round := rounds.getOrNewRound(pDetId.DetId, pDetId.Timestamp)
			if round == nil {
				//this sourceId has reach the limitation of different detId, or has confirmed
				continue
			}

			roundPrice, _ := new(big.Int).SetString(pDetId.Price, 10)

			//			fmt.Printf("debug calculator.fillPrice before updatePriceAndPower. power%s, price%s\n", power.String(), roundPrice.String())
			updated, confirmed := round.updatePriceAndPower(&priceAndPower{roundPrice, power}, c.totalPower)
			if updated && confirmed {
				//sourceId, detId, price
				confirmedRounds = append(confirmedRounds, &confirmedPrice{pSource.SourceId, round.detId, round.price, round.timestamp}) //TODO: just in v1 with mode==1, we use asap, so we just ignore any further data from this DS, even higher detId may get to consensus, in this way, in most case, we can complete the calculation in the transaction execution process. Release the pressure in EndBlocker
				//TODO: this may delay to current block finish
				break
			}
		}
	}
	//	fmt.Println("debug calculator.fillPrice, after calculator.ds[1]", c.deterministicSource[1])
	return
}

func newCalculator(validatorSetLength int, totalPower *big.Int) *calculator {
	return &calculator{
		deterministicSource: make(map[int32]*roundPricesList),
		validatorLength:     validatorSetLength,
		totalPower:          totalPower,
	}
}
