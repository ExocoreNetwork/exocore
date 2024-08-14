package types

const (
	// NonceKeyPrefix is the prefix to retrieve all KeyNonce
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
