package types

const (
	// ModuleName defines the module name
	ModuleName = "exoslash"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_exoslash"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	prefixParams = iota + 1
)

var (
	KeyPrefixParams = []byte{prefixParams}

	ParamsKey = []byte("Params")
)
