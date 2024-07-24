package avs

import (
	"fmt"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	util "github.com/ExocoreNetwork/exocore/utils"
	avskeep "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) GetAVSParamsFromInputs(_ sdk.Context, args []interface{}) (*avskeep.AVSRegisterOrDeregisterParams, error) {
	if len(args) != len(p.ABI.Methods[MethodRegisterAVS].Inputs) {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodRegisterAVS].Inputs), len(args))
	}
	avsParams := &avskeep.AVSRegisterOrDeregisterParams{}
	var err error
	avsName, ok := args[0].(string)
	if !ok || avsName == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", avsName)
	}
	avsParams.AvsName = avsName
	// When creating tasks in AVS, check the minimum requirements,minStakeAmount at least greater than 0
	minStakeAmount, ok := args[1].(uint64)
	if !ok || minStakeAmount == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "uint64", minStakeAmount)
	}
	avsParams.MinStakeAmount = minStakeAmount

	taskAddr, ok := args[2].(string)
	if !ok || taskAddr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", taskAddr)
	}
	if err != nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", taskAddr)
	}
	avsParams.TaskAddr = taskAddr

	slashContractAddr, ok := args[3].(string)
	if !ok || slashContractAddr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "string", slashContractAddr)
	}

	slashContractAddr, err = util.ProcessAddress(slashContractAddr)
	if err != nil || slashContractAddr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "string", slashContractAddr)
	}
	avsParams.SlashContractAddr = slashContractAddr

	rewardContractAddr, ok := args[4].(string)
	if !ok || rewardContractAddr == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "string", rewardContractAddr)
	}

	rewardContractAddr, err = util.ProcessAddress(rewardContractAddr)
	if err != nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "string", rewardContractAddr)
	}
	avsParams.RewardContractAddr = rewardContractAddr

	avsOwnerAddress, ok := args[5].([]string)
	if !ok || avsOwnerAddress == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 5, "[]string", avsOwnerAddress)
	}

	exoAddresses := make([]string, len(avsOwnerAddress))

	for i, addr := range avsOwnerAddress {
		exoAddresses[i], err = util.ProcessAddress(addr)
		if err != nil {
			return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 5, "[]string", avsOwnerAddress)
		}
	}

	avsParams.AvsOwnerAddress = exoAddresses

	assetID, ok := args[6].([]string)
	if !ok || assetID == nil || len(assetID) == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 6, "[]string", assetID)
	}
	avsParams.AssetID = assetID

	unbondingPeriod, ok := args[7].(uint64)
	if !ok || unbondingPeriod == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 7, "uint64", unbondingPeriod)
	}
	avsParams.UnbondingPeriod = unbondingPeriod

	minSelfDelegation, ok := args[8].(uint64)
	if !ok || minSelfDelegation == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 8, "uint64", minSelfDelegation)
	}
	avsParams.MinSelfDelegation = minSelfDelegation

	epochIdentifier, ok := args[9].(string)
	if !ok || epochIdentifier == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 9, "string", epochIdentifier)
	}

	avsParams.EpochIdentifier = epochIdentifier
	// When creating tasks in AVS, check the minimum requirements,minOptInOperators at least greater than 0
	minOptInOperators, ok := args[10].(uint64)
	if !ok || minOptInOperators == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 10, "uint64", minOptInOperators)
	}
	avsParams.MinOptInOperators = minOptInOperators
	// When creating tasks in AVS, check the minimum requirements,minTotalStakeAmount at least greater than 0
	minTotalStakeAmount, ok := args[11].(uint64)
	if !ok || minTotalStakeAmount == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 11, "uint64", minTotalStakeAmount)
	}
	avsParams.MinTotalStakeAmount = minTotalStakeAmount

	avsReward, ok := args[12].(uint64)
	if !ok || avsReward == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 12, "uint64", avsReward)
	}
	avsParams.AvsReward = avsReward

	avsSlash, ok := args[13].(uint64)
	if !ok || avsSlash == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 13, "uint64", avsSlash)
	}
	avsParams.AvsSlash = avsSlash

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
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 7, "string", epochIdentifier)
	}

	avsParams.EpochIdentifier = epochIdentifier

	return avsParams, nil
}

func (p Precompile) GetTaskParamsFromInputs(_ sdk.Context, args []interface{}) (*avskeep.TaskParams, error) {
	if len(args) != len(p.ABI.Methods[MethodCreateAVSTask].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}
	taskParams := &avskeep.TaskParams{}

	name, ok := args[0].(string)
	if !ok || name == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", name)
	}
	taskParams.TaskName = name

	data, ok := args[1].([]byte)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", data)
	}
	taskParams.Data = data

	taskID, ok := args[2].(string)
	if !ok || taskID == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", taskID)
	}
	taskParams.TaskID = taskID

	taskResponsePeriod, ok := args[3].(uint64)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "uint64", taskResponsePeriod)
	}
	taskParams.TaskResponsePeriod = taskResponsePeriod

	taskChallengePeriod, ok := args[4].(uint64)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "uint64", taskChallengePeriod)
	}
	taskParams.TaskChallengePeriod = taskChallengePeriod

	thresholdPercentage, ok := args[5].(uint64)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 5, "uint64", thresholdPercentage)
	}
	taskParams.ThresholdPercentage = thresholdPercentage
	return taskParams, nil
}
