package aggregator

import (
	"errors"
	"math/big"
	"time"

	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/cache"
	"github.com/ExocoreNetwork/exocore/x/oracle/keeper/common"
	"github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type priceItemKV struct {
	TokenId int32
	PriceTR types.PriceWithTimeAndRound
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

// AggregatorContext keeps memory cache for state params, validatorset, and updatedthese values as they udpated on chain. And it keeps the infomation to track all tokenFeeders' status and data collection
type AggregatorContext struct {
	params *common.Params

	//validator->power
	validatorsPower map[string]*big.Int
	totalPower      *big.Int

	//each active feederToken has a roundInfo
	rounds map[int32]*roundInfo

	//each roundInfo has a worker
	aggregators map[int32]*worker
}

func (agc *AggregatorContext) sanityCheck(msg *types.MsgCreatePrice) error {
	//sanity check
	//TODO: check nonce [1,3] in anteHandler, related to params, may not able
	//TODO: check the msgCreatePrice's Decimal is correct with params setting
	//TODO: check len(price.prices)>0, len(price.prices._range_eachPriceWithSource.Prices)>0, at least has one source, and for each source has at least one price
	//TODO: check for each source, at most maxDetId count price (now in filter, ->anteHandler)

	if agc.validatorsPower[msg.Creator] == nil {
		return errors.New("signer is not validator")
	}

	if msg.Nonce < 1 || msg.Nonce > common.MaxNonce {
		return errors.New("nonce invalid")
	}

	//TODO: sanity check for price(no more than maxDetId count for each source, this should be take care in anteHandler)
	if msg.Prices == nil || len(msg.Prices) == 0 {
		return errors.New("msg should provide at least one price")
	}

	for _, pSource := range msg.Prices {
		if pSource.Prices == nil || len(pSource.Prices) == 0 || len(pSource.Prices) > common.MaxDetId || !agc.params.IsValidSource(pSource.SourceId) {
			return errors.New("source should be valid and provide at least one price")
		}
		//check with params is coressponding source is deteministic
		if agc.params.IsDeterministicSource(pSource.SourceId) {
			for _, pDetID := range pSource.Prices {
				//TODO: verify the format of DetId is correct, since this is string, and we will make consensus with validator's power, so it's ok not to verify the format
				//just make sure the DetId won't mess up with NS's placeholder id, the limitation of maximum count one validator can submit will be check by filter
				if len(pDetID.DetId) == 0 {
					//deterministic must have specified deterministicId
					return errors.New("ds should have roundid")
				}
				//DS's price value will go through consensus process, so it's safe to skip the check here
			}
			//sanity check: NS submit only one price with detId==""
		} else if len(pSource.Prices) > 1 || len(pSource.Prices[0].DetId) > 0 {
			return errors.New("ns should not have roundid")
		}
	}
	return nil
}

func (agc *AggregatorContext) checkMsg(msg *types.MsgCreatePrice) error {
	if err := agc.sanityCheck(msg); err != nil {
		return err
	}

	//check feeder is active
	feederContext := agc.rounds[msg.FeederId]
	if feederContext == nil || feederContext.status != 1 {
		//feederId does not exist or not alive
		return errors.New("context not exist or not available")
	}
	//senity check on basedBlock
	if msg.BasedBlock != feederContext.basedBlock {
		return errors.New("baseblock not match")
	}

	//check sources rule matches
	if ok, err := agc.params.CheckRules(msg.FeederId, msg.Prices); !ok {
		return err
	}
	return nil
}

func (agc *AggregatorContext) FillPrice(msg *types.MsgCreatePrice) (*priceItemKV, *cache.CacheItemM, error) {
	feederWorker := agc.aggregators[msg.FeederId]
	//worker initialzed here reduce workload for Endblocker
	if feederWorker == nil {
		feederWorker = newWorker(msg.FeederId, agc)
		agc.aggregators[msg.FeederId] = feederWorker
	}

	if feederWorker.sealed {
		return nil, nil, types.ErrPriceProposalIgnored.Wrap("price aggregation for this round has sealed")
	}

	if listFilled := feederWorker.do(msg); listFilled != nil {
		if finalPrice := feederWorker.aggregate(); finalPrice != nil {
			agc.rounds[msg.FeederId].status = 2
			feederWorker.seal()
			return &priceItemKV{agc.params.GetTokenFeeder(msg.FeederId).TokenId, types.PriceWithTimeAndRound{
				Price:   finalPrice.String(),
				Decimal: agc.params.GetTokenInfo(msg.FeederId).Decimal,
				//TODO: check the format
				Timestamp: time.Now().String(),
				RoundId:   agc.rounds[msg.FeederId].nextRoundId,
			}}, &cache.CacheItemM{FeederId: msg.FeederId}, nil
		}
		return nil, &cache.CacheItemM{msg.FeederId, listFilled, msg.Creator}, nil
	}

	//return nil, nil, errors.New("no valid price proposal to add for aggregation")
	return nil, nil, types.ErrPriceProposalIgnored
}

// NewCreatePrice receives msgCreatePrice message, and goes process: filter->aggregator, filter->calculator->aggregator
// non-deterministic data will goes directly into aggregator, and deterministic data will goes into calculator first to get consensus on the deterministic id.
func (agc *AggregatorContext) NewCreatePrice(ctx sdk.Context, msg *types.MsgCreatePrice) (*priceItemKV, *cache.CacheItemM, error) {
	if err := agc.checkMsg(msg); err != nil {
		return nil, nil, types.ErrInvalidMsg.Wrap(err.Error())
	}
	return agc.FillPrice(msg)
}

// prepare for new roundInfo, just update the status kept in memory
// executed at EndBlock stage, seall all success or expired roundInfo
// including possible aggregation and state update
// when validatorSet update, set force to true, to seal all alive round
// returns: 1st successful sealed, need to be written to KVStore, 2nd: failed sealed tokenId, use previous price to write to KVStore
func (agc *AggregatorContext) SealRound(ctx sdk.Context, force bool) (success []*priceItemKV, failed []int32) {
	//1. check validatorSet udpate
	//TODO: if validatoSet has been updated in current block, just seal all active rounds and return
	//1. for sealed worker, the KVStore has been updated
	for feederId, round := range agc.rounds {
		if round.status == 1 {
			feeder := agc.params.GetTokenFeeder(feederId)
			//TODO: for mode=1, we don't do aggregate() here, since if it donesn't success in the transaction execution stage, it won't success here
			//but it's not always the same for other modes, switch modes
			switch common.Mode {
			case 1:
				expired := feeder.EndBlock > 0 && ctx.BlockHeight() >= feeder.EndBlock
				outOfWindow := uint64(ctx.BlockHeight())-round.basedBlock >= uint64(common.MaxNonce)
				if expired || outOfWindow || force {
					failed = append(failed, feeder.TokenId)
					if expired {
						delete(agc.rounds, feederId)
						delete(agc.aggregators, feederId)
					} else {
						round.status = 2
						agc.aggregators[feederId] = nil
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

//func (agc *AggregatorContext) ForceSeal(ctx sdk.Context) (success []*priceItemKV, failed []int32) {
//
//}

func (agc *AggregatorContext) PrepareRound(ctx sdk.Context, block uint64) {
	//block>0 means recache initialization, all roundInfo is empty
	if block == 0 {
		block = uint64(ctx.BlockHeight())
	}

	for feederId, feeder := range agc.params.GetTokenFeeders() {
		if feederId == 0 {
			continue
		}
		if (feeder.EndBlock > 0 && uint64(feeder.EndBlock) <= block) || uint64(feeder.StartBaseBlock) > block {

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
			if left >= common.MaxNonce {
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
			} else if round.status == 1 && left >= common.MaxNonce {
				//this shouldn't happend, if do sealround properly before prepareRound, basically for test only
				round.status = 2
				//TODO: just modify the status here, since sealRound should do all the related seal actios already when parepare invoked
			}
		}
	}
}

func (agc *AggregatorContext) SetParams(p *common.Params) {
	agc.params = p
}

func (agc *AggregatorContext) SetValidatorPowers(vp map[string]*big.Int) {
	//	t := big.NewInt(0)
	agc.totalPower = big.NewInt(0)
	agc.validatorsPower = make(map[string]*big.Int)
	for addr, power := range vp {
		agc.validatorsPower[addr] = power
		agc.totalPower = new(big.Int).Add(agc.totalPower, power)
	}
}
func (agc *AggregatorContext) GetValidatorPowers() (vp map[string]*big.Int) {
	return agc.validatorsPower
}

//func (agc *AggregatorContext) SetTotalPower(power *big.Int) {
//	agc.totalPower = power
//}

func NewAggregatorContext() *AggregatorContext {
	return &AggregatorContext{
		validatorsPower: make(map[string]*big.Int),
		totalPower:      big.NewInt(0),
		rounds:          make(map[int32]*roundInfo),
		aggregators:     make(map[int32]*worker),
	}
}
