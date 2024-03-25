package task

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

const (
	// EventTypeNewTaskCreated defines the event type for the task create transaction.
	EventTypeNewTaskCreated = "NewTaskCreated"
)

// EmitNewTaskCreatedEvent creates a new task transaction.
func (p Precompile) EmitNewTaskCreatedEvent(
	ctx sdk.Context,
	stateDB vm.StateDB,
	taskIndex uint32,
	numberToBeSquared uint64,
	quorumNumbers []byte,
	quorumThresholdPercentage uint32,
) error {
	event := p.ABI.Events[EventTypeNewTaskCreated]
	topics := make([]common.Hash, 3)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error
	// sender and receiver are indexed
	topics[1], err = cmn.MakeTopic(taskIndex)
	if err != nil {
		return err
	}

	// Prepare the event data: denom, amount, memo
	arguments := abi.Arguments{event.Inputs[2], event.Inputs[3], event.Inputs[4]}
	packed, err := arguments.Pack(numberToBeSquared, quorumNumbers, quorumThresholdPercentage)
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
