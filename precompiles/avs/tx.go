package avs

import (
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	avskeeper "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

const (
	MethodRegisterAVS               = "registerAVS"
	MethodUpdateAVS                 = "updateAVS"
	MethodDeregisterAVS             = "deregisterAVS"
	MethodRegisterOperatorToAVS     = "registerOperatorToAVS"
	MethodDeregisterOperatorFromAVS = "deregisterOperatorFromAVS"
	MethodCreateAVSTask             = "createTask"
	MethodSubmitProof               = "submitProof"
	MethodRegisterBLSPublicKey      = "registerBLSPublicKey"
	MethodGetRegisteredPubkey       = "getRegisteredPubkey"
	MethodGetOptinOperators         = "getOptInOperators"
)

// AVSInfoRegister register the avs related information and change the state in avs keeper module.
func (p Precompile) RegisterAVS(
	ctx sdk.Context,
	origin common.Address,
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
	if !slices.Contains(avsParams.AvsOwnerAddress, sdk.AccAddress(origin.Bytes()).String()) {
		return nil, errorsmod.Wrap(err, "not qualified to registerOrDeregister")
	}
	// The AVS registration is done by the calling contract.
	avsParams.AvsAddress = contract.CallerAddress.String()
	avsParams.Action = avskeeper.RegisterAction
	// Finally, update the AVS information in the keeper.
	err = p.avsKeeper.AVSInfoUpdate(ctx, avsParams)
	if err != nil {
		fmt.Println("Failed to update AVS info", err)
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (p Precompile) DeregisterAVS(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodDeregisterAVS].Inputs) {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodDeregisterAVS].Inputs), len(args))
	}
	avsParams := &avstypes.AVSRegisterOrDeregisterParams{}

	avsName, ok := args[0].(string)
	if !ok || avsName == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", avsName)
	}
	avsParams.AvsName = avsName

	avsParams.AvsAddress = contract.CallerAddress.String()
	avsParams.Action = avskeeper.DeRegisterAction
	// validates that this is owner
	avsParams.CallerAddress = sdk.AccAddress(origin[:]).String()

	err := p.avsKeeper.AVSInfoUpdate(ctx, avsParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) UpdateAVS(
	ctx sdk.Context,
	origin common.Address,
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

	avsParams.AvsAddress = contract.CallerAddress.String()
	avsParams.CallerAddress = sdk.AccAddress(origin[:]).String()
	avsParams.Action = avskeeper.UpdateAction
	err = p.avsKeeper.AVSInfoUpdate(ctx, avsParams)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (p Precompile) BindOperatorToAVS(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	_ []interface{},
) ([]byte, error) {
	operatorParams := &avskeeper.OperatorOptParams{}
	operatorParams.OperatorAddress = sdk.AccAddress(origin[:]).String()
	operatorParams.AvsAddress = contract.CallerAddress.String()
	operatorParams.Action = avskeeper.RegisterAction
	err := p.avsKeeper.OperatorOptAction(ctx, operatorParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) UnbindOperatorToAVS(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	_ []interface{},
) ([]byte, error) {
	operatorParams := &avskeeper.OperatorOptParams{}
	operatorParams.OperatorAddress = sdk.AccAddress(origin[:]).String()
	operatorParams.AvsAddress = contract.CallerAddress.String()
	operatorParams.Action = avskeeper.DeRegisterAction
	err := p.avsKeeper.OperatorOptAction(ctx, operatorParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// CreateAVSTask Middleware uses exocore's default avstask template to create tasks in avstask module.
func (p Precompile) CreateAVSTask(
	ctx sdk.Context,
	origin common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	params, err := p.GetTaskParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	params.TaskContractAddress = contract.CallerAddress.String()
	params.CallerAddress = sdk.AccAddress(origin[:]).String()
	err = p.avsKeeper.CreateAVSTask(ctx, params)
	if err != nil {
		return nil, err
	}
	if err = p.EmitCreateAVSTaskEvent(ctx, stateDB, params); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// RegisterBLSPublicKey
func (p Precompile) RegisterBLSPublicKey(
	ctx sdk.Context,
	origin common.Address,
	_ *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodRegisterBLSPublicKey].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodRegisterBLSPublicKey].Inputs), len(args))
	}
	blsParams := &avskeeper.BlsParams{}
	blsParams.Operator = sdk.AccAddress(origin[:]).String()
	name, ok := args[1].(string)
	if !ok || name == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", name)
	}
	blsParams.Name = name

	pubkeyBz, ok := args[2].([]byte)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", pubkeyBz)
	}
	blsParams.PubKey = pubkeyBz

	pubkeyRegistrationSignature, ok := args[3].([]byte)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "[]byte", pubkeyRegistrationSignature)
	}
	blsParams.PubkeyRegistrationSignature = pubkeyRegistrationSignature

	pubkeyRegistrationMessageHash, ok := args[4].([]byte)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "[]byte", pubkeyRegistrationMessageHash)
	}
	blsParams.PubkeyRegistrationMessageHash = pubkeyRegistrationMessageHash

	err := p.avsKeeper.RegisterBLSPublicKey(ctx, blsParams)
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
	if len(args) != len(p.ABI.Methods[MethodGetRegisteredPubkey].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 1, len(args))
	}
	// the key is set using the operator's acc address so the same logic should apply here
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

// SubmitProof
func (p Precompile) SubmitProof(
	_ sdk.Context,
	_ common.Address,
	_ *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != len(p.ABI.Methods[MethodSubmitProof].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodSubmitProof].Inputs), len(args))
	}
	// TODO: check whether this address is acc or hex and decide the type accordingly
	addr, ok := args[0].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", addr)
	}
	// TODO implement SubmitProof
	// err := p.avsKeeper.SubmitProof(ctx, addr)
	// if err != nil {
	//	return nil, err
	//}
	return method.Outputs.Pack(true)
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
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", addr)
	}

	list, err := p.avsKeeper.GetOptInOperators(ctx, addr.String())
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(list)
}
