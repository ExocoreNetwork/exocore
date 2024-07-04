package types

const (
	// ModuleName defines the module name
	ModuleName = "feedistribution"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_feedistribute"
)

var (
	ParamsKey = []byte("p_feedistribute")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
