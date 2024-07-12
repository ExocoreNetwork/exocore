package avs

import (
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	util "github.com/ExocoreNetwork/exocore/utils"
	avskeeper "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

const (
	MethodRegisterOperatorToAVS     = "RegisterOperatorToAVS"
	MethodDeregisterOperatorFromAVS = "DeregisterOperatorFromAVS"
	MethodRegisterAVS               = "RegisterAVS"
	MethodUpdateAVS                 = "UpdateAVS"
	MethodDeregisterAVS             = "DeregisterAVS"
	MethodRegisterAVSTask           = "registerAVSTask"
	MethodRegisterBLSPublicKey      = "registerBLSPublicKey"
	MethodGetRegisteredPubkey       = "getRegisteredPubkey"
)

// AVSInfoRegister register the avs related information and change the state in avs keeper module.
func (p Precompile) RegisterAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// parse the avs input params first.
	avsParams, err := p.GetAVSParamsFromInputs(ctx, args)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse args error")
	}
	avsAddress, err := util.ProcessAddress(contract.Address().String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse avsAddress error")
	}

	callerAddress, err := util.ProcessAddress(contract.CallerAddress.String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse callerAddress error")
	}

	if !slices.Contains(avsParams.AvsOwnerAddress, callerAddress) {
		return nil, errorsmod.Wrap(err, "not qualified to registerOrDeregister")
	}

	avsParams.AvsAddress = avsAddress
	avsParams.Action = avskeeper.RegisterAction

	if err != nil {
		return nil, err
	}
	err = p.avsKeeper.AVSInfoUpdate(ctx, avsParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) DeregisterAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodDeregisterAVS].Inputs) {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodDeregisterAVS].Inputs), len(args))
	}
	avsParams := &avskeeper.AVSRegisterOrDeregisterParams{}

	avsName, ok := args[0].(string)
	if !ok || avsName == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", avsName)
	}
	avsParams.AvsName = avsName

	avsAddress, err := util.ProcessAddress(contract.Address().String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse avsAddress error")
	}
	avsParams.AvsAddress = avsAddress

	callerAddress, err := util.ProcessAddress(contract.CallerAddress.String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse callerAddress error")
	}
	avsParams.CallerAddress = callerAddress

	avsParams.Action = avskeeper.DeRegisterAction

	if err != nil {
		return nil, err
	}
	err = p.avsKeeper.AVSInfoUpdate(ctx, avsParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) UpdateAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// parse the avs input params first.
	avsParams, err := p.GetAVSParamsFromUpdateInputs(ctx, args)
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse args error")
	}
	avsAddress, err := util.ProcessAddress(contract.Address().String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse avsAddress error")
	}
	avsParams.AvsAddress = avsAddress

	callerAddress, err := util.ProcessAddress(contract.CallerAddress.String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse callerAddress error")
	}
	avsParams.CallerAddress = callerAddress

	avsParams.Action = avskeeper.UpdateAction

	if err != nil {
		return nil, err
	}
	err = p.avsKeeper.AVSInfoUpdate(ctx, avsParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) BindOperatorToAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	_ []interface{},
) ([]byte, error) {
	operatorParams := &avskeeper.OperatorOptParams{}

	callerAddress, err := util.ProcessAddress(contract.CallerAddress.String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse callerAddress error")
	}

	operatorParams.OperatorAddress = callerAddress

	avsAddress, err := util.ProcessAddress(contract.Address().String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse avsAddress error")
	}
	operatorParams.AvsAddress = avsAddress
	operatorParams.Action = avskeeper.RegisterAction
	err = p.avsKeeper.AVSInfoUpdateWithOperator(ctx, operatorParams)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (p Precompile) UnbindOperatorToAVS(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	_ []interface{},
) ([]byte, error) {
	operatorParams := &avskeeper.OperatorOptParams{}
	callerAddress, err := util.ProcessAddress(contract.CallerAddress.String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse callerAddress error")
	}

	operatorParams.OperatorAddress = callerAddress

	avsAddress, err := util.ProcessAddress(contract.Address().String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse avsAddress error")
	}
	operatorParams.AvsAddress = avsAddress
	operatorParams.Action = avskeeper.DeRegisterAction

	if err != nil {
		return nil, err
	}
	err = p.avsKeeper.AVSInfoUpdateWithOperator(ctx, operatorParams)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

// RegisterAVSTask Middleware uses exocore's default avstask template to create tasks in avstask module.
func (p Precompile) RegisterAVSTask(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract
	callerAddress, _ := util.ProcessAddress(contract.CallerAddress.String())
	params, err := p.GetTaskParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	params.FromAddress = callerAddress
	_, err = p.avsKeeper.RegisterAVSTask(ctx, params)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// RegisterBLSPublicKey
func (p Precompile) RegisterBLSPublicKey(
	ctx sdk.Context,
	_ common.Address,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}

	addr, ok := args[0].(string)
	if !ok || addr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", addr)
	}

	pubkeyBz, ok := args[1].([]byte)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "[]byte", pubkeyBz)
	}

	err := p.avsKeeper.SetOperatorPubKey(ctx, addr, pubkeyBz)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// GetRegisteredPubkey
func (p Precompile) GetRegisteredPubkey(
	ctx sdk.Context,
	_ *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}

	addr, ok := args[0].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", addr)
	}

	pubkey, err := p.avsKeeper.GetOperatorPubKey(ctx, addr)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(pubkey)
}
