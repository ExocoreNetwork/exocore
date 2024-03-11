package deposit

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	// MethodDepositTo defines the ABI method name for the deposit
	// DepositTo transaction.
	MethodDepositTo = "depositTo"
)

// DepositTo deposit the client chain assets to the staker, that will change the state in deposit module.
func (p Precompile) DepositTo(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract,the caller must be exoCore LzApp contract
	depositModuleParam, err := p.depositKeeper.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	exoCoreLzAppAddr := common.HexToAddress(depositModuleParam.ExoCoreLzAppAddress)
	if contract.CallerAddress != exoCoreLzAppAddr {
		return nil, fmt.Errorf(ErrContractCaller, contract.CallerAddress, exoCoreLzAppAddr)
	}

	// parse the depositTo input params
	depositParams, err := p.GetDepositToParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	// call depositKeeper to execute the deposit action
	err = p.depositKeeper.Deposit(ctx, depositParams)
	if err != nil {
		return nil, err
	}

	// get the latest asset state of staker to return.
	stakerID, assetID := types.GetStakeIDAndAssetID(depositParams.ClientChainLzID, depositParams.StakerAddress, depositParams.AssetsAddress)
	info, err := p.stakingStateKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetID)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true, info.TotalDepositAmount.BigInt())
}
