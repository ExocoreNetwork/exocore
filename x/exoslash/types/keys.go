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
<<<<<<< HEAD
	prefixParams       = iota + 1
	prefixOperatorInfo = iota + 1
=======
	prefixParams = iota + 1
>>>>>>> eebca7f (implement slash interface)
)

var (
	KeyPrefixParams = []byte{prefixParams}
<<<<<<< HEAD
	// KeyPrefixOperatorInfo key-value: operatorAddr->operatorInfo
	KeyPrefixOperatorInfo = []byte{prefixOperatorInfo}
	ParamsKey             = []byte("Params")
=======

	ParamsKey = []byte("Params")
>>>>>>> eebca7f (implement slash interface)
)
