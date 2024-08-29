package types

const (
	// NonceKeyPrefix is used as a prefix for storing key nonces.
	NonceKeyPrefix = "KeyNonce/value/"
)

func NonceKey(
	validator string,
) []byte {
	var key []byte

	key = append(key, validator...)
	key = append(key, []byte("/")...)

	return key
}
