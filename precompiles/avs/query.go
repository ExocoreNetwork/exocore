package avs

import (
	"fmt"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

const (
	MethodGetRegisteredPubkey      = "getRegisteredPubkey"
	MethodGetOptinOperators        = "getOptInOperators"
	MethodGetAVSUSDValue           = "getAVSUSDValue"
	MethodGetOperatorOptedUSDValue = "getOperatorOptedUSDValue"
)

// GetRegisteredPubkey
func (p Precompile) GetRegisteredPubkey(
	ctx sdk.Context,
	_ *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodGetRegisteredPubkey].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}
	// the key is set using the operator's acc address so the same logic should apply here
	addr, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", addr)
	}
	pubkey, err := p.avsKeeper.GetOperatorPubKey(ctx, addr)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(pubkey)
}

// GetOptedInOperatorAccAddrs
func (p Precompile) GetOptedInOperatorAccAddrs(
	ctx sdk.Context,
	_ *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodGetOptinOperators].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	addr, ok := args[0].(common.Address)
	if !ok || addr == (common.Address{}) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", addr)
	}

	list, err := p.avsKeeper.GetOptInOperators(ctx, addr.String())
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(list)
}

// GetAVSUSDValue is a function to retrieve the USD share of specified Avs,
func (p Precompile) GetAVSUSDValue(
	ctx sdk.Context,
	_ *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodGetAVSUSDValue].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}
	addr, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", addr)
	}
	amount, err := p.operatorKeeper.GetAVSUSDValue(ctx, addr)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(amount.String())
}

// GetOperatorOptedUSDValue is a function to retrieve the USD share of specified operator and Avs,
func (p Precompile) GetOperatorOptedUSDValue(
	ctx sdk.Context,
	_ *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodGetOperatorOptedUSDValue].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}
	avsAddr, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", avsAddr)
	}
	operatorAddr, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", operatorAddr)
	}
	amount, err := p.operatorKeeper.GetOperatorOptedUSDValue(ctx, avsAddr, operatorAddr)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(amount.ActiveUSDValue.String())
}
