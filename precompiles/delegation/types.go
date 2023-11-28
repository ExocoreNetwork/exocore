package delegation

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	keeper2 "github.com/exocore/x/delegation/keeper"
	"math/big"
	"reflect"
)

func GetDelegationParamsFromInputs(args []interface{}) (*keeper2.DelegationOrUnDelegationParams, error) {
	if len(args) != 6 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 6, len(args))
	}
	delegationParams := &keeper2.DelegationOrUnDelegationParams{}
	clientChainLzID, ok := args[0].(uint16)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), clientChainLzID)
	}
	delegationParams.ClientChainLzId = uint64(clientChainLzID)

	txLzNonce, ok := args[1].(uint64)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 1, reflect.TypeOf(args[1]), txLzNonce)
	}
	delegationParams.LzNonce = txLzNonce

	assetAddr, ok := args[2].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 2, reflect.TypeOf(args[2]), assetAddr)
	}
	delegationParams.AssetsAddress = assetAddr

	stakerAddr, ok := args[3].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 3, reflect.TypeOf(args[3]), stakerAddr)
	}
	delegationParams.StakerAddress = stakerAddr

	operatorAddr, ok := args[4].([]byte)
	if !ok || operatorAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 4, reflect.TypeOf(args[4]), operatorAddr)
	}
	opAccAddr, err := sdk.AccAddressFromBech32(string(operatorAddr[:]))
	if err != nil {
		return nil, errorsmod.Wrap(err, "error occurred when parse acc address from Bech32")
	}
	delegationParams.OperatorAddress = opAccAddr

	opAmount, ok := args[5].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 5, reflect.TypeOf(args[5]), opAmount)
	}
	delegationParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)
	return delegationParams, nil
}
