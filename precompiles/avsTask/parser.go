package task

import (
	"fmt"
	"reflect"

	"github.com/ExocoreNetwork/exocore/x/taskmanageravs/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

func (p Precompile) GetTaskParamsFromInputs(ctx sdk.Context, args []interface{}) (*keeper.CreateNewTaskParams, error) {
	if len(args) != 8 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	taskParams := &keeper.CreateNewTaskParams{}
	numberToBeSquared, ok := args[0].(uint16)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), numberToBeSquared)
	}
	taskParams.NumberToBeSquared = uint64(numberToBeSquared)
	taskParams.QuorumThresholdPercentage = args[1].(uint32)
	qnums, ok := args[2].([]byte)
	taskParams.QuorumNumbers = qnums

	return taskParams, nil
}
