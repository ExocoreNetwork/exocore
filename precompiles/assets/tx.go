package assets

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

const (
	// MethodDepositTo defines the ABI method name for the deposit
	// DepositAndWithdraw transaction.
	MethodDepositTo = "depositTo"
	MethodWithdraw  = "withdrawPrinciple"

	MethodGetClientChains = "getClientChains"
)

// DepositAndWithdraw deposit and withdraw the client chain assets for the staker,
// that will change the state in assets module.
func (p Precompile) DepositAndWithdraw(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract,the caller must be exoCore LzApp contract
	err := p.assetsKeeper.CheckExocoreLzAppAddr(ctx, contract.CallerAddress)
	if err != nil {
		return nil, errorsmod.Wrap(err, exocmn.ErrContractCaller)
	}
	// parse the depositTo input params
	depositWithdrawParams, err := p.GetDepositWithdrawParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	// call assets keeper to perform the deposit or withdraw action
	switch method.Name {
	// deposit transactions
	case MethodDepositTo:
		depositWithdrawParams.Action = assetstypes.Deposit
	case MethodWithdraw:
		depositWithdrawParams.Action = assetstypes.WithdrawPrinciple
	default:
		return nil, fmt.Errorf(cmn.ErrUnknownMethod, method.Name)
	}
	err = p.assetsKeeper.PerformDepositOrWithdraw(ctx, depositWithdrawParams)
	if err != nil {
		return nil, err
	}

	// get the latest asset state of staker to return.
	stakerID, assetID := assetstypes.GetStakeIDAndAssetID(depositWithdrawParams.ClientChainLzID, depositWithdrawParams.StakerAddress, depositWithdrawParams.AssetsAddress)
	info, err := p.assetsKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetID)
	if err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true, info.TotalDepositAmount.BigInt())
}

func (p Precompile) GetClientChains(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) > 0 {
		ctx.Logger().Error(
			"GetClientChains",
			"err", errorsmod.Wrapf(
				assetstypes.ErrInvalidInputParameter, "no input is required",
			),
		)
		return method.Outputs.Pack(false, nil)
	}
	infos, err := p.assetsKeeper.GetAllClientChainInfo(ctx)
	if err != nil {
		ctx.Logger().Error(
			"GetClientChains",
			"err", err,
		)
		return method.Outputs.Pack(false, nil)
	}
	ids := make([]uint32, 0, len(infos))
	for id := range infos {
		// #nosec G701 // already checked
		convID := uint32(id)
		ids = append(ids, convID)
	}
	return method.Outputs.Pack(true, ids)
}
