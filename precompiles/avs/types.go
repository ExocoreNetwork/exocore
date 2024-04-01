package avs

import (
	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	"github.com/ExocoreNetwork/exocore/x/avs/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) GetAVSParamsFromInputs(ctx sdk.Context, args []interface{}) (*keeper.AVSParams, error) {
	if len(args) != 4 {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	avsParams := &keeper.AVSParams{}
	avsName, ok := args[0].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", avsName)
	}
	avsParams.AVSName = avsName

	avsAddress, ok := args[1].([]byte)
	if !ok || avsAddress == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", avsAddress)
	}
	avsParams.AVSAddress = avsAddress

	operatorAddress, ok := args[2].([]byte)
	if !ok || operatorAddress == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", operatorAddress)
	}
	avsParams.OperatorAddress = operatorAddress

	action, ok := args[3].(uint64)
	if !ok || action != keeper.RegisterAction || action != keeper.DeRegisterAction {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "uint64", action)
	}
	avsParams.Action = action
	return avsParams, nil
}
