package types

const (
	// ModuleName is the name of the module
	ModuleName = "appchain"
	// Version is the current version of the module
	Version = ModuleName + "-1"
	// CoordinatorPortID is the default port id to which the coordinator module binds
	CoordinatorPortID = "coordinator"
	// SubscriberPortID is the default port id to which the subscriber module binds
	SubscriberPortID = "subscriber"
	// StoreKey defines the store key for the module (used in tests)
	StoreKey = ModuleName
	// MemStoreKey defines the in-memory store key (used in tests)
	MemStoreKey = "mem_appchain"
)
