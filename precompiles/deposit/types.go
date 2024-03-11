package deposit

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) GetDepositToParamsFromInputs(ctx sdk.Context, args []interface{}) (*keeper.DepositParams, error) {
	if len(args) != 4 {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	depositParams := &keeper.DepositParams{}
	clientChainLzID, ok := args[0].(uint16)
	if !ok {
		return nil, xerrors.Errorf(ErrContractInputParaOrType, 0, "uint16", clientChainLzID)
	}
	depositParams.ClientChainLzID = uint64(clientChainLzID)

	info, err := p.stakingStateKeeper.GetClientChainInfoByIndex(ctx, depositParams.ClientChainLzID)
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	// the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, xerrors.Errorf(ErrContractInputParaOrType, 1, "[]byte", assetAddr)
	}
	if len(assetAddr) != types.GeneralAssetsAddrLength {
		return nil, xerrors.Errorf(ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	depositParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, xerrors.Errorf(ErrContractInputParaOrType, 2, "[]byte", stakerAddr)
	}
	if len(stakerAddr) != types.GeneralClientChainAddrLength {
		return nil, xerrors.Errorf(ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	depositParams.StakerAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, xerrors.Errorf(ErrContractInputParaOrType, 3, "*big.Int", opAmount)
	}
	depositParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)

	return depositParams, nil
}
