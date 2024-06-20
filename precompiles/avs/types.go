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
	if len(args) != 8 {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, 8, len(args))
	}
	avsParams := &avstypes.AVSRegisterOrDeregisterParams{}
	avsOwnerAddress, ok := args[0].([]string)
	if !ok || avsOwnerAddress == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "[]string", avsOwnerAddress)
	}

	exoAddresses := make([]string, len(avsOwnerAddress))
	var err error
	for i, addr := range avsOwnerAddress {
		exoAddresses[i], err = util.ProcessAddress(addr)
		if err != nil {
			return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "[]string", avsOwnerAddress)
		}
	}

	avsParams.AvsOwnerAddress = exoAddresses

	avsName, ok := args[1].(string)
	if !ok || avsName == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", avsName)
	}
	avsParams.AvsName = avsName

	slashContractAddr, ok := args[2].(string)
	if !ok || slashContractAddr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", slashContractAddr)
	}

	slashContractAddr, err = util.ProcessAddress(slashContractAddr)
	if err != nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", slashContractAddr)
	}
	avsParams.SlashContractAddr = slashContractAddr

	assetID, ok := args[3].([]string)
	if !ok || assetID == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "[]string", assetID)
	}
	avsParams.AssetID = assetID

	action, ok := args[4].(uint64)
	if !ok || (action != avstypes.RegisterAction && action != avstypes.DeRegisterAction) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "uint64", action)
	}
	avsParams.Action = action

	minSelfDelegation, ok := args[5].(uint64)
	if !ok || (action != avstypes.RegisterAction && action != avstypes.DeRegisterAction) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 5, "uint64", minSelfDelegation)
	}
	avsParams.MinSelfDelegation = minSelfDelegation

	unbondingPeriod, ok := args[6].(uint64)
	if !ok || (action != avstypes.RegisterAction && action != avstypes.DeRegisterAction) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 6, "uint64", unbondingPeriod)
	}
	avsParams.UnbondingPeriod = unbondingPeriod

	epochIdentifier, ok := args[7].(string)
	if !ok || (action != avstypes.RegisterAction && action != avstypes.DeRegisterAction) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 7, "string", epochIdentifier)
	}

	avsParams.EpochIdentifier = epochIdentifier

	return avsParams, nil
}
