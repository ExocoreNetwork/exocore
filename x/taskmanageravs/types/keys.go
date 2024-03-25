package types

const (
	// ModuleName defines the module name
	ModuleName = "taskmanageravs"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_taskmanageravs"
)

const (
	prefixAVSTaskInfo = iota + 1
	prefixAVSTaskMap
	PrefixAvsTaskIdMap
)

var (
	// KeyPrefixAVSTaskInfo key-value: avsAddr->AVSTaskInfo
	KeyPrefixAVSTaskInfo  = []byte{prefixAVSTaskInfo}
	KeyPrefixAVSTaskMap   = []byte{prefixAVSTaskMap}
	KeyPrefixAvsTaskIdMap = []byte{PrefixAvsTaskIdMap}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
