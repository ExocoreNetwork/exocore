package avs

import (
	avskeep "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v16/precompiles/common"
)

const (
	// EventTypeRegisterAVSTask defines the event type for the avs CreateAVSTask transaction.
	EventTypeRegisterAVSTask = "CreateTask"
)

// EmitCreateAVSTaskEvent creates a new event emitted on a EmitCreateAVSTaskEvent transaction.
func (p Precompile) EmitCreateAVSTaskEvent(ctx sdk.Context, stateDB vm.StateDB, task *avskeep.TaskParams) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeRegisterAVSTask]

	topics := make([]common.Hash, 3)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error
	topics[1], err = cmn.MakeTopic(common.HexToAddress(task.TaskContractAddress))
	if err != nil {
		return err
	}

	topics[2], err = cmn.MakeTopic(task.TaskID)
	if err != nil {
		return err
	}

	// Pack the arguments to be used as the Data field
	arguments := event.Inputs[2:7]
	packed, err := arguments.Pack(task.TaskName, task.Data, task.TaskResponsePeriod, task.TaskChallengePeriod, task.ThresholdPercentage)
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
