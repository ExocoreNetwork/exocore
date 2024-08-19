package keeper_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"testing"
)

var (
	structThing, _ = abi.NewType("tuple", "struct", []abi.ArgumentMarshaling{
		{Name: "field_one", Type: "uint256"},
		{Name: "field_two", Type: "address"},
	})

	args = abi.Arguments{
		{Type: structThing, Name: "param_one"},
	}
)

func TestReceiptMarshalBinary(t *testing.T) {

	record := struct {
		FieldOne *big.Int
		FieldTwo common.Address
	}{
		big.NewInt(2e18),
		common.HexToAddress("0x0002"),
	}

	packed, err := args.Pack(&record)
	if err != nil {
		fmt.Println("bad bad ", err)
		return
	} else {
		fmt.Println("abi encoded", hexutil.Encode(packed))
	}
	b, _ := args.Unpack(packed)
	fmt.Println("unpacked", b)
}
