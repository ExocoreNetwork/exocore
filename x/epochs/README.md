# Epochs

## References

This module is backported from Cosmos SDK, which itself [upstreamed](https://github.com/cosmos/cosmos-sdk/pull/19697) it
from Osmosis. Much of the documentation below is sourced from the [README in that pull request.]( https://github.com/cosmos/cosmos-sdk/blob/8b171c5a4300fb3287ab591b5e6cf221019ddb07/x/epochs/README.md)

## Abstract

Often in the SDK, we would like to run certain code every-so often. The
purpose of `epochs` module is to allow other modules to set that they
would like to be signaled once every period. So another module can
specify it wants to execute code once a week, starting at UTC-time = x.
`epochs` creates a generalized epoch interface to other modules so that
they can easily be signaled upon such events.

## Contents

1. **[Concept](#concepts)**
2. **[State](#state)**
3. **[Events](#events)**
4. **[Keeper](#keepers)**
5. **[Hooks](#hooks)**
6. **[Queries](#queries)**

## Concepts

The epochs module defines on-chain timers that execute at fixed time intervals.
Other SDK modules can then register logic to be executed at the timer ticks.
We refer to the period in between two timer ticks as an "epoch".

Every timer has a unique identifier.
Every epoch will have a start time, and an end time, where `end time = start time + timer interval`.

The timer will tick at the first block whose block time is greater than the timer end time,
and set the start as the prior timer end time. (Notably, it's _not_ set to the block time!)
This means that if the chain has been down for a while, you will get **one timer tick per block**,
until the timer has caught up.

## State

The Epochs module keeps a single `EpochInfo` per identifier.
This contains the current state of the timer with the corresponding identifier.
Its fields are modified at every timer tick.
EpochInfos are initialized as part of genesis initialization or upgrade logic,
and are only modified during the begin blocker.

## Hooks

```go
  // the first block whose timestamp is after the duration is counted as the end of the epoch
  AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64)
  // new epoch is next block of epoch end block
  BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64)
```

### How modules receive hooks

On hook receiver function of other modules, they need to filter
`epochIdentifier` and only do executions for only specific
epochIdentifier. Filtering epochIdentifier could be in `Params` of other
modules so that they can be modified by governance.

This is the standard dev UX of this:

```golang
func (k MyModuleKeeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
    params := k.GetParams(ctx)
    if epochIdentifier == params.DistrEpochIdentifier {
    // my logic
  }
}
```

## Queries

The Epochs module provides the following queries to check the module's state.

```protobuf
service Query {
  // EpochInfos provide running epochInfos
  rpc EpochInfos(QueryEpochsInfoRequest) returns (QueryEpochsInfoResponse) {}
  // CurrentEpoch provide current epoch of specified identifier
  rpc CurrentEpoch(QueryCurrentEpochRequest) returns (QueryCurrentEpochResponse) {}
}
```
