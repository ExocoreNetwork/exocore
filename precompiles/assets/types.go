package assets

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	sdkmath "cosmossdk.io/math"
	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	assetskeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	oracletypes "github.com/ExocoreNetwork/exocore/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

func (p Precompile) DepositWithdrawParamsFromInputs(ctx sdk.Context, args []interface{}) (*assetskeeper.DepositWithdrawParams, error) {
	inputsLen := len(p.ABI.Methods[MethodDepositTo].Inputs)
	if len(args) != inputsLen {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
	}
	depositWithdrawParams := &assetskeeper.DepositWithdrawParams{}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
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
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", args[1])
	}
	if uint32(len(assetAddr)) < clientChainAddrLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddr), clientChainAddrLength)
	}
	depositWithdrawParams.AssetsAddress = assetAddr[:clientChainAddrLength]

	stakerAddr, ok := args[2].([]byte)
	if !ok || stakerAddr == nil {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "[]byte", args[2])
	}
	if uint32(len(stakerAddr)) < clientChainAddrLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidAddrLength, len(stakerAddr), clientChainAddrLength)
	}
	depositWithdrawParams.StakerAddress = stakerAddr[:clientChainAddrLength]

	opAmount, ok := args[3].(*big.Int)
	if !ok || opAmount == nil || !(opAmount.Cmp(big.NewInt(0)) == 1) {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, "*big.Int", args[3])
	}
	depositWithdrawParams.OpAmount = sdkmath.NewIntFromBigInt(opAmount)

	return depositWithdrawParams, nil
}

func (p Precompile) ClientChainInfoFromInputs(_ sdk.Context, args []interface{}) (*types.ClientChainInfo, error) {
	inputsLen := len(p.ABI.Methods[MethodRegisterOrUpdateClientChain].Inputs)
	if len(args) != inputsLen {
		return nil, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
	}
	clientChain := types.ClientChainInfo{}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
	}
	clientChain.LayerZeroChainID = uint64(clientChainID)

	addressLength, ok := args[1].(uint8)
	if !ok || addressLength == 0 {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "uint8", args[1])
	}
	if addressLength < types.MinClientChainAddrLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidAddrLength, addressLength, types.MinClientChainAddrLength)
	}
	clientChain.AddressLength = uint32(addressLength)

	name, ok := args[2].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "string", args[2])
	}
	if name == "" || len(name) > types.MaxChainTokenNameLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidNameLength, name, len(name), types.MaxChainTokenNameLength)
	}
	clientChain.Name = name

	metaInfo, ok := args[3].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, "string", args[2])
	}
	if metaInfo == "" || len(metaInfo) > types.MaxChainTokenMetaInfoLength {
		return nil, fmt.Errorf(exocmn.ErrInvalidMetaInfoLength, metaInfo, len(metaInfo), types.MaxChainTokenMetaInfoLength)
	}
	clientChain.MetaInfo = metaInfo

	signatureType, ok := args[4].(string)
	if !ok {
		return nil, fmt.Errorf(exocmn.ErrContractInputParaOrType, 4, "string", args[4])
	}
	clientChain.SignatureType = signatureType

	return &clientChain, nil
}

func (p Precompile) TokenFromInputs(ctx sdk.Context, args []interface{}) (asset types.AssetInfo, oInfo oracletypes.OracleInfo, err error) {
	inputsLen := len(p.ABI.Methods[MethodRegisterOrUpdateTokens].Inputs)
	if len(args) != inputsLen {
		err = fmt.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
		return asset, oInfo, err
	}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		err = fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
		return asset, oInfo, err
	}
	asset.LayerZeroChainID = uint64(clientChainID)
	info, err := p.assetsKeeper.GetClientChainInfoByIndex(ctx, asset.LayerZeroChainID)
	if err != nil {
		return asset, oInfo, err
	}
	clientChainAddrLength := info.AddressLength

	assetAddr, ok := args[1].([]byte)
	if !ok || assetAddr == nil {
		err = fmt.Errorf(exocmn.ErrContractInputParaOrType, 1, "[]byte", args[1])
		return asset, oInfo, err
	}
	if uint32(len(assetAddr)) < clientChainAddrLength {
		err = fmt.Errorf(exocmn.ErrInvalidAddrLength, len(assetAddr), clientChainAddrLength)
		return asset, oInfo, err
	}
	asset.Address = hexutil.Encode(assetAddr[:clientChainAddrLength])

	decimal, ok := args[2].(uint8)
	if !ok {
		err = fmt.Errorf(exocmn.ErrContractInputParaOrType, 2, "uint8", args[2])
		return asset, oInfo, err
	}
	asset.Decimals = uint32(decimal)

	tvlLimit, ok := args[3].(*big.Int)
	if !ok || tvlLimit == nil || !(tvlLimit.Cmp(big.NewInt(0)) == 1) {
		err = fmt.Errorf(exocmn.ErrContractInputParaOrType, 3, "*big.Int", args[3])
		return asset, oInfo, err
	}
	asset.TotalSupply = sdkmath.NewIntFromBigInt(tvlLimit)

	name, ok := args[4].(string)
	if !ok {
		err = fmt.Errorf(exocmn.ErrContractInputParaOrType, 4, "string", args[4])
		return asset, oInfo, err
	}
	if name == "" || len(name) > types.MaxChainTokenNameLength {
		err = fmt.Errorf(exocmn.ErrInvalidNameLength, name, len(name), types.MaxChainTokenNameLength)
		return asset, oInfo, err
	}
	asset.Name = name

	metaInfo, ok := args[5].(string)
	if !ok {
		err = fmt.Errorf(exocmn.ErrContractInputParaOrType, 5, "string", args[5])
		return asset, oInfo, err
	}
	if metaInfo == "" || len(metaInfo) > types.MaxChainTokenMetaInfoLength {
		err = fmt.Errorf(exocmn.ErrInvalidMetaInfoLength, metaInfo, len(metaInfo), types.MaxChainTokenMetaInfoLength)
		return asset, oInfo, err
	}
	asset.MetaInfo = metaInfo

	oInfoStr, ok := args[6].(string)
	if !ok {
		err = fmt.Errorf(exocmn.ErrContractInputParaOrType, 6, "string", args[6])
		return asset, oInfo, err
	}

	if err = json.Unmarshal([]byte(oInfoStr), &oInfo); err != nil {
		return asset, oInfo, err
	}
	if len(oInfo.Token.Name) == 0 ||
		len(oInfo.Token.Chain.Name) == 0 ||
		len(oInfo.Token.Decimal) == 0 {
		err = errors.New(exocmn.ErrInvalidOracleInfo)
		return asset, oInfo, err
	}

	return asset, oInfo, err
}

func (p Precompile) ClientChainIDFromInputs(_ sdk.Context, args []interface{}) (uint32, error) {
	inputsLen := len(p.ABI.Methods[MethodIsRegisteredClientChain].Inputs)
	if len(args) != inputsLen {
		return 0, fmt.Errorf(cmn.ErrInvalidNumberOfArgs, inputsLen, len(args))
	}
	clientChainID, ok := args[0].(uint32)
	if !ok {
		return 0, fmt.Errorf(exocmn.ErrContractInputParaOrType, 0, "uint32", args[0])
	}
	return clientChainID, nil
}
