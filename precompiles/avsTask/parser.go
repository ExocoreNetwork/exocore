package task

import (
	"fmt"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	types "github.com/ExocoreNetwork/exocore/x/avstask/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) GetTaskParamsFromInputs(_ sdk.Context, args []interface{}) (*types.RegisterAVSTaskReq, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 3, len(args))
	}
	taskParams := &types.RegisterAVSTaskReq{}
	taskinfo := &types.TaskContractInfo{}

	taskaddr, ok := args[0].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "string", taskaddr)
	}
	taskinfo.TaskContractAddress = taskaddr

	taskName, ok := args[1].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "string", taskName)
	}
	taskinfo.Name = taskName

	metainfo, ok := args[2].(string)
	if !ok || metainfo == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", metainfo)
	}
	taskinfo.MetaInfo = metainfo
	taskParams.Task = taskinfo
	return taskParams, nil
}
