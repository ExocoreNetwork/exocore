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
<<<<<<< HEAD
	prefixParams       = iota + 1
	prefixOperatorInfo = iota + 1
=======
	prefixParams = iota + 1
>>>>>>> eebca7f (implement slash interface)
=======
	prefixParams       = iota + 1
	prefixOperatorInfo = iota + 1
>>>>>>> 5429dca (add unti test for slash and fix some  bugs)
)

var (
	KeyPrefixParams = []byte{prefixParams}
<<<<<<< HEAD
<<<<<<< HEAD
	// KeyPrefixOperatorInfo key-value: operatorAddr->operatorInfo
	KeyPrefixOperatorInfo = []byte{prefixOperatorInfo}
	ParamsKey             = []byte("Params")
=======

	ParamsKey = []byte("Params")
>>>>>>> eebca7f (implement slash interface)
=======
	// KeyPrefixOperatorInfo key-value: operatorAddr->operatorInfo
	KeyPrefixOperatorInfo = []byte{prefixOperatorInfo}
	ParamsKey             = []byte("Params")
>>>>>>> 5429dca (add unti test for slash and fix some  bugs)
)
