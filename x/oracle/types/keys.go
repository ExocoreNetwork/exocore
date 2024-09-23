package types

const (
	// ModuleName defines the module name
	ModuleName = "oracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_oracle"
)

var (
	// ParamsKey defines the key to store the params in store
	ParamsKey = []byte{0x11}
	// BlockKey stores the last validator update block
	BlockKey = []byte{0x0}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	ValidatorsKey = "Validators/value/"
)

const (
	ValidatorUpdateBlockKey = "ValidatorUpdateBlock/value/"
)

const (
	IndexRecentParamsKey = "IndexRecentParams/value/"
)

const (
	IndexRecentMsgKey = "IndexRecentMsg/value/"
)
