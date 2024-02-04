package deposit

import (
	"fmt"
	"math/big"
	"reflect"

	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

func (p Precompile) GetDepositToParamsFromInputs(ctx sdk.Context, args []interface{}) (*keeper.DepositParams, error) {
	if len(args) != 4 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	depositParams := &keeper.DepositParams{}
	clientChainLzID, ok := args[0].(uint16)
	if !ok {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), clientChainLzID)
	}
	depositParams.ClientChainLzId = uint64(clientChainLzID)

	info, err := p.stakingStateKeeper.GetClientChainInfoByIndex(ctx, depositParams.ClientChainLzId)
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	// the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 1, reflect.TypeOf(args[0]), assetAddr)
	}
	if len(assetAddr) != types.GeneralAssetsAddrLength {
		return nil, fmt.Errorf(ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	depositParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 2, reflect.TypeOf(args[0]), stakerAddr)
	}
	if len(stakerAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	depositParams.StakerAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf(ErrContractInputParaOrType, 3, reflect.TypeOf(args[0]), opAmount)
	}
	depositParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)

	return depositParams, nil
}
