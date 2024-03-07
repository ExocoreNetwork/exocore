package keeper

import (
	"errors"
	"math/big"
	"time"

	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: these consts should be defined in params
const (
	//maxNonce indicates how many messages a validator can submit in a single roudn to offer price
	//current we use this as a mock distance
	maxNonce = 3
	//these two threshold value used to set the threshold to tell when the price had come to consensus and was able to get a final price of that round
	threshold_a = 2
	threshold_b = 3
	//maxDetId each validator can submit, so the calculator can cache maximum of maxDetId*count(validators) values, this is for resistance of malicious validator submmiting invalid detId
	maxDetId = 5
	//consensus mode: v1: as soon as possbile
	mode = 1
)

var agc *aggregatorContext

type priceItemKV struct {
	tokenId int32
	priceTR types.PriceWithTimeAndRound
}

// worker is the actual instance used to calculate final price for each tokenFeeder's round. Which means, every tokenFeeder corresponds to a specified token, and for that tokenFeeder, each round we use a worker instance to calculate the final price
type worker struct {
	sealed  bool
	price   string
	decimal int32
	//mainly used for deterministic source data to check conflics and validation
	f *filter
	//used to get to consensus on deterministic source's data
	c *calculator
	//when enough data(exceeds threshold) collected, aggregate to conduct the final price
	a   *aggregator
	ctx *aggregatorContext
}

func (w *worker) do(msg *types.MsgCreatePrice) []*types.PriceWithSource {
	validator := msg.Creator
	power := w.ctx.validatorsPower[validator]
	list4Calculator, list4Aggregator := w.f.filtrate(msg)
	if list4Aggregator != nil {
		w.a.fillPrice(list4Aggregator, validator, power)
		if confirmedRounds := w.c.fillPrice(list4Calculator, validator, power); confirmedRounds != nil {
			w.a.confirmDSPrice(confirmedRounds)
		}
	}
	return list4Aggregator
}

func (w *worker) aggregate() *big.Int {
	return w.a.aggregate()
}

// not concurrency safe
func (w *worker) seal() {
	if w.sealed {
		return
	}
	w.sealed = true
	w.price = w.a.aggregate().String()
	w.f = nil
	w.c = nil
	w.a = nil
}

func (w *worker) getPrice() (string, int32) {
	if w.sealed {
		return w.price, w.decimal
	}
	return "", 0
}

type roundInfo struct {
	//this round of price will start from block basedBlock+1, the basedBlock served as a trigger to notify validators to submit prices
	basedBlock uint64
	//next round id of the price oracle service, price with thie id will be record on block basedBlock+1 if all prices submitted by validators(for v1, validators serve as oracle nodes) get to consensus immedately
	nextRoundId uint64
	//indicate if this round is open for collecting prices or closed in either condition that success with a consensused price or not
	//1: open, 2: closed
	status int32
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
}

func (agc *aggregatorContext) sanityCheck(msg *types.MsgCreatePrice) error {
	//sanity check
	//TODO: check nonce [1,3] in anteHandler, related to params, may not able
	//TODO: check the msgCreatePrice's Decimal is correct with params setting
	//TODO: check len(price.prices)>0, len(price.prices._range_eachPriceWithSource.Prices)>0, at least has one source, and for each source has at least one price
	//TODO: check for each source, at most maxDetId count price (now in filter, ->anteHandler)
	if agc.validatorsPower[msg.Creator] == nil {
		return errors.New("")
	}

	if msg.Nonce < 1 || msg.Nonce > maxNonce {
		return errors.New("")
	}

	//TODO: sanity check for price(no more than maxDetId count for each source, this should be take care in anteHandler)
	if msg.Prices == nil || len(msg.Prices) == 0 {
		return errors.New("")
	}

	for _, pSource := range msg.Prices {
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

func (agc *aggregatorContext) checkMsg(msg *types.MsgCreatePrice) error {
	if err := agc.sanityCheck(msg); err != nil {
		return err
	}

	//check feeder is active
	feederContext := agc.rounds[msg.FeederId]
	if feederContext == nil || feederContext.status != 1 {
		//feederId does not exist or not alive
		return errors.New("")
	}
	//senity check on basedBlock
	if msg.BasedBlock != feederContext.basedBlock {
		return errors.New("")
	}

	//check sources rule matches
	if ok, err := agc.params.checkRules(msg.FeederId, msg.Prices); !ok {
		return err
	}
	return nil
}

func (agc *aggregatorContext) fillPrice(msg *types.MsgCreatePrice) (*priceItemKV, *cacheItem, error) {
	feederWorker := agc.aggregators[msg.FeederId]
	//worker initialzed here reduce workload for Endblocker
	if feederWorker == nil {
		feederWorker = agc.newWorker(msg.FeederId)
		agc.aggregators[msg.FeederId] = feederWorker
	}

	if feederWorker.sealed {
		return nil, nil, errors.New("")
	}

	if listFilled := feederWorker.do(msg); listFilled != nil {
		if finalPrice := feederWorker.aggregate(); finalPrice != nil {
			agc.rounds[msg.FeederId].status = 2
			feederWorker.seal()
			return &priceItemKV{agc.params.getTokenFeeder(msg.FeederId).TokenId, types.PriceWithTimeAndRound{
				Price:   finalPrice.String(),
				Decimal: agc.params.getTokenInfo(msg.FeederId).Decimal,
				//TODO: check the format
				Timestamp: time.Now().String(),
				RoundId:   agc.rounds[msg.FeederId].nextRoundId,
			}}, nil, nil
		}
		return nil, &cacheItem{msg.FeederId, listFilled, msg.Creator}, nil
	}

	return nil, nil, errors.New("")
}

// NewCreatePrice receives msgCreatePrice message, and goes process: filter->aggregator, filter->calculator->aggregator
// non-deterministic data will goes directly into aggregator, and deterministic data will goes into calculator first to get consensus on the deterministic id.
func (agc *aggregatorContext) newCreatePrice(ctx sdk.Context, msg *types.MsgCreatePrice) (*priceItemKV, *cacheItem, error) {

	if err := agc.checkMsg(msg); err != nil {
		return nil, nil, err
	}

	return agc.fillPrice(msg)
}

// newWorker new a instance for a tokenFeeder's specific round
func (agc *aggregatorContext) newWorker(feederId int32) *worker {
	return &worker{
		f:       newFilter(maxNonce, maxDetId),
		c:       newCalculator(len(agc.validatorsPower), agc.totalPower),
		a:       newAggregator(len(agc.validatorsPower), agc.totalPower),
		decimal: agc.params.getTokenInfo(feederId).Decimal,
		ctx:     agc,
	}
}

// prepare for new roundInfo, just update the status kept in memory
// executed at EndBlock stage, seall all success or expired roundInfo
// including possible aggregation and state update
// returns: 1st successful sealed, need to be written to KVStore, 2nd: failed sealed tokenId, use previous price to write to KVStore
func (agc *aggregatorContext) sealRound(ctx sdk.Context) (success []*priceItemKV, failed []int32) {
	//1. check validatorSet udpate
	//TODO: if validatoSet has been updated in current block, just seal all active rounds and return
	//1. for sealed worker, the KVStore has been updated
	for feederId, round := range agc.rounds {
		if round.status == 1 {
			feeder := agc.params.getTokenFeeder(feederId)
			//TODO: for mode=1, we don't do aggregate() here, since if it donesn't success in the transaction execution stage, it won't success here
			//but it's not always the same for other modes, switch modes
			switch mode {
			case 1:
				expired := ctx.BlockHeight() >= feeder.EndBlock
				outOfWindow := uint64(ctx.BlockHeight())-round.basedBlock >= uint64(maxNonce)
				if expired || outOfWindow {
					//TODO: WRITE TO KVSTORE with previous round data for this round
					failed = append(failed, feeder.TokenId)
					if expired {
						delete(agc.rounds, feederId)
						delete(agc.aggregators, feederId)
					} else {
						round.status = 2
						agc.aggregators[feederId] = nil
						//TODO: WRITE TO KVSTORE with previous round data for this round
						failed = append(failed, feeder.TokenId)
					}
				}
			}
		}
		//all status: 1->2, remove its aggregator
		if agc.aggregators[feederId] != nil && agc.aggregators[feederId].sealed {
			agc.aggregators[feederId] = nil
		}
	}
	return
}

func (agC *aggregatorContext) prepareRound(ctx sdk.Context, block uint64) {
	//block>0 means recache initialization, all roundInfo is empty
	if block == 0 {
		block = uint64(ctx.BlockHeight())
	}

	for feederId, feeder := range agc.params.TokenFeeders {
		if uint64(feeder.EndBlock) <= block || uint64(feeder.StartBaseBlock) > block {
			//this feeder is inactive
			continue
		}

		delta := (block - uint64(feeder.StartBaseBlock))
		left := delta % uint64(feeder.Interval)
		count := delta / uint64(feeder.Interval)
		latestBasedblock := block - left
		latestNextRoundId := uint64(feeder.StartRoundId) + count

		feederIdInt32 := int32(feederId)
		round := agc.rounds[feederIdInt32]
		if round == nil {
			round = &roundInfo{
				basedBlock:  latestBasedblock,
				nextRoundId: latestNextRoundId,
			}
			if left >= maxNonce {
				round.status = 2
			} else {
				round.status = 1
			}
			agc.rounds[feederIdInt32] = round
		} else {
			//prepare a new round for exist roundInfo
			if left == 0 {
				round.basedBlock = latestBasedblock
				round.nextRoundId = latestNextRoundId
				round.status = 1
				//drop previous worker
				agc.aggregators[feederIdInt32] = nil
			}
		}
	}
}

func (agc *aggregatorContext) recache(from, to uint64, k Keeper) {
	//block_from'endblocker --> from_to'endblock
	//f.prepare->(f+1).msgs->(f+1).seal
	//...(to-2).prepare->(to-1).msgs->(to-1).seal
	//to.prepare-> return
	validatorsPower := k.GetValidators
}
