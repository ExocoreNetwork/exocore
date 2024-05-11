package avs

import (
	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	util "github.com/ExocoreNetwork/exocore/utils"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) GetAVSParamsFromInputs(_ sdk.Context, args []interface{}) (*avstypes.AVSRegisterOrDeregisterParams, error) {
	if len(args) != 6 {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, 6, len(args))
	}
	avsParams := &avstypes.AVSRegisterOrDeregisterParams{}
	avsName, ok := args[0].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", avsName)
	}
	avsParams.AvsName = avsName

	avsAddress, ok := args[1].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", avsAddress)
	}
	avsAddress, err := util.ProcessAvsAddress(avsAddress)
	if err != nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", avsAddress)
	}
	avsParams.AvsAddress = avsAddress

	operatorAddress, ok := args[2].(string)
	if !ok || operatorAddress == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", operatorAddress)
	}

	operatorAddress, err = util.ProcessAvsAddress(operatorAddress)
	if err != nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", operatorAddress)
	}
	avsParams.OperatorAddress = operatorAddress
	action, ok := args[3].(uint64)
	if !ok || (action != avstypes.RegisterAction && action != avstypes.DeRegisterAction) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "uint64", action)
	}
	avsParams.Action = action

	avsOwnerAddress, ok := args[4].(string)
	if !ok || avsOwnerAddress == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "string", avsOwnerAddress)
	}
	avsOwnerAddress, err = util.ProcessAvsAddress(avsOwnerAddress)
	if err != nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", avsOwnerAddress)
	}
	avsParams.AvsOwnerAddress = avsOwnerAddress

	assetID, ok := args[5].(string)
	if !ok || assetID == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "uint64", action)
	}
	avsParams.AssetID = assetID
	return avsParams, nil
}
