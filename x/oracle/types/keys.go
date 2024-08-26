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

	// TODO: rename for prefix and keys
	ValidatorsKey = "Validators/value/"

	ValidatorUpdateBlockKey = "ValidatorUpdateBlock/value/"

	IndexRecentParamsKey = "IndexRecentParams/value/"

	IndexRecentMsgKey = "IndexRecentMsg/value/"
)

var ParamsKey = []byte{0x11}

func KeyPrefix(p string) []byte {
	return []byte(p)
}
