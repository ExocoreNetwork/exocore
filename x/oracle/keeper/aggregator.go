package keeper

import (
	"errors"
	"math/big"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
)

// TODO: these consts should be defined in params
const (
	//maxNonce indicates how many messages a validator can submit in a single roudn to offer price
	maxNonce = 3
	//these two threshold value used to set the threshold to tell when the price had come to consensus and was able to get a final price of that round
	threshold_a = 2
	threshold_b = 3
	//maxDetId each validator can submit, so the calculator can cache maximum of maxDetId*count(validators) values, this is for resistance of malicious validator submmiting invalid detId
	maxDetId = 5
	//consensus mode: v1: as soon as possbile
	mode = 1
)

type roundInfo struct {
	//this round of price will start from block basedBlock+1, the basedBlock served as a trigger to notify validators to submit prices
	basedBlock uint64
	//next round id of the price oracle service, price with thie id will be record on block basedBlock+1 if all prices submitted by validators(for v1, validators serve as oracle nodes) get to consensus immedately
	nextRoundId uint64
	//indicate if this round is open for collecting prices or closed in either condition that success with a consensused price or not
	//1: open, 2: closed
	status int32
}

// worker is the actual instance used to calculate final price for each tokenFeeder's round. Which means, every tokenFeeder corresponds to a specified token, and for that tokenFeeder, each round we use a worker instance to calculate the final price
type worker struct {
	//mainly used for deterministic source data to check conflics and validation
	f *filter
	//used to get to consensus on deterministic source's data
	c *calculator
	//when enough data(exceeds threshold) collected, aggregate to conduct the final price
	a *aggregator
}

// aggregatorContext keeps memory cache for state params, validatorset, and updatedthese values as they udpated on chain. And it keeps the infomation to track all tokenFeeders' status and data collection
type aggregatorContext struct {
	params *params

	//validator->power
	validatorsPower map[string]*big.Int
	totalPower      *big.Int

	//each active feederToken has a roundInfo
	rounds map[int32]*roundInfo

	//each roundInfo has a worker
	aggregators map[int32]*worker

	k *Keeper
}

//workerInstance chan *type.MsgCreatePrice

func (agc *aggregatorContext) sanityCheck(price types.MsgCreatePrice) error {
	//sanity check
	//TODO: check nonce [1,3] in anteHandler, related to params, may not able
	//TODO: check the msgCreatePrice's Decimal is correct with params setting
	//TODO: check len(price.prices)>0, len(price.prices._range_eachPriceWithSource.Prices)>0, at least has one source, and for each source has at least one price
	//TODO: check for each source, at most maxDetId count price (now in filter, ->anteHandler)
	if price.Nonce < 1 || price.Nonce > maxNonce {
		return errors.New("")
	}

	//TODO: sanity check for price(no more than maxDetId count for each source, this should be take care in anteHandler)
	if price.Prices == nil || len(price.Prices) == 0 {
		return errors.New("")
	}

	for _, pSource := range price.Prices {
		if pSource.Prices == nil || len(pSource.Prices) == 0 || len(pSource.Prices) > maxDetId || !agc.params.isValidSource(pSource.SourceId) {
			return errors.New("")
		}
		//check with params is coressponding source is deteministic
		if agc.params.isDeterministicSource(pSource.SourceId) {
			for _, pDetId := range pSource.Prices {
				//TODO: verify the format of DetId is correct, since this is string, and we will make consensus with validator's power, so it's ok not to verify the format
				//just make sure the DetId won't mess up with NS's placeholder id, the limitation of maximum count one validator can submit will be check by filter
				if len(pDetId.DetId) == 0 {
					//deterministic must have specified deterministicId
					return errors.New("")
				}
				//DS's price value will go through consensus process, so it's safe to skip the check here
			}
		} else {
			//sanity check: NS submit only one price with detId==""
			if len(pSource.Prices) > 1 || len(pSource.Prices[0].DetId) > 0 {
				return errors.New("")
			}
		}
	}

	return nil
}

// NewCreatePrice receives msgCreatePrice message, and goes process: filter->aggregator, filter->calculator->aggregator
// non-deterministic data will goes directly into aggregator, and deterministic data will goes into calculator first to get consensus on the deterministic id.
func (agc *aggregatorContext) NewCreatePrice(price types.MsgCreatePrice) (bool, error) {

	if err := agc.sanityCheck(price); err != nil {
		return false, err
	}

	validator := price.Creator
	//validator exists in current active validator set
	if power := agc.validatorsPower[validator]; power != nil {
		//check feeder is active
		feederContext := agc.rounds[price.FeederId]
		if feederContext == nil || feederContext.status != 1 {
			//feederId does not exist or not alive
			return false, errors.New("")
		}
		//senity check on basedBlock
		if price.BasedBlock != feederContext.basedBlock {
			return false, errors.New("")
		}

		//check sources rule matches
		if ok, err := agc.params.checkRules(price.FeederId, price.Prices); !ok {
			return false, err
		}

		feederWorker := agc.aggregators[price.FeederId]
		//worker initialzed here reduce workload for Endblocker
		if feederWorker == nil {
			feederWorker = agc.newWorker()
			agc.aggregators[price.FeederId] = feederWorker
		}

		list4Calculator, list4Aggregator := feederWorker.f.filtrate(price)

		feederWorker.a.fillPrice(list4Aggregator, validator, power)
		if confirmedRounds := feederWorker.c.fillPrice(list4Calculator, validator, power); confirmedRounds != nil {
			feederWorker.a.confirmDSPrice(confirmedRounds)
		}

		agc.k.addCache(&cacheItem{list4Aggregator, validator, power})
	}

	//invalid creator, require validator to be the price reporter
	return false, errors.New("")
}

// newWorker new a instance for a tokenFeeder's specific round
func (aggC *aggregatorContext) newWorker() *worker {
	return &worker{
		f: newFilter(maxNonce, maxDetId),
		c: newCalculator(len(aggC.validatorsPower), aggC.totalPower),
		a: newAggregator(len(aggC.validatorsPower), aggC.totalPower),
	}
}
