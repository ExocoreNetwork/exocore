package avs

import (
	"fmt"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	avskeep "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) GetAVSParamsFromInputs(_ sdk.Context, args []interface{}) (*avstypes.AVSRegisterOrDeregisterParams, error) {
	if len(args) != len(p.ABI.Methods[MethodRegisterAVS].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodRegisterAVS].Inputs), len(args))
	}
	avsParams := &avstypes.AVSRegisterOrDeregisterParams{}
	avsName, ok := args[0].(string)
	if !ok || avsName == "" {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", avsName)
	}
	avsParams.AvsName = avsName
	// When creating tasks in AVS, check the minimum requirements,minStakeAmount at least greater than 0
	minStakeAmount, ok := args[1].(uint64)
	if !ok || minStakeAmount == 0 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "uint64", minStakeAmount)
	}
	avsParams.MinStakeAmount = minStakeAmount

	taskAddr, ok := args[2].(common.Address)
	if !ok || taskAddr == (common.Address{}) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "common.Address", taskAddr)
	}
	avsParams.TaskAddr = taskAddr.String()

	slashContractAddr, ok := args[3].(common.Address)
	if !ok || (slashContractAddr == common.Address{}) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "common.Address", slashContractAddr)
	}
	avsParams.SlashContractAddr = slashContractAddr.String()

	rewardContractAddr, ok := args[4].(common.Address)
	if !ok || (rewardContractAddr == common.Address{}) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "common.Address", rewardContractAddr)
	}
	avsParams.RewardContractAddr = rewardContractAddr.String()

	// bech32
	avsOwnerAddress, ok := args[5].([]string)
	if !ok || avsOwnerAddress == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 5, "[]string", avsOwnerAddress)
	}
	exoAddresses := make([]string, len(avsOwnerAddress))
	for i, addr := range avsOwnerAddress {
		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 5, "[]string", avsOwnerAddress)
		}
		exoAddresses[i] = accAddr.String()
	}
	avsParams.AvsOwnerAddress = exoAddresses

	// string, since it is the address_id representation
	assetID, ok := args[6].([]string)
	if !ok || assetID == nil || len(assetID) == 0 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 6, "[]string", assetID)
	}
	avsParams.AssetID = assetID

	unbondingPeriod, ok := args[7].(uint64)
	if !ok || unbondingPeriod == 0 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 7, "uint64", unbondingPeriod)
	}
	avsParams.UnbondingPeriod = unbondingPeriod

	minSelfDelegation, ok := args[8].(uint64)
	if !ok || minSelfDelegation == 0 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 8, "uint64", minSelfDelegation)
	}
	avsParams.MinSelfDelegation = minSelfDelegation

	epochIdentifier, ok := args[9].(string)
	if !ok || epochIdentifier == "" {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 9, "string", epochIdentifier)
	}
	avsParams.EpochIdentifier = epochIdentifier

	// The parameters below are used when creating tasks, to ensure that the minimum criteria are met by the set
	// of operators.

	taskParam, ok := args[10].([]uint64)
	if !ok || taskParam == nil || len(taskParam) != 4 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 10, "[]string", taskParam)
	}
	minOptInOperators := taskParam[0]
	avsParams.MinOptInOperators = minOptInOperators

	minTotalStakeAmount := taskParam[1]
	avsParams.MinTotalStakeAmount = minTotalStakeAmount

	avsReward := taskParam[2]
	avsParams.AvsReward = avsReward

	avsSlash := taskParam[3]
	avsParams.AvsSlash = avsSlash

	return avsParams, nil
}

func (p Precompile) GetAVSParamsFromUpdateInputs(_ sdk.Context, args []interface{}) (*avstypes.AVSRegisterOrDeregisterParams, error) {
	if len(args) != len(p.ABI.Methods[MethodRegisterAVS].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, len(p.ABI.Methods[MethodRegisterAVS].Inputs), len(args))
	}
	avsParams := &avstypes.AVSRegisterOrDeregisterParams{}
	avsName, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", avsName)
	}
	avsParams.AvsName = avsName
	// When creating tasks in AVS, check the minimum requirements,minStakeAmount at least greater than 0
	minStakeAmount, ok := args[1].(uint64)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "uint64", minStakeAmount)
	}
	avsParams.MinStakeAmount = minStakeAmount

	taskAddr, ok := args[2].(common.Address)
	if !ok || taskAddr == (common.Address{}) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "common.Address", taskAddr)
	}
	avsParams.TaskAddr = taskAddr.String()

	slashContractAddr, ok := args[3].(common.Address)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "common.Address", slashContractAddr)
	}
	avsParams.SlashContractAddr = slashContractAddr.String()

	rewardContractAddr, ok := args[4].(common.Address)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "common.Address", rewardContractAddr)
	}
	avsParams.RewardContractAddr = rewardContractAddr.String()

	// bech32
	avsOwnerAddress, ok := args[5].([]string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 5, "[]string", avsOwnerAddress)
	}
	exoAddresses := make([]string, len(avsOwnerAddress))
	for i, addr := range avsOwnerAddress {
		accAddr, err := sdk.AccAddressFromBech32(addr)
		if err != nil {
			return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 5, "[]string", avsOwnerAddress)
		}
		exoAddresses[i] = accAddr.String()
	}
	avsParams.AvsOwnerAddress = exoAddresses

	// string, since it is the address_id representation
	assetID, ok := args[6].([]string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 6, "[]string", assetID)
	}
	avsParams.AssetID = assetID

	unbondingPeriod, ok := args[7].(uint64)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 7, "uint64", unbondingPeriod)
	}
	avsParams.UnbondingPeriod = unbondingPeriod

	minSelfDelegation, ok := args[8].(uint64)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 8, "uint64", minSelfDelegation)
	}
	avsParams.MinSelfDelegation = minSelfDelegation

	epochIdentifier, ok := args[9].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 9, "string", epochIdentifier)
	}
	avsParams.EpochIdentifier = epochIdentifier

	// The parameters below are used when creating tasks, to ensure that the minimum criteria are met by the set
	// of operators.

	taskParam, ok := args[10].([]uint64)
	if !ok || taskParam == nil || len(taskParam) != 4 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 10, "[]string", taskParam)
	}
	minOptInOperators := taskParam[0]
	avsParams.MinOptInOperators = minOptInOperators

	minTotalStakeAmount := taskParam[1]
	avsParams.MinTotalStakeAmount = minTotalStakeAmount

	avsReward := taskParam[2]
	avsParams.AvsReward = avsReward

	avsSlash := taskParam[3]
	avsParams.AvsSlash = avsSlash

	return avsParams, nil
}

func (p Precompile) GetTaskParamsFromInputs(_ sdk.Context, args []interface{}) (*avskeep.TaskParams, error) {
	if len(args) != len(p.ABI.Methods[MethodCreateAVSTask].Inputs) {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}
	taskParams := &avskeep.TaskParams{}

	name, ok := args[0].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", name)
	}
	taskParams.TaskName = name

	data, ok := args[1].([]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", data)
	}
	taskParams.Hash = data

	taskID, ok := args[2].(string)
	if !ok || taskID == "" {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", taskID)
	}
	taskParams.TaskID = taskID

	taskResponsePeriod, ok := args[3].(uint64)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, "uint64", taskResponsePeriod)
	}
	taskParams.TaskResponsePeriod = taskResponsePeriod

	taskChallengePeriod, ok := args[4].(uint64)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 4, "uint64", taskChallengePeriod)
	}
	taskParams.TaskChallengePeriod = taskChallengePeriod

	thresholdPercentage, ok := args[5].(uint64)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 5, "uint64", thresholdPercentage)
	}
	taskParams.ThresholdPercentage = thresholdPercentage

	return taskParams, nil
}
