package types

const (
	// ModuleName defines the module name
	ModuleName = "feedistribution"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_feedistribute"
)

const (
	// EpochIdentifier defines the epoch identifier for fee distribution module
	EpochIdentifierKey = "epoch_identifier_feedistribute"
)

var (
	ParamsKey                = []byte("p_feedistribute")
	KeyPrefixEpochIdentifier = []byte(EpochIdentifierKey)
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
