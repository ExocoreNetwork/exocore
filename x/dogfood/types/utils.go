package types

import (
	"bytes"
	"reflect"
)

// RemoveFromBytesList removes an address from a list of addresses
// or a byte slice from a list of byte slices.
func RemoveFromBytesList(list [][]byte, addr []byte) [][]byte {
	for i, a := range list {
		if bytes.Equal(a, addr) {
			return append(list[:i], list[i+1:]...)
		}
	}
	panic("address not found in list")
}

func PanicIfZeroOrNil(x interface{}, msg string) {
	if x == nil || reflect.ValueOf(x).IsZero() {
		panic("zero or nil value for " + msg)
	}
}
