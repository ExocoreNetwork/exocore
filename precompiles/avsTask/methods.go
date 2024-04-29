package task

import (
	"fmt"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	"github.com/cosmos/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

const (
	// MethodRegisterAVSTask defines the ABI method name for the avstask
	//  transaction.
	MethodRegisterAVSTask      = "registerAVSTask"
	MethodRegisterBLSPublicKey = "registerBLSPublicKey"
	MethodGetRegisteredPubkey  = "getRegisteredPubkey"
)

// RegisterAVSTask Middleware uses exocore's default avstask template to create tasks in avstask module.
func (p Precompile) RegisterAVSTask(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract
	callerAddress, _ := bech32.EncodeFromBase256("exo", contract.CallerAddress.Bytes())
	_, err := p.avsKeeper.GetAVSInfo(ctx, callerAddress)

	params, err := p.GetTaskParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	params.FromAddress = callerAddress
	_, err = p.taskKeeper.RegisterAVSTask(ctx, params)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true)
}

// RegisterBLSPublicKey
func (p Precompile) RegisterBLSPublicKey(
	ctx sdk.Context,
	_ common.Address,
	stateDB vm.StateDB,
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

	err := p.taskKeeper.SetOperatorPubKey(ctx, addr, pubkeyBz)
	if err != nil {
		return nil, err
	}
	err = p.EmitEventTypeNewPubkeyRegistration(
		ctx,
		stateDB,
		addr,
		pubkeyBz,
	)
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

	pubkey, err := p.taskKeeper.GetOperatorPubKey(ctx, addr)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack([]byte(pubkey))
}
