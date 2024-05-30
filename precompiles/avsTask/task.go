package task

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/avs/keeper"
	taskKeeper "github.com/ExocoreNetwork/exocore/x/avstask/keeper"
	"github.com/cometbft/cometbft/libs/log"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
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

// Precompile defines the precompiled contract for avstask.
type Precompile struct {
	cmn.Precompile
	taskKeeper taskKeeper.Keeper
	avsKeeper  keeper.Keeper
}

// NewPrecompile creates a new avstask Precompile instance as a
// PrecompiledContract interface.
func NewPrecompile(
	authzKeeper authzkeeper.Keeper,
	taskKeeper taskKeeper.Keeper,
	avsKeeper keeper.Keeper,
) (*Precompile, error) {
	abiBz, err := f.ReadFile("abi.json")
	if err != nil {
		return nil, fmt.Errorf("error loading the avstask ABI %s", err)
	}

	newAbi, err := abi.JSON(bytes.NewReader(abiBz))
	if err != nil {
		return nil, fmt.Errorf(cmn.ErrInvalidABI, err)
	}

	return &Precompile{
		Precompile: cmn.Precompile{
			ABI:                  newAbi,
			AuthzKeeper:          authzKeeper,
			KvGasConfig:          storetypes.KVGasConfig(),
			TransientKVGasConfig: storetypes.TransientGasConfig(),
			ApprovalExpiration:   cmn.DefaultExpirationDuration, // should be configurable in the future.
		},
		taskKeeper: taskKeeper,
		avsKeeper:  avsKeeper,
	}, nil
}

// Address defines the address of the avstask compile contract.
// address: 0x0000000000000000000000000000000000000901
func (p Precompile) Address() common.Address {
	return common.HexToAddress("0x0000000000000000000000000000000000000901")
}

// RequiredGas calculates the precompiled contract's base gas rate.
func (p Precompile) RequiredGas(input []byte) uint64 {
	methodID := input[:4]

	method, err := p.MethodById(methodID)
	if err != nil {
		// This should never happen since this method is going to fail during Run
		return 0
	}

	return p.Precompile.RequiredGas(input, p.IsTransaction(method.Name))
}

// Run executes the precompiled contract RegisterOrDeregisterAVSInfo methods defined in the ABI.
func (p Precompile) Run(evm *vm.EVM, contract *vm.Contract, readOnly bool) (bz []byte, err error) {
	ctx, stateDB, method, initialGas, args, err := p.RunSetup(evm, contract, readOnly, p.IsTransaction)
	if err != nil {
		return nil, err
	}

	// This handles any out of gas errors that may occur during the execution of a precompile tx or query.
	// It avoids panics and returns the out of gas error so the EVM can continue gracefully.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	switch method.Name {
	case MethodRegisterAVSTask:
		bz, err = p.RegisterAVSTask(ctx, evm.Origin, contract, method, args)
	case MethodRegisterBLSPublicKey:
		bz, err = p.RegisterBLSPublicKey(ctx, evm.Origin, stateDB, method, args)
	case MethodGetRegisteredPubkey:
		bz, err = p.GetRegisteredPubkey(ctx, contract, method, args)
	}

	if err != nil {
		return nil, err
	}

	cost := ctx.GasMeter().GasConsumed() - initialGas

	if !contract.UseGas(cost) {
		return nil, vm.ErrOutOfGas
	}

	return bz, nil
}

// IsTransaction checks if the given methodID corresponds to a transaction or query.
//
// Available avstask transactions are:
//   - RegisterAVSTask
func (Precompile) IsTransaction(methodID string) bool {
	switch methodID {
	case MethodRegisterAVSTask,
		MethodRegisterBLSPublicKey,
		MethodGetRegisteredPubkey:
		return true
	default:
		return false
	}
}

// Logger returns a precompile-specific logger.
func (p Precompile) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("ExoCore module", "avstask")
}
