package deposit

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"github.com/exocore/x/deposit/keeper"
	"math/big"
	"reflect"
)

func GetDepositToParamsFromInputs(args []interface{}) (*keeper.DepositParams, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	depositParams := &keeper.DepositParams{}
	clientChainLzID, ok := args[0].(uint16)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), clientChainLzID)
	}
	depositParams.ClientChainLzId = uint64(clientChainLzID)

	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 1, reflect.TypeOf(args[0]), assetAddr)
	}
	depositParams.AssetsAddress = assetAddr

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 2, reflect.TypeOf(args[0]), stakerAddr)
	}
	depositParams.StakerAddress = stakerAddr

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 3, reflect.TypeOf(args[0]), opAmount)
	}
	depositParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)
	return depositParams, nil
}
