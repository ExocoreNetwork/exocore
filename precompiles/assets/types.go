package assets

import (
	"math/big"

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
	inputsLen := len(p.ABI.Methods[MethodDepositTo].Inputs)
	if len(args) != inputsLen {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
	}
	depositWithdrawParams := &assetskeeper.DepositWithdrawParams{}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
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
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", args[1])
	}
	if uint32(len(assetAddr)) < clientChainAddrLength {
		return nil, xerrors.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddr), clientChainAddrLength)
	}
	depositWithdrawParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", args[2])
	}
	if uint32(len(stakerAddr)) < clientChainAddrLength {
		return nil, xerrors.Errorf(exocmn.ErrInvalidAddrLength, len(stakerAddr), clientChainAddrLength)
	}
	depositWithdrawParams.StakerAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || !(opAmount.Cmp(big.NewInt(0)) == 1) {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "*big.Int", args[3])
	}
	depositWithdrawParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)

	return depositWithdrawParams, nil
}

func (p Precompile) ClientChainInfoFromInputs(_ sdk.Context, args []interface{}) (*types.ClientChainInfo, error) {
	inputsLen := len(p.ABI.Methods[MethodRegisterClientChain].Inputs)
	if len(args) != inputsLen {
		return nil, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
	}
	clientChain := types.ClientChainInfo{}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
	}
	clientChain.LayerZeroChainID = uint64(clientChainID)

	addressLength, ok := args[1].(uint8)
	if !ok || addressLength == 0 {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "uint8", args[1])
	}
	if addressLength < types.MinClientChainAddrLength {
		return nil, xerrors.Errorf(exocmn.ErrInvalidAddrLength, addressLength, types.MinClientChainAddrLength)
	}
	clientChain.AddressLength = uint32(addressLength)

	name, ok := args[2].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", args[2])
	}
	if name == "" || len(name) > types.MaxChainTokenNameLength {
		return nil, xerrors.Errorf(exocmn.ErrInvalidNameLength, name, len(name), types.MaxChainTokenNameLength)
	}
	clientChain.Name = name

	metaInfo, ok := args[3].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "string", args[2])
	}
	if metaInfo == "" || len(metaInfo) > types.MaxChainTokenMetaInfoLength {
		return nil, xerrors.Errorf(exocmn.ErrInvalidMetaInfoLength, metaInfo, len(metaInfo), types.MaxChainTokenMetaInfoLength)
	}
	clientChain.MetaInfo = metaInfo

	signatureType, ok := args[4].(string)
	if !ok {
		return nil, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "string", args[4])
	}
	clientChain.SignatureType = signatureType

	return &clientChain, nil
}

func (p Precompile) TokenFromInputs(ctx sdk.Context, args []interface{}) (types.AssetInfo, error) {
	inputsLen := len(p.ABI.Methods[MethodRegisterToken].Inputs)
	if len(args) != inputsLen {
		return types.AssetInfo{}, xerrors.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
	}
	asset := types.AssetInfo{}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
	}
	asset.LayerZeroChainID = uint64(clientChainID)
	info, err := p.assetsKeeper.GetClientChainInfoByIndex(ctx, asset.LayerZeroChainID)
	if err != nil {
		return types.AssetInfo{}, err
	}
	clientChainAddrLength := info.AddressLength

	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", args[1])
	}
	if uint32(len(assetAddr)) < clientChainAddrLength {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddr), clientChainAddrLength)
	}
	asset.Address = hexutil.Encode(assetAddr[:clientChainAddrLength])

	decimal, ok := args[2].(uint8)
	if !ok {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 2, "uint8", args[2])
	}
	asset.Decimals = uint32(decimal)

	tvlLimit, ok := args[3].(*big.Int)
	if !ok || tvlLimit == nil || !(tvlLimit.Cmp(big.NewInt(0)) == 1) {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 3, "*big.Int", args[3])
	}
	asset.TotalSupply = sdkmath.NewIntFromBigInt(tvlLimit)

	name, ok := args[4].(string)
	if !ok {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 4, "string", args[4])
	}
	if name == "" || len(name) > types.MaxChainTokenNameLength {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrInvalidNameLength, name, len(name), types.MaxChainTokenNameLength)
	}
	asset.Name = name

	metaInfo, ok := args[5].(string)
	if !ok {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrContractInputParaOrType, 5, "string", args[5])
	}
	if metaInfo == "" || len(metaInfo) > types.MaxChainTokenMetaInfoLength {
		return types.AssetInfo{}, xerrors.Errorf(exocmn.ErrInvalidMetaInfoLength, metaInfo, len(metaInfo), types.MaxChainTokenMetaInfoLength)
	}
	asset.MetaInfo = metaInfo

	return asset, nil
}
