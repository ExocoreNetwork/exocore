package types

func ContainsString(slice []string, target string) bool {
	for _, str := range slice {
		if str == target {
			return true
		}
	}
	return false
}
func RemoveOperatorAddress(addresses []string, addressToRemove string) []string {
	for i, address := range addresses {
		if address == addressToRemove {
			addresses[i] = addresses[len(addresses)-1]
			return addresses[:len(addresses)-1]
		}
	}
	return addresses
}
