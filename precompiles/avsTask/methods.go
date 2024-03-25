package task

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MethodCreateNewTask defines the ABI method name for the task
	//  transaction.
	MethodCreateNewTask   = "createNewTask"
	MethodIsOperatorOptin = "isOperatorOptin"
)

// CreateNewTask Middleware uses exocore's default task template to create tasks in task module.
func (p Precompile) CreateNewTask(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	stateDB vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract
	flag := p.avsKeeper.IsAVS(ctx, sdk.AccAddress(contract.CallerAddress.String()))
	if !flag {
		return nil, fmt.Errorf(ErrNotYetRegistered, contract.CallerAddress)
	}

	createNewTaskParams, err := p.GetTaskParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	createNewTaskParams.ContractAddr = contract.CallerAddress.String()
	createNewTaskParams.TaskCreatedBlock = ctx.BlockHeight()
	_, err = p.taskKeeper.CreateNewTask(ctx, createNewTaskParams)
	if err != nil {
		return nil, err
	}
	if err = p.EmitNewTaskCreatedEvent(
		ctx,
		stateDB,
		createNewTaskParams.TaskIndex,
		createNewTaskParams.NumberToBeSquared,
		createNewTaskParams.QuorumNumbers,
		createNewTaskParams.QuorumThresholdPercentage,
	); err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// IsOperatorOptin Middleware uses exocore's default task template to create tasks in task module.
func (p Precompile) IsOperatorOptin(
	ctx sdk.Context,
	contract *vm.Contract,
	_ *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract
	flag := p.avsKeeper.IsAVS(ctx, sdk.AccAddress(contract.CallerAddress.String()))
	if !flag {
		return nil, fmt.Errorf(ErrNotYetRegistered, contract.CallerAddress)
	}

	return method.Outputs.Pack(true)
}
