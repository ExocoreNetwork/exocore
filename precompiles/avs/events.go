package avs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

const (
	// EventTypeRegisterAVSTask defines the event type for the avs RegisterAVSTask transaction.
	EventTypeRegisterAVSTask = "RegisterAVSTask"
)

// EmitRegisterAVSTaskEvent creates a new event emitted on a SetWithdrawAddressMethod transaction.
func (p Precompile) EmitRegisterAVSTaskEvent(ctx sdk.Context, stateDB vm.StateDB, taskContractAddress string, metaInfo string, name string) error {
	// Prepare the event topics
	event := p.ABI.Events[EventTypeRegisterAVSTask]
	topics := make([]common.Hash, 2)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error
	topics[1], err = cmn.MakeTopic(taskContractAddress)
	if err != nil {
		return err
	}

	// Pack the arguments to be used as the Data field
	arguments := abi.Arguments{event.Inputs[2], event.Inputs[3]}
	packed, err := arguments.Pack(metaInfo, name)
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
