package types

const (
	// ModuleName defines the module name
	ModuleName = "coordinator"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_coordinator"

	// PortId is the default port id that the module binds to
	PortId = "coordinator"

	// SubscriberRewardsPool is the address that receives the rewards from the subscriber
	// chains. Technically, it is possible for the subscriber chain to send these rewards
	// directly to the FeeCollector, but this intermediate step allows the coordinator
	// module to ensure that the rewards are actually sent to us.
	SubscriberRewardsPool = "subscriber_rewards_pool"
)

const (
	ParamsBytePrefix byte = iota + 1
)

func ParamsKey() []byte {
	return []byte{ParamsBytePrefix}
}
