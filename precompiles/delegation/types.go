package delegation

import (
	"fmt"
	"math/big"
	"reflect"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/ExocoreNetwork/exocore/precompiles/deposit"
	keeper2 "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

func (p Precompile) GetDelegationParamsFromInputs(ctx sdk.Context, args []interface{}) (*keeper2.DelegationOrUndelegationParams, error) {
	if len(args) != 6 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 6, len(args))
	}
	delegationParams := &keeper2.DelegationOrUndelegationParams{}
	clientChainLzID, ok := args[0].(uint16)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), clientChainLzID)
	}
	delegationParams.ClientChainLzID = uint64(clientChainLzID)

	info, err := p.stakingStateKeeper.GetClientChainInfoByIndex(ctx, delegationParams.ClientChainLzID)
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	txLzNonce, ok := args[1].(uint64)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 1, reflect.TypeOf(args[1]), txLzNonce)
	}
	delegationParams.LzNonce = txLzNonce

	// the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[2].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 2, reflect.TypeOf(args[2]), assetAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(deposit.ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	delegationParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[3].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 3, reflect.TypeOf(args[3]), stakerAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(deposit.ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	delegationParams.StakerAddress = stakerAddr[:clientChainAddrLength]

	// the input operator address is cosmos accAddress type,so we need to check the length and decode it through Bench32
	operatorAddr, ok := args[4].([]byte)
	if !ok || operatorAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 4, reflect.TypeOf(args[4]), operatorAddr)
	}
	if len(operatorAddr) != types.ExoCoreOperatorAddrLength {
		return nil, fmt.Errorf(ErrInputOperatorAddrLength, len(operatorAddr), types.ExoCoreOperatorAddrLength)
	}

	opAccAddr, err := sdk.AccAddressFromBech32(string(operatorAddr))
	if err != nil {
		return nil, errorsmod.Wrap(err, fmt.Sprintf("error occurred when parse acc address from Bech32,the addr is:%s", string(operatorAddr)))
	}
	delegationParams.OperatorAddress = opAccAddr

	opAmount, ok := args[5].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 5, reflect.TypeOf(args[5]), opAmount)
	}
	delegationParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)
	return delegationParams, nil
}
