# State for Aggregator rechache
## []Params: used as circle
- TokenFeeder: (active/inactive)
- block -> Params
- length = x
### Spec.
only latest Params matter for recache
## ValidatorSet
- block -> []validator{address, power}
## []msgCreatePrice: used as circle
- block -> []msgCreatePrice
---
x: configured in Params
# Recache
Current block : h
## params update:
1,2,[3],4,5
a. start feeder: no affect on former round/txs
b. stop feeder: related round/txs sealed at H=3
Just recache from earliest block-2
- msg in or before block-1 must have been seald in or before block-5, means the windows are closed, don't need cache for the collection
## validatoset update:
1,2,[3],4,5
Validator set changed on block-3, this would be done in EndBlock of block-3, and all live round/txs would be seald here, means we don't need to recache any msgs from 1,2,3
Just remove all block/txs info in 1,2,3 in block-3 Endblock after sealed all live round/txs
## workflow =>
h0 = h-x
reProcess blocks from max{h0, min_[]msgCreatePrice.block} to fill the cache
---
# Aggregator
1. initCache if cache is nil
- if []msgCreatePrice is empty just load Params and validatorSet
- if []msgCreatePrice is not empty, recache the `Aggregator`

## Triggerd
1. create-price-service(): msg accepted call collectPrice()-> initCache if nil
2. EndBlock()
Seal()-> initCache if nil
- check params update
-- stop feeder(update live feeder's EndBlock): seal related round/txs
--- if any feeder stop at current block, no realted txs will be accepted,
-- new feeder(for new token from currently service)

- check validatorset update:
-- true: seal all alive round/txs here, update aggregator's status: 1->2
-- clear all []historyBlock-mem (which would be persit in KV)

- aggregate()
-- initCache if nil
-- check all live round and:
--- consensus reached, then the round seal with sucess new price(basedBlock), status:1->2
--- consensus not reached yet, but
---- feeder stops here(set by params), seal with previous price(basedBlock), status: 1->2//and clear corresponding roundInfo[feederId], dealed in postAggregation
---- window ends here, seal with previous price(basedBlock), udpate the related roundInfo[feederId].status:1->2
---- validatorset changed: seal with previous price(basedBlock), update the related roundInfo[feeerId].status:1->2
** this change will be active on next block, so we should seal here in front

- postAggregation()
-- remove all stopped feeders related roundInfo[feederId]
-- if validator changed, remove all roundData

- Prepare() -> initCache if nil
-- params is up to date already
-- based on current blockHeight and params(feeder)_ fromAggregatorCache, is there any possible status: 0->1(with info filled), 2->1(with info update)
-- update validatorset if changed

- Persist memCache for recache
- if validatorset changed:
-- clear all the inMem_[]msgCreatePrice, nothing to persist, and clear KV's []msgCreatePrice
-- set inMem_[]validatorSet
- if Parmas cahged:
-- set Params
- msgCreatePrice, if not nil, append to KV
- params, if not nil, append to KV
- validatorset, if not nil, Set to KV(only keep one validatorSet)

# params update
## updateTime, activeTime
TODO: When set up a new feeder, the feeder activeTime should be later than updateTime(may be dynamic)
