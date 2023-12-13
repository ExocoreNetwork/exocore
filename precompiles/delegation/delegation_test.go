package delegation_test

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/evmos/v14/utils"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
	"github.com/exocore/app"
	"github.com/exocore/precompiles/delegation"
	"github.com/exocore/precompiles/deposit"
	keeper2 "github.com/exocore/x/delegation/keeper"
	types2 "github.com/exocore/x/delegation/types"
	"github.com/exocore/x/deposit/keeper"
	types3 "github.com/exocore/x/deposit/types"
	"github.com/exocore/x/restaking_assets_manage/types"
	"math/big"
	"strings"
)

func (s *PrecompileTestSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			delegation.MethodDelegateToThroughClientChain,
			s.precompile.Methods[delegation.MethodDelegateToThroughClientChain].Name,
			true,
		},
		{
			delegation.MethodUndelegateFromThroughClientChain,
			s.precompile.Methods[delegation.MethodUndelegateFromThroughClientChain].Name,
			true,
		},
		{
			"invalid",
			"invalid",
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().Equal(s.precompile.IsTransaction(tc.method), tc.isTx)
		})
	}
}

func paddingClientChainAddress(input []byte, outputLength int) []byte {
	if len(input) < outputLength {
		padding := make([]byte, outputLength-len(input))
		return append(input, padding...)
	}
	return input
}

// TestRun tests DelegateToThroughClientChain method through calling Run function.
func (s *PrecompileTestSuite) TestRunDelegateToThroughClientChain() {
	//deposit params for test
	exoCoreLzAppAddress := "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	exoCoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	usdtAddress := common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	opAccAddr := "evmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3h6cprl"
	clientChainLzId := 101
	lzNonce := 0
	delegationAmount := big.NewInt(50)
	depositAmount := big.NewInt(100)
	smallDepositAmount := big.NewInt(20)
	assetAddr := paddingClientChainAddress(usdtAddress, types.GeneralClientChainAddrLength)
	depositAsset := func(staker []byte, depositAmount sdkmath.Int) {
		//deposit asset for delegation test
		params := &keeper.DepositParams{
			ClientChainLzId: 101,
			Action:          types.Deposit,
			StakerAddress:   staker,
			AssetsAddress:   usdtAddress,
			OpAmount:        depositAmount,
		}
		err := s.app.DepositKeeper.Deposit(s.ctx, params)
		s.Require().NoError(err)
	}
	registerOperator := func() {
		registerReq := &types2.RegisterOperatorReq{
			FromAddress: opAccAddr,
			Info: &types2.OperatorInfo{
				EarningsAddr: opAccAddr,
			},
		}
		_, err := s.app.DelegationKeeper.RegisterOperator(s.ctx, registerReq)
		s.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		//prepare the call input for delegation test
		valAddr, err := sdk.ValAddressFromBech32(s.validators[0].OperatorAddress)
		s.Require().NoError(err)
		val, _ := s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
		coins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(1e18)))
		s.app.DistrKeeper.AllocateTokensToValidator(s.ctx, val, sdk.NewDecCoinsFromCoins(coins...))
		input, err := s.precompile.Pack(
			delegation.MethodDelegateToThroughClientChain,
			uint16(clientChainLzId),
			uint64(lzNonce),
			assetAddr,
			paddingClientChainAddress(s.address.Bytes(), types.GeneralClientChainAddrLength),
			[]byte(opAccAddr),
			delegationAmount,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.address, input
	}
	successRet, err := s.precompile.Methods[delegation.MethodDelegateToThroughClientChain].Outputs.Pack(true)
	s.Require().NoError(err)

	testcases := []struct {
		name        string
		malleate    func() (common.Address, []byte)
		readOnly    bool
		expPass     bool
		errContains string
		returnBytes []byte
	}{
		{
			name: "fail - delegateToThroughClientChain transaction will fail because the exoCoreLzAppAddress haven't been stored",
			malleate: func() (common.Address, []byte) {
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: types3.ErrNoParamsKey.Error(),
		},
		{
			name: "fail - delegateToThroughClientChain transaction will fail because the contract caller isn't the exoCoreLzAppAddr",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &types3.Params{
					ExoCoreLzAppAddress:    exoCoreLzAppAddress,
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.DepositKeeper.SetParams(s.ctx, depositModuleParam)
				s.Require().NoError(err)
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: strings.Split(deposit.ErrContractCaller, ",")[0],
		},
		{
			name: "fail - delegateToThroughClientChain transaction will fail because the delegated operator hasn't been registered",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &types3.Params{
					ExoCoreLzAppAddress:    s.address.String(),
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.DepositKeeper.SetParams(s.ctx, depositModuleParam)
				s.Require().NoError(err)
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: types2.ErrOperatorNotExist.Error(),
		},
		{
			name: "fail - delegateToThroughClientChain transaction will fail because the delegated asset hasn't been deposited",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &types3.Params{
					ExoCoreLzAppAddress:    s.address.String(),
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.DepositKeeper.SetParams(s.ctx, depositModuleParam)
				s.Require().NoError(err)
				registerOperator()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: types.ErrNoStakerAssetKey.Error(),
		},
		{
			name: "fail - delegateToThroughClientChain transaction will fail because the delegation amount is bigger than the canWithdraw amount",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &types3.Params{
					ExoCoreLzAppAddress:    s.address.String(),
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.DepositKeeper.SetParams(s.ctx, depositModuleParam)
				s.Require().NoError(err)
				registerOperator()
				depositAsset(s.address.Bytes(), sdkmath.NewIntFromBigInt(smallDepositAmount))
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: types2.ErrDelegationAmountTooBig.Error(),
		},
		{
			name: "pass - delegateToThroughClientChain transaction",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &types3.Params{
					ExoCoreLzAppAddress:    s.address.String(),
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.DepositKeeper.SetParams(s.ctx, depositModuleParam)
				s.Require().NoError(err)
				registerOperator()
				depositAsset(s.address.Bytes(), sdkmath.NewIntFromBigInt(depositAmount))
				return commonMalleate()
			},
			returnBytes: successRet,
			readOnly:    false,
			expPass:     true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
			// setup basic test suite
			s.SetupTest()

			baseFee := s.app.FeeMarketKeeper.GetBaseFee(s.ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), s.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   s.app.EvmKeeper.ChainID(),
				Nonce:     0,
				To:        &contractAddr,
				Amount:    nil,
				GasLimit:  100000,
				GasPrice:  app.MainnetMinGasPrices.BigInt(),
				GasFeeCap: baseFee,
				GasTipCap: big.NewInt(1),
				Accesses:  &ethtypes.AccessList{},
			}
			msgEthereumTx := evmtypes.NewTx(&txArgs)

			msgEthereumTx.From = s.address.String()
			err := msgEthereumTx.Sign(s.ethSigner, s.signer)
			s.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := s.ctx.BlockHeader().ProposerAddress
			cfg, err := s.app.EvmKeeper.EVMConfig(s.ctx, proposerAddress, s.app.EvmKeeper.ChainID())
			s.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(s.ethSigner, baseFee)
			s.Require().NoError(err, "failed to instantiate Ethereum message")

			// Create StateDB
			s.stateDB = statedb.New(s.ctx, s.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.ctx.HeaderHash().Bytes())))
			// Instantiate EVM
			evm := s.app.EvmKeeper.NewEVM(
				s.ctx, msg, cfg, nil, s.stateDB,
			)

			params := s.app.EvmKeeper.GetParams(s.ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := s.app.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			s.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := s.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				s.Require().NoError(err, "expected no error when running the precompile")
				s.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				s.Require().Error(err, "expected error to be returned when running the precompile")
				s.Require().Nil(bz, "expected returned bytes to be nil")
				s.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}

// TestRun tests DelegateToThroughClientChain method through calling Run function.
func (s *PrecompileTestSuite) TestRunUnDelegateFromThroughClientChain() {
	//deposit params for test
	exoCoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	usdtAddress := common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	operatorAddr := "evmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3h6cprl"
	clientChainLzId := 101
	lzNonce := uint64(0)
	delegationAmount := big.NewInt(50)
	depositAmount := big.NewInt(100)
	assetAddr := paddingClientChainAddress(usdtAddress, types.GeneralClientChainAddrLength)
	depositAsset := func(staker []byte, depositAmount sdkmath.Int) {
		//deposit asset for delegation test
		params := &keeper.DepositParams{
			ClientChainLzId: 101,
			Action:          types.Deposit,
			StakerAddress:   staker,
			AssetsAddress:   usdtAddress,
			OpAmount:        depositAmount,
		}
		err := s.app.DepositKeeper.Deposit(s.ctx, params)
		s.Require().NoError(err)
	}

	delegateAsset := func(staker []byte, delegateAmount sdkmath.Int) {
		//deposit asset for delegation test
		delegateToParams := &keeper2.DelegationOrUndelegationParams{
			ClientChainLzId: 101,
			Action:          types.DelegateTo,
			StakerAddress:   staker,
			AssetsAddress:   usdtAddress,
			OpAmount:        delegateAmount,
			LzNonce:         lzNonce,
		}
		opAccAddr, err := sdk.AccAddressFromBech32(operatorAddr)
		s.Require().NoError(err)
		delegateToParams.OperatorAddress = opAccAddr
		err = s.app.DelegationKeeper.DelegateTo(s.ctx, delegateToParams)
		s.Require().NoError(err)
	}
	registerOperator := func() {
		registerReq := &types2.RegisterOperatorReq{
			FromAddress: operatorAddr,
			Info: &types2.OperatorInfo{
				EarningsAddr: operatorAddr,
			},
		}
		_, err := s.app.DelegationKeeper.RegisterOperator(s.ctx, registerReq)
		s.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		//prepare the call input for delegation test
		valAddr, err := sdk.ValAddressFromBech32(s.validators[0].OperatorAddress)
		s.Require().NoError(err)
		val, _ := s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
		coins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(1e18)))
		s.app.DistrKeeper.AllocateTokensToValidator(s.ctx, val, sdk.NewDecCoinsFromCoins(coins...))
		input, err := s.precompile.Pack(
			delegation.MethodUndelegateFromThroughClientChain,
			uint16(clientChainLzId),
			lzNonce+1,
			assetAddr,
			paddingClientChainAddress(s.address.Bytes(), types.GeneralClientChainAddrLength),
			[]byte(operatorAddr),
			delegationAmount,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.address, input
	}
	successRet, err := s.precompile.Methods[delegation.MethodUndelegateFromThroughClientChain].Outputs.Pack(true)
	s.Require().NoError(err)

	testcases := []struct {
		name        string
		malleate    func() (common.Address, []byte)
		readOnly    bool
		expPass     bool
		errContains string
		returnBytes []byte
	}{
		{
			name: "pass - undelegateFromThroughClientChain transaction",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &types3.Params{
					ExoCoreLzAppAddress:    s.address.String(),
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.DepositKeeper.SetParams(s.ctx, depositModuleParam)
				s.Require().NoError(err)
				registerOperator()
				depositAsset(s.address.Bytes(), sdkmath.NewIntFromBigInt(depositAmount))
				delegateAsset(s.address.Bytes(), sdkmath.NewIntFromBigInt(delegationAmount))
				return commonMalleate()
			},
			returnBytes: successRet,
			readOnly:    false,
			expPass:     true,
		},
	}
	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
			// setup basic test suite
			s.SetupTest()

			baseFee := s.app.FeeMarketKeeper.GetBaseFee(s.ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), s.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   s.app.EvmKeeper.ChainID(),
				Nonce:     0,
				To:        &contractAddr,
				Amount:    nil,
				GasLimit:  100000,
				GasPrice:  app.MainnetMinGasPrices.BigInt(),
				GasFeeCap: baseFee,
				GasTipCap: big.NewInt(1),
				Accesses:  &ethtypes.AccessList{},
			}
			msgEthereumTx := evmtypes.NewTx(&txArgs)

			msgEthereumTx.From = s.address.String()
			err := msgEthereumTx.Sign(s.ethSigner, s.signer)
			s.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := s.ctx.BlockHeader().ProposerAddress
			cfg, err := s.app.EvmKeeper.EVMConfig(s.ctx, proposerAddress, s.app.EvmKeeper.ChainID())
			s.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(s.ethSigner, baseFee)
			s.Require().NoError(err, "failed to instantiate Ethereum message")

			//set txHash for delegation module
			fmt.Println("the txHash is:", msgEthereumTx.Hash)
			s.ctx = s.ctx.WithValue(delegation.CtxKeyTxHash, common.HexToHash(msgEthereumTx.Hash))
			// Create StateDB
			s.stateDB = statedb.New(s.ctx, s.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.ctx.HeaderHash().Bytes())))
			// Instantiate EVM
			evm := s.app.EvmKeeper.NewEVM(
				s.ctx, msg, cfg, nil, s.stateDB,
			)
			params := s.app.EvmKeeper.GetParams(s.ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := s.app.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			s.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := s.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				s.Require().NoError(err, "expected no error when running the precompile")
				s.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				s.Require().Error(err, "expected error to be returned when running the precompile")
				s.Require().Nil(bz, "expected returned bytes to be nil")
				s.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}
