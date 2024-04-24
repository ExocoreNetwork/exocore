package slash

import (
	"fmt"
	"math/big"
	"reflect"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"

	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/slash/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

func (p Precompile) GetSlashParamsFromInputs(ctx sdk.Context, args []interface{}) (*keeper.SlashParams, error) {
	if len(args) != 8 {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	slashParams := &keeper.SlashParams{}
	clientChainLzID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), clientChainLzID)
	}
	slashParams.ClientChainLzID = uint64(clientChainLzID)

	info, err := p.assetsKeeper.GetClientChainInfoByIndex(ctx, slashParams.ClientChainLzID)
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	// the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, reflect.TypeOf(args[0]), assetAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(exocmn.ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	slashParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, reflect.TypeOf(args[0]), stakerAddr)
	}
	if len(assetAddr) != types.GeneralClientChainAddrLength {
		return nil, fmt.Errorf(exocmn.ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	slashParams.StakerAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, reflect.TypeOf(args[0]), opAmount)
	}

	slashParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)
	return slashParams, nil
}
