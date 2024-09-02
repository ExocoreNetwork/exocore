package assets

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"

	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
)

const (
	MethodDepositTo                   = "depositTo"
	MethodWithdraw                    = "withdrawPrincipal"
	MethodGetClientChains             = "getClientChains"
	MethodRegisterOrUpdateClientChain = "registerOrUpdateClientChain"
	MethodRegisterToken               = "registerToken"
	MethodUpdateToken                 = "updateToken"
	MethodIsRegisteredClientChain     = "isRegisteredClientChain"
)

// DepositOrWithdraw deposit and withdraw the client chain assets for the staker,
// that will change the state in assets module.
func (p Precompile) DepositOrWithdraw(
	ctx sdk.Context,
	_ common.Address,
	contract *vm.Contract,
	_ vm.StateDB,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract,the caller must be exoCore LzApp contract
	err := p.assetsKeeper.CheckExocoreGatewayAddr(ctx, contract.CallerAddress)
	if err != nil {
		return nil, fmt.Errorf(exocmn.ErrContractCaller, err.Error())
	}
	// parse the depositTo input params
	depositWithdrawParams, err := p.DepositWithdrawParamsFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	// call assets keeper to perform the deposit or withdraw action
	switch method.Name {
	// deposit transactions
	case MethodDepositTo:
		depositWithdrawParams.Action = assetstypes.Deposit
	case MethodWithdraw:
		depositWithdrawParams.Action = assetstypes.WithdrawPrincipal
	default:
		return nil, fmt.Errorf(cmn.ErrUnknownMethod, method.Name)
	}
	err = p.assetsKeeper.PerformDepositOrWithdraw(ctx, depositWithdrawParams)
	if err != nil {
		return nil, err
	}

	// get the latest asset state of staker to return.
	stakerID, assetID := assetstypes.GetStakeIDAndAssetID(depositWithdrawParams.ClientChainLzID, depositWithdrawParams.StakerAddress, depositWithdrawParams.AssetsAddress)
	info, err := p.assetsKeeper.GetStakerSpecifiedAssetInfo(ctx, stakerID, assetID)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true, info.TotalDepositAmount.BigInt())
}

func (p Precompile) GetClientChains(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	if len(args) > 0 {
		ctx.Logger().Error(
			"GetClientChains",
			"err", errors.New("no input is required"),
		)
		return method.Outputs.Pack(false, nil)
	}
	ids, err := p.assetsKeeper.GetAllClientChainID(ctx)
	if err != nil {
		ctx.Logger().Error(
			"GetClientChains",
			"err", err,
		)
		return method.Outputs.Pack(false, nil)
	}
	return method.Outputs.Pack(true, ids)
}

func (p Precompile) RegisterOrUpdateClientChain(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// check the invalidation of caller contract,the caller must be exoCore LzApp contract
	err := p.assetsKeeper.CheckExocoreGatewayAddr(ctx, contract.CallerAddress)
	if err != nil {
		return nil, fmt.Errorf(exocmn.ErrContractCaller, err.Error())
	}

	clientChainInfo, err := p.ClientChainInfoFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	updated := p.assetsKeeper.ClientChainExists(ctx, clientChainInfo.LayerZeroChainID)
	err = p.assetsKeeper.SetClientChainInfo(ctx, clientChainInfo)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Pack(true, updated)
}

func (p Precompile) RegisterToken(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// the caller must be the ExocoreGateway contract
	err := p.assetsKeeper.CheckExocoreGatewayAddr(ctx, contract.CallerAddress)
	if err != nil {
		return nil, fmt.Errorf(exocmn.ErrContractCaller, err.Error())
	}

	// parse inputs
	asset, oInfo, err := p.TokenFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	_, assetID := assetstypes.GetStakeIDAndAssetIDFromStr(asset.LayerZeroChainID, "", asset.Address)
	oInfo.AssetID = assetID

	if p.assetsKeeper.IsStakingAsset(ctx, assetID) {
		return nil, fmt.Errorf("asset %s already exists", assetID)
	}

	stakingAsset := &assetstypes.StakingAssetInfo{
		AssetBasicInfo:     &asset,
		StakingTotalAmount: sdkmath.NewInt(0),
	}

	if err := p.assetsKeeper.RegisterNewTokenAndSetTokenFeeder(ctx, &oInfo); err != nil {
		return nil, err
	}

	// this is where the magic happens
	if err := p.assetsKeeper.SetStakingAssetInfo(ctx, stakingAsset); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (p Precompile) UpdateToken(
	ctx sdk.Context,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	// the caller must be the ExocoreGateway contract
	err := p.assetsKeeper.CheckExocoreGatewayAddr(ctx, contract.CallerAddress)
	if err != nil {
		return nil, fmt.Errorf(exocmn.ErrContractCaller, err.Error())
	}

	// parse inputs
	clientChainID, hexAssetAddr, tvlLimit, metadata, err := p.UpdateTokenFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}

	// check that the asset being updated actually exists
	_, assetID := assetstypes.GetStakeIDAndAssetIDFromStr(uint64(clientChainID), "", hexAssetAddr)
	assetInfo, err := p.assetsKeeper.GetStakingAssetInfo(ctx, assetID)
	if err != nil {
		// fails if asset does not exist with ErrNoClientChainAssetKey
		return nil, err
	}

	// finally, execute the update
	if len(metadata) > 0 {
		// if metadata is not empty, update it
		assetInfo.AssetBasicInfo.MetaInfo = metadata
	}
	// always update TVL, since a value of 0 is used to reject deposits
	assetInfo.AssetBasicInfo.TotalSupply = tvlLimit

	if err := p.assetsKeeper.SetStakingAssetInfo(ctx, assetInfo); err != nil {
		return nil, err
	}

	return method.Outputs.Pack(true)
}

func (p Precompile) IsRegisteredClientChain(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) ([]byte, error) {
	clientChainID, err := p.ClientChainIDFromInputs(ctx, args)
	if err != nil {
		return nil, err
	}
	exists := p.assetsKeeper.ClientChainExists(ctx, uint64(clientChainID))
	return method.Outputs.Pack(true, exists)
}
