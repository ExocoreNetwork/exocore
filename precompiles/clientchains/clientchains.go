package clientchains

import (
	"bytes"
	"embed"
	"fmt"

	assetskeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
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

// Precompile defines the precompiled contract for deposit.
type Precompile struct {
	cmn.Precompile
	assetsKeeper assetskeeper.Keeper
}

// NewPrecompile creates a new deposit Precompile instance as a
// PrecompiledContract interface.
func NewPrecompile(
	assetsKeeper assetskeeper.Keeper,
	authzKeeper authzkeeper.Keeper,
) (*Precompile, error) {
	abiBz, err := f.ReadFile("abi.json")
	if err != nil {
		return nil, fmt.Errorf("error loading the client chains ABI %s", err)
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
			// should be configurable in the future.
			ApprovalExpiration: cmn.DefaultExpirationDuration,
		},
		assetsKeeper: assetsKeeper,
	}, nil
}

// Address defines the address of the client chains precompile contract.
// address: 0x0000000000000000000000000000000000000801
func (p Precompile) Address() common.Address {
	return common.HexToAddress("0x0000000000000000000000000000000000000801")
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

// Run executes the precompiled contract client chain methods defined in the ABI.
func (p Precompile) Run(
	evm *vm.EVM, contract *vm.Contract, readOnly bool,
) (bz []byte, err error) {
	// if the user calls instead of staticcalls, it is their problem. we don't validate that.
	ctx, _, method, initialGas, args, err := p.RunSetup(
		evm, contract, readOnly, p.IsTransaction,
	)
	if err != nil {
		return nil, err
	}

	// This handles any out of gas errors that may occur during the execution of a precompile tx
	// or query. It avoids panics and returns the out of gas error so the EVM can continue
	// gracefully.
	defer cmn.HandleGasError(ctx, contract, initialGas, &err)()

	bz, err = p.GetClientChains(ctx, method, args)

	if err != nil {
		ctx.Logger().Error(
			"call client chains precompile error",
			"module", "client chains precompile",
			"err", err,
		)
		return nil, err
	}

	cost := ctx.GasMeter().GasConsumed() - initialGas

	if !contract.UseGas(cost) {
		return nil, vm.ErrOutOfGas
	}

	return bz, nil
}

// IsTransaction checks if the given methodID corresponds to a transaction (true)
// or query (false).
func (Precompile) IsTransaction(string) bool {
	// there are no transaction methods in this precompile, only queries.
	return false
}
