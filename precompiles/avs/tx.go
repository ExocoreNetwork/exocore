package avs

import (
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	util "github.com/ExocoreNetwork/exocore/utils"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

const (
	MethodAVSAction      = "AVSAction"
	MethodOperatorAction = "OperatorOptAction"
)

// AVSInfoRegister register the avs related information and change the state in avs keeper module.
func (p Precompile) AVSInfoRegisterOrDeregister(
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
	avsAddress, err := util.ProcessAvsAddress(contract.Address().String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse avsAddress error")
	}

	callerAddress, err := util.ProcessAvsAddress(contract.CallerAddress.String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse callerAddress error")
	}

	if !slices.Contains(avsParams.AvsOwnerAddress, callerAddress) {
		return nil, errorsmod.Wrap(err, "not qualified to registerOrDeregister")
	}

	avsParams.AvsAddress = avsAddress
	if err != nil {
		return nil, err
	}
	err = p.avsKeeper.AVSInfoUpdate(ctx, avsParams)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

func (p Precompile) OperatorBindingAvs(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 2, len(args))
	}
	operatorParams := &avstypes.OperatorOptParams{}
	action, ok := args[0].(uint64)
	if !ok || (action != avstypes.RegisterAction && action != avstypes.DeRegisterAction) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint64", action)
	}
	operatorParams.Action = action

	callerAddress, err := util.ProcessAvsAddress(contract.CallerAddress.String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse callerAddress error")
	}

	operatorParams.OperatorAddress = callerAddress

	avsAddress, err := util.ProcessAvsAddress(contract.Address().String())
	if err != nil {
		return nil, errorsmod.Wrap(err, "parse avsAddress error")
	}
	operatorParams.AvsAddress = avsAddress
	if err != nil {
		return nil, err
	}
	err = p.avsKeeper.AVSInfoUpdateWithOperator(ctx, operatorParams)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}
