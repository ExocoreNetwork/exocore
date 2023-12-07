package slash

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"github.com/exocore/x/exoslash/keeper"
	"math/big"
	"reflect"
)

func GetSlashParamsFromInputs(args []interface{}) (*keeper.SlashParams, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	slashParams := &keeper.SlashParams{}
	clientChainLzID, ok := args[0].(uint16)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), clientChainLzID)
	}
	slashParams.ClientChainLzId = uint64(clientChainLzID)

	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 1, reflect.TypeOf(args[0]), assetAddr)
	}
	slashParams.AssetsAddress = assetAddr

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 2, reflect.TypeOf(args[0]), stakerAddr)
	}
	slashParams.OperatorAddress = stakerAddr

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 3, reflect.TypeOf(args[0]), opAmount)
	}
	slashParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)
	return slashParams, nil
}
