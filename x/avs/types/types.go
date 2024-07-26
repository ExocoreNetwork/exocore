package types

const (
	ChainID = "chain-id-code"
)

func RemoveOperatorAddress(addresses []string, addressToRemove string) []string {
	for i, address := range addresses {
		if address == addressToRemove {
			addresses[i] = addresses[len(addresses)-1]
			return addresses[:len(addresses)-1]
		}
	}
	return addresses
}
