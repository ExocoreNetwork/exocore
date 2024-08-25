package types

const (
	// ModuleName defines the module name
	ModuleName = "subscriber"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_subscriber"

	// PortID is the default port id that module binds to
	PortID = "subscriber"

	SubscriberRedistributeName = "subscriber_redistribute"

	SubscriberToSendToCoordinatorName = "subscriber_to_send_to_coordinator"
)

const (
	ParamsBytePrefix byte = iota + 1
)

func ParamsKey() []byte {
	return []byte{ParamsBytePrefix}
}
