package keeper_test

import (
	"fmt"
	utiltx "github.com/ExocoreNetwork/exocore/testutil/tx"
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func Test_difference(t *testing.T) {
	arr1 := []string{"apple", "banana", "cherry"}
	arr2 := []string{"apple", "cherry", "date"}

	diff := types.Difference(arr1, arr2)
	fmt.Println("Differences:", diff)
	num1 := sdk.MustNewDecFromStr("1.3")
	num2 := sdk.MustNewDecFromStr("12.3")

	// Perform division
	result := num1.Quo(num2)

	// Convert result to percentage
	percentage := result.Mul(sdk.NewDec(100))

	// Print the result
	fmt.Print(percentage)
}

func Test_genKey(t *testing.T) {
	addresses := make([]string, 5)

	for i := 0; i < 5; i++ {
		address := utiltx.GenerateAddress()
		exoAddress := sdk.AccAddress(address.Bytes()).String()
		addresses[i] = exoAddress
	}

	fmt.Println("Generated EXO addresses:")
	for _, address := range addresses {
		fmt.Println(address)
	}
}
