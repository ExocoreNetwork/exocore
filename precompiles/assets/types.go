package assets

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"

	sdkmath "cosmossdk.io/math"
	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	assetskeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	"golang.org/x/xerrors"
)

func (p Precompile) DepositWithdrawParamsFromInputs(ctx sdk.Context, args []interface{}) (*assetskeeper.DepositWithdrawParams, error) {
	if len(args) != 4 {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	depositWithdrawParams := &assetskeeper.DepositWithdrawParams{}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), args[0])
	}
	depositWithdrawParams.ClientChainLzID = uint64(clientChainID)

	info, err := p.assetsKeeper.GetClientChainInfoByIndex(ctx, depositWithdrawParams.ClientChainLzID)
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	// the length of client chain address inputted by caller is 32, so we need to check the length and remove the padding according to the actual length.
	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, reflect.TypeOf(args[1]), args[1])
	}
	if len(assetAddr) != types.GeneralAssetsAddrLength {
		return nil, xerrors.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	depositWithdrawParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, reflect.TypeOf(args[2]), args[2])
	}
	if len(stakerAddr) != types.GeneralClientChainAddrLength {
		return nil, xerrors.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddr), types.GeneralClientChainAddrLength)
	}
	depositWithdrawParams.StakerAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || opAmount.Cmp(big.NewInt(0)) == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, reflect.TypeOf(args[3]), args[3])
	}
	depositWithdrawParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)

	return depositWithdrawParams, nil
}

func (p Precompile) ClientChainInfoFromInputs(_ sdk.Context, args []interface{}) (*types.ClientChainInfo, error) {
	if len(args) != 5 {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, 5, len(args))
	}
	clientChain := types.ClientChainInfo{}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), args[0])
	}
	clientChain.LayerZeroChainID = uint64(clientChainID)

	addressLength, ok := args[1].(uint32)
	if !ok || addressLength == 0 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, reflect.TypeOf(args[1]), args[1])
	}
	if addressLength > types.GeneralAssetsAddrLength {
		return nil, xerrors.Errorf(exocmn.ErrInvalidAddrLength, addressLength, "not greater than 32")
	}
	clientChain.AddressLength = addressLength

	name, ok := args[2].(string)
	if !ok || name == "" {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, reflect.TypeOf(args[2]), args[2])
	}
	clientChain.Name = name

	metaInfo, ok := args[3].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, reflect.TypeOf(args[2]), args[2])
	}
	clientChain.MetaInfo = metaInfo

	signatureType, ok := args[4].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, reflect.TypeOf(args[4]), args[4])
	}
	clientChain.SignatureType = signatureType

	return &clientChain, nil
}

func (p Precompile) TokensFromInputs(ctx sdk.Context, args []interface{}) ([]types.AssetInfo, error) {
	if len(args) != 4 {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, 4, len(args))
	}
	assets := make([]types.AssetInfo, 0)
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, reflect.TypeOf(args[0]), args[0])
	}
	info, err := p.assetsKeeper.GetClientChainInfoByIndex(ctx, uint64(clientChainID))
	if err != nil {
		return nil, err
	}
	clientChainAddrLength := info.AddressLength

	assetAddrList, ok := args[1].([][]byte)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, reflect.TypeOf(args[1]), args[1])
	}
	assetNumber := len(assetAddrList)
	if assetNumber == 0 {
		return nil, fmt.Errorf(exocmn.ErrInvalidInputList, len(assetAddrList), "greater than 0")
	}
	decimalList, ok := args[2].([]uint8)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, reflect.TypeOf(args[2]), args[2])
	}
	if assetNumber != len(decimalList) {
		return nil, fmt.Errorf(exocmn.ErrInvalidInputList, len(decimalList), assetNumber)
	}
	tvlLimitList, ok := args[3].([]*big.Int)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, reflect.TypeOf(args[3]), args[3])
	}
	if assetNumber != len(tvlLimitList) {
		return nil, fmt.Errorf(exocmn.ErrInvalidInputList, len(tvlLimitList), assetNumber)
	}

	for i := 0; i < assetNumber; i++ {
		if len(assetAddrList[i]) != types.GeneralAssetsAddrLength {
			return nil, xerrors.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddrList[i]), types.GeneralClientChainAddrLength)
		}
		assetAddr := assetAddrList[i][:clientChainAddrLength]
		assets = append(assets, types.AssetInfo{
			Address:          hexutil.Encode(assetAddr),
			Decimals:         uint32(decimalList[i]),
			LayerZeroChainID: uint64(clientChainID),
			TotalSupply:      sdkmath.NewIntFromBigInt(tvlLimitList[i]),
		})
	}
	return assets, nil
}
