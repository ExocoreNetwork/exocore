package bls

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

var _ vm.PrecompiledContract = &Precompile{}

// Embed abi json file to the executable binary. Needed when importing as dependency.
//
//go:embed abi.json
var f embed.FS

// Precompile defines the precompiled contract for deposit.
type Precompile struct {
	abi.ABI
	baseGas uint64
}

// NewPrecompile creates a new BLS Precompile instance as a
// PrecompiledContract interface.
func NewPrecompile(baseGas uint64) (*Precompile, error) {
	abiBz, err := f.ReadFile("abi.json")
	if err != nil {
		return nil, fmt.Errorf("error loading the deposit ABI %s", err)
	}

	newABI, err := abi.JSON(bytes.NewReader(abiBz))
	if err != nil {
		return nil, fmt.Errorf(cmn.ErrInvalidABI, err)
	}

	if baseGas == 0 {
		return nil, fmt.Errorf("baseGas cannot be zero")
	}

	return &Precompile{
		ABI:     newABI,
		baseGas: baseGas,
	}, nil
}

// Address defines the address of the deposit compile contract.
// address: 0x0000000000000000000000000000000000000809
func (p Precompile) Address() common.Address {
	return common.HexToAddress("0x0000000000000000000000000000000000000809")
}

// RequiredGas calculates the precompiled contract's base gas rate.
func (p Precompile) RequiredGas(_ []byte) uint64 {
	return p.baseGas
}

// Run executes the precompiled contract deposit methods defined in the ABI.
func (p Precompile) Run(_ *vm.EVM, contract *vm.Contract, _ bool) (bz []byte, err error) {
	methodID := contract.Input[:4]
	// NOTE: this function iterates over the method map and returns
	// the method with the given ID
	method, err := p.MethodById(methodID)
	if err != nil {
		return nil, err
	}

	argsBz := contract.Input[4:]
	args, err := method.Inputs.Unpack(argsBz)
	if err != nil {
		return nil, err
	}

	switch method.Name {
	case MethodVerify:
		bz, err = p.Verify(method, args)
	case MethodFastAggregateVerify:
		bz, err = p.FastAggregateVerify(method, args)
	case MethodGeneratePrivateKey:
		bz, err = p.GeneratePrivateKey(method, args)
	case MethodPublicKey:
		bz, err = p.PublicKey(method, args)
	case MethodSign:
		bz, err = p.Sign(method, args)
	case MethodAggregatePubkeys:
		bz, err = p.AggregatePubkeys(method, args)
	case MethodAggregateSignatures:
		bz, err = p.AggregateSignatures(method, args)
	case MethodAddTwoPubkeys:
		bz, err = p.AddTwoPubkeys(method, args)
	default:
		return nil, fmt.Errorf("invalid method")
	}

	if err != nil {
		return nil, err
	}

	// fmt.Printf("calling method: %s\n", method.Name)
	// fmt.Printf("precompile result: %x\n", bz)
	return bz, nil
}
