package avs

import (
	"encoding/hex"
	avskeep "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	avstypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

const (
	EventTypeAVSRegistered               = "AVSRegistered"
	EventTypeAVSUpdated                  = "AVSUpdated"
	EventTypeAVSDeregistered             = "AVSDeregistered"
	EventTypeOperatorJoined              = "OperatorJoined"
	EventTypeOperatorOuted               = "OperatorOuted"
	EventTypeTaskCreated                 = "TaskCreated"
	EventTypeChallengeInitiated          = "ChallengeInitiated"
	EventTypePublicKeyRegistered         = "PublicKeyRegistered"
	EventTypeOperatorRegisteredToExocore = "OperatorRegisteredToExocore"
	EventTypeTaskSubmittedByOperator     = "TaskSubmittedByOperator"
)

func (p Precompile) EmitAVSRegistered(ctx sdk.Context, stateDB vm.StateDB, avs *avstypes.AVSRegisterOrDeregisterParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeAVSRegistered]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:2]
	packed, err := arguments.Pack(
		common.HexToAddress(avs.CallerAddress),
		avs.AvsName,
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (p Precompile) EmitAVSUpdated(ctx sdk.Context, stateDB vm.StateDB, avs *avstypes.AVSRegisterOrDeregisterParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeAVSUpdated]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:2]
	packed, err := arguments.Pack(
		common.HexToAddress(avs.CallerAddress),
		avs.AvsName,
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (p Precompile) EmitAVSDeregistered(ctx sdk.Context, stateDB vm.StateDB, avs *avstypes.AVSRegisterOrDeregisterParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeAVSDeregistered]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:2]
	packed, err := arguments.Pack(
		common.HexToAddress(avs.CallerAddress),
		avs.AvsName,
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (p Precompile) EmitOperatorJoined(ctx sdk.Context, stateDB vm.StateDB, params *avskeep.OperatorOptParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeOperatorJoined]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:1]
	packed, err := arguments.Pack(
		common.HexToAddress(params.OperatorAddress),
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (p Precompile) EmitOperatorOuted(ctx sdk.Context, stateDB vm.StateDB, params *avskeep.OperatorOptParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeOperatorOuted]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:1]
	packed, err := arguments.Pack(
		common.HexToAddress(params.OperatorAddress),
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (p Precompile) EmitTaskCreated(ctx sdk.Context, stateDB vm.StateDB, task *avskeep.TaskInfoParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeTaskCreated]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:9]
	packed, err := arguments.Pack(
		common.HexToAddress(task.CallerAddress),
		task.TaskID,
		common.HexToAddress(task.TaskContractAddress),
		task.TaskName,
		task.Hash,
		task.TaskResponsePeriod,
		task.TaskChallengePeriod,
		task.ThresholdPercentage,
		task.TaskStatisticalPeriod)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (p Precompile) EmitChallengeInitiated(ctx sdk.Context, stateDB vm.StateDB, params *avskeep.ChallengeParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeChallengeInitiated]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:5]
	packed, err := arguments.Pack(
		common.HexToAddress(params.CallerAddress),
		params.TaskHash,
		params.TaskID,
		params.TaskResponseHash,
		params.OperatorAddress,
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}
func (p Precompile) EmitPublicKeyRegistered(ctx sdk.Context, stateDB vm.StateDB, params *avskeep.BlsParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypePublicKeyRegistered]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:2]
	packed, err := arguments.Pack(
		common.HexToAddress(params.Operator),
		params.Name,
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (p Precompile) EmitOperatorRegisteredToExocore(ctx sdk.Context, stateDB vm.StateDB, params *avskeep.OperatorParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeOperatorRegisteredToExocore]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:2]
	packed, err := arguments.Pack(
		common.HexToAddress(params.CallerAddress),
		params.OperatorMetaInfo,
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}

func (p Precompile) EmitTaskSubmittedByOperator(ctx sdk.Context, stateDB vm.StateDB, params *avskeep.TaskResultParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeTaskSubmittedByOperator]

	topics := make([]common.Hash, 1)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[0:6]
	packed, err := arguments.Pack(
		common.HexToAddress(params.CallerAddress),
		params.TaskID,
		hex.EncodeToString(params.TaskResponse),
		hex.EncodeToString(params.BlsSignature),
		params.TaskContractAddress.String(),
		params.Stage,
		true)
	if err != nil {
		return err
	}

	stateDB.AddLog(&ethtypes.Log{
		Address:     p.Address(),
		Topics:      topics,
		Data:        packed,
		BlockNumber: uint64(ctx.BlockHeight()),
	})

	return nil
}
