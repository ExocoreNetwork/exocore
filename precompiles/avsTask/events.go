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
	// EventTypeNewTaskCreated defines the event type for the avstask create transaction.
	EventTypeNewPubkeyRegistration = "NewPubkeyRegistration"
)

// EmitEventTypeNewPubkeyRegistration new bls pubkey reg
func (p Precompile) EmitEventTypeNewPubkeyRegistration(
	ctx sdk.Context,
	stateDB vm.StateDB,
	operator string,
	pubkey []byte,
) error {
	event := p.ABI.Events[EventTypeNewPubkeyRegistration]
	topics := make([]common.Hash, 2)

	// The first topic is always the signature of the event.
	topics[0] = event.ID

	var err error
	// sender and receiver are indexed
	topics[1], err = cmn.MakeTopic(operator)
	if err != nil {
		return err
	}

	// Pack the arguments to be used as the Data field
	arguments := abi.Arguments{event.Inputs[0], event.Inputs[1]}
	packed, err := arguments.Pack(operator, pubkey)
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
