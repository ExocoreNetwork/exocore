package assets

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	assetskeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) GetDepositWithdrawParamsFromInputs(ctx sdk.Context, args []interface{}) (*assetskeeper.DepositWithdrawParams, error) {
	if len(args) != 4 {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	depositWithdrawParams := &assetskeeper.DepositWithdrawParams{}
	clientChainLzID, ok := args[0].(uint32)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint16", clientChainLzID)
	}
	depositWithdrawParams.ClientChainLzID = uint64(clientChainLzID)

	info, err := p.assetsKeeper.GetClientChainInfoByIndex(ctx, depositWithdrawParams.ClientChainLzID)
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	// the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", assetAddr)
	}
	if len(assetAddr) != types.GeneralAssetsAddrLength {
		return nil, xerrors.Errorf(exocmn.ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	depositWithdrawParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", stakerAddr)
	}
	if len(stakerAddr) != types.GeneralClientChainAddrLength {
		return nil, xerrors.Errorf(exocmn.ErrInputClientChainAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	depositWithdrawParams.StakerAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "*big.Int", opAmount)
	}
	depositWithdrawParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)

	return depositWithdrawParams, nil
}
