package avs

import (
	"fmt"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	util "github.com/ExocoreNetwork/exocore/utils"
	avskeep "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) GetAVSParamsFromInputs(_ sdk.Context, args []interface{}) (*avskeep.AVSRegisterOrDeregisterParams, error) {
	if len(args) != len(p.ABI.Methods[MethodRegisterAVS].Inputs) {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodRegisterAVS].Inputs), len(args))
	}
	avsParams := &avskeep.AVSRegisterOrDeregisterParams{}
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

	rewardContractAddr, ok := args[2].(string)
	if !ok || rewardContractAddr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", rewardContractAddr)
	}

	rewardContractAddr, err = util.ProcessAddress(rewardContractAddr)
	if err != nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", rewardContractAddr)
	}
	avsParams.RewardContractAddr = rewardContractAddr

	slashContractAddr, ok := args[3].(string)
	if !ok || slashContractAddr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "string", slashContractAddr)
	}

	slashContractAddr, err = util.ProcessAddress(slashContractAddr)
	if err != nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "string", slashContractAddr)
	}
	avsParams.SlashContractAddr = slashContractAddr

	assetID, ok := args[4].([]string)
	if !ok || assetID == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "[]string", assetID)
	}
	avsParams.AssetID = assetID

	minSelfDelegation, ok := args[5].(uint64)
	if !ok || minSelfDelegation == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 5, "uint64", minSelfDelegation)
	}
	avsParams.MinSelfDelegation = minSelfDelegation

	unbondingPeriod, ok := args[6].(uint64)
	if !ok || unbondingPeriod == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 6, "uint64", unbondingPeriod)
	}
	avsParams.UnbondingPeriod = unbondingPeriod

	epochIdentifier, ok := args[7].(string)
	if !ok || epochIdentifier == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 7, "string", epochIdentifier)
	}

	avsParams.EpochIdentifier = epochIdentifier

	return avsParams, nil
}

func (p Precompile) GetAVSParamsFromUpdateInputs(_ sdk.Context, args []interface{}) (*avskeep.AVSRegisterOrDeregisterParams, error) {
	if len(args) != len(p.ABI.Methods[MethodUpdateAVS].Inputs) {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodUpdateAVS].Inputs), len(args))
	}
	avsParams := &avskeep.AVSRegisterOrDeregisterParams{}
	avsOwnerAddress, ok := args[0].([]string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "[]string", avsOwnerAddress)
	}
	avsParams.AvsOwnerAddress = nil
	var err error
	if avsOwnerAddress != nil {
		exoAddresses := make([]string, len(avsOwnerAddress))

		for i, addr := range avsOwnerAddress {
			exoAddresses[i], err = util.ProcessAddress(addr)
			if err != nil {
				return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "[]string", avsOwnerAddress)
			}
		}
		avsParams.AvsOwnerAddress = exoAddresses
	}

	avsName, ok := args[1].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", avsName)
	}
	avsParams.AvsName = avsName

	rewardContractAddr, ok := args[2].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", rewardContractAddr)
	}
	if rewardContractAddr != "" {
		rewardContractAddr, err = util.ProcessAddress(rewardContractAddr)
		if err != nil {
			return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", rewardContractAddr)
		}
	}
	avsParams.RewardContractAddr = rewardContractAddr

	slashContractAddr, ok := args[3].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "string", slashContractAddr)
	}
	if slashContractAddr != "" {
		slashContractAddr, err = util.ProcessAddress(slashContractAddr)
		if err != nil {
			return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "string", slashContractAddr)
		}
	}
	avsParams.SlashContractAddr = slashContractAddr

	assetID, ok := args[4].([]string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "[]string", assetID)
	}
	avsParams.AssetID = assetID

	minSelfDelegation, ok := args[5].(uint64)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 5, "uint64", minSelfDelegation)
	}
	avsParams.MinSelfDelegation = minSelfDelegation

	unbondingPeriod, ok := args[6].(uint64)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 6, "uint64", unbondingPeriod)
	}
	avsParams.UnbondingPeriod = unbondingPeriod

	epochIdentifier, ok := args[7].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 7, "string", epochIdentifier)
	}

	avsParams.EpochIdentifier = epochIdentifier

	return avsParams, nil
}

func (p Precompile) GetTaskParamsFromInputs(_ sdk.Context, args []interface{}) (*avstypes.RegisterAVSTaskReq, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}
	taskParams := &avstypes.RegisterAVSTaskReq{}
	taskinfo := &avstypes.TaskInfo{}

	taskaddr, ok := args[0].(string)
	if !ok || taskaddr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", taskaddr)
	}
	taskinfo.TaskContractAddress = taskaddr

	taskName, ok := args[1].(string)
	if !ok || taskName == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", taskName)
	}
	taskinfo.Name = taskName

	name, ok := args[2].(string)
	if !ok || name == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", name)
	}
	taskinfo.Name = name
	taskParams.Task = taskinfo
	return taskParams, nil
}
