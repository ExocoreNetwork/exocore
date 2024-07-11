package types

const (
	// ModuleName defines the module name
	ModuleName = "epochs"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

const (
	// prefixEpoch is the byte prefix used by the epoch store.
	prefixEpoch = iota + 1
)

// KeyPrefixEpoch is the byte-array prefix used by the epoch store.
var KeyPrefixEpoch = []byte{prefixEpoch}

// KeyEpoch returns the key for the epoch with the given identifier.
func KeyEpoch(identifier string) []byte {
	return append(KeyPrefixEpoch, []byte(identifier)...)
}
