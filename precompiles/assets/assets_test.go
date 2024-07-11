package assets_test

import (
	"math/big"

	sdkmath "cosmossdk.io/math"
	assetsprecompile "github.com/ExocoreNetwork/exocore/precompiles/assets"
	assetskeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	"github.com/evmos/evmos/v14/x/evm/statedb"

	"github.com/ExocoreNetwork/exocore/app"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

func (s *AssetsPrecompileSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			assetsprecompile.MethodDepositTo,
			s.precompile.Methods[assetsprecompile.MethodDepositTo].Name,
			true,
		},
		{
			assetsprecompile.MethodWithdraw,
			s.precompile.Methods[assetsprecompile.MethodWithdraw].Name,
			true,
		},
		{
			assetsprecompile.MethodGetClientChains,
			s.precompile.Methods[assetsprecompile.MethodGetClientChains].Name,
			false,
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

// TestRunDepositTo tests DepositOrWithdraw method through calling Run function..
func (s *AssetsPrecompileSuite) TestRunDepositTo() {
	// assetsprecompile params for test
	exocoreLzAppAddress := "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	exocoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	usdtAddress := paddingClientChainAddress(common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"), assetstype.GeneralClientChainAddrLength)
	usdcAddress := paddingClientChainAddress(common.FromHex("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), assetstype.GeneralClientChainAddrLength)
	clientChainLzID := 101
	stakerAddr := paddingClientChainAddress(s.Address.Bytes(), assetstype.GeneralClientChainAddrLength)
	opAmount := big.NewInt(100)
	assetAddr := usdtAddress
	commonMalleate := func() (common.Address, []byte) {
		input, err := s.precompile.Pack(
			assetsprecompile.MethodDepositTo,
			uint32(clientChainLzID),
			assetAddr,
			stakerAddr,
			opAmount,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.Address, input
	}
	successRet, err := s.precompile.Methods[assetsprecompile.MethodDepositTo].Outputs.Pack(true, opAmount)
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
			name: "fail - depositTo transaction will fail because the exocoreLzAppAddress is mismatched",
			malleate: func() (common.Address, []byte) {
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: assetstype.ErrNotEqualToLzAppAddr.Error(),
		},
		{
			name: "fail - depositTo transaction will fail because the contract caller isn't the exoCoreLzAppAddr",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &assetstype.Params{
					ExocoreLzAppAddress:    exocoreLzAppAddress,
					ExocoreLzAppEventTopic: exocoreLzAppEventTopic,
				}
				err := s.App.AssetsKeeper.SetParams(s.Ctx, depositModuleParam)
				s.Require().NoError(err)
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: assetstype.ErrNotEqualToLzAppAddr.Error(),
		},
		{
			name: "fail - depositTo transaction will fail because the staked assetsprecompile hasn't been registered",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &assetstype.Params{
					ExocoreLzAppAddress:    s.Address.String(),
					ExocoreLzAppEventTopic: exocoreLzAppEventTopic,
				}
				err := s.App.AssetsKeeper.SetParams(s.Ctx, depositModuleParam)
				s.Require().NoError(err)
				assetAddr = usdcAddress
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: assetstype.ErrNoClientChainAssetKey.Error(),
		},
		{
			name: "pass - depositTo transaction",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &assetstype.Params{
					ExocoreLzAppAddress:    s.Address.String(),
					ExocoreLzAppEventTopic: exocoreLzAppEventTopic,
				}
				assetAddr = usdtAddress
				err := s.App.AssetsKeeper.SetParams(s.Ctx, depositModuleParam)
				s.Require().NoError(err)
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

			baseFee := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), s.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   s.App.EvmKeeper.ChainID(),
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

			msgEthereumTx.From = s.Address.String()
			err := msgEthereumTx.Sign(s.EthSigner, s.Signer)
			s.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := s.Ctx.BlockHeader().ProposerAddress
			cfg, err := s.App.EvmKeeper.EVMConfig(s.Ctx, proposerAddress, s.App.EvmKeeper.ChainID())
			s.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(s.EthSigner, baseFee)
			s.Require().NoError(err, "failed to instantiate Ethereum message")

			// Instantiate EVM
			evm := s.App.EvmKeeper.NewEVM(
				s.Ctx, msg, cfg, nil, s.StateDB,
			)

			params := s.App.EvmKeeper.GetParams(s.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := s.App.EvmKeeper.Precompiles(activePrecompiles...)
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
				/*		s.Require().Error(err, "expected error to be returned when running the precompile")
						s.Require().Nil(bz, "expected returned bytes to be nil")
						s.Require().ErrorContains(err, tc.errContains)*/
				// for failed cases we expect it returns bool value instead of error
				// this is a workaround because the error returned by precompile can not be caught in EVM
				// see https://github.com/ExocoreNetwork/exocore/issues/70
				// TODO: we should figure out root cause and fix this issue to make precompiles work normally
				result, err := s.precompile.ABI.Unpack(assetsprecompile.MethodDepositTo, bz)
				s.Require().NoError(err)
				s.Require().Equal(len(result), 2)
				success, ok := result[0].(bool)
				s.Require().True(ok)
				s.Require().False(success)
			}
		})
	}
}

// TestRun tests the precompiled Run method withdraw.
func (s *AssetsPrecompileSuite) TestRunWithdrawPrincipal() {
	// deposit params for test
	exocoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	usdtAddress := common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7")
	clientChainLzID := 101
	withdrawAmount := big.NewInt(10)
	depositAmount := big.NewInt(100)
	assetAddr := paddingClientChainAddress(usdtAddress, assetstype.GeneralClientChainAddrLength)
	depositAsset := func(staker []byte, depositAmount sdkmath.Int) {
		// deposit asset for withdraw test
		params := &assetskeeper.DepositWithdrawParams{
			ClientChainLzID: 101,
			Action:          assetstype.Deposit,
			StakerAddress:   staker,
			AssetsAddress:   usdtAddress,
			OpAmount:        depositAmount,
		}
		err := s.App.AssetsKeeper.PerformDepositOrWithdraw(s.Ctx, params)
		s.Require().NoError(err)
	}

	commonMalleate := func() (common.Address, []byte) {
		// Prepare the call input for withdraw test
		input, err := s.precompile.Pack(
			assetsprecompile.MethodWithdraw,
			uint32(clientChainLzID),
			assetAddr,
			paddingClientChainAddress(s.Address.Bytes(), assetstype.GeneralClientChainAddrLength),
			withdrawAmount,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.Address, input
	}
	successRet, err := s.precompile.Methods[assetsprecompile.MethodWithdraw].Outputs.Pack(true, new(big.Int).Sub(depositAmount, withdrawAmount))
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
			name: "pass - withdraw via pre-compiles",
			malleate: func() (common.Address, []byte) {
				depositModuleParam := &assetstype.Params{
					ExocoreLzAppAddress:    s.Address.String(),
					ExocoreLzAppEventTopic: exocoreLzAppEventTopic,
				}
				err := s.App.AssetsKeeper.SetParams(s.Ctx, depositModuleParam)
				s.Require().NoError(err)
				depositAsset(s.Address.Bytes(), sdkmath.NewIntFromBigInt(depositAmount))
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

			baseFee := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), s.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   s.App.EvmKeeper.ChainID(),
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

			msgEthereumTx.From = s.Address.String()
			err := msgEthereumTx.Sign(s.EthSigner, s.Signer)
			s.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := s.Ctx.BlockHeader().ProposerAddress
			cfg, err := s.App.EvmKeeper.EVMConfig(s.Ctx, proposerAddress, s.App.EvmKeeper.ChainID())
			s.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(s.EthSigner, baseFee)
			s.Require().NoError(err, "failed to instantiate Ethereum message")

			// Create StateDB
			s.StateDB = statedb.New(s.Ctx, s.App.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(s.Ctx.HeaderHash().Bytes())))
			// Instantiate EVM
			evm := s.App.EvmKeeper.NewEVM(
				s.Ctx, msg, cfg, nil, s.StateDB,
			)
			params := s.App.EvmKeeper.GetParams(s.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := s.App.EvmKeeper.Precompiles(activePrecompiles...)
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

func (s *AssetsPrecompileSuite) TestGetClientChains() {
	input, err := s.precompile.Pack("getClientChains")
	s.Require().NoError(err, "failed to pack input")
	output, err := s.precompile.Methods["getClientChains"].Outputs.Pack(true, []uint32{101})
	s.Require().NoError(err, "failed to pack output")
	s.Run("get client chains", func() {
		s.SetupTest()
		baseFee := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx)
		contract := vm.NewPrecompile(
			vm.AccountRef(s.Address),
			s.precompile,
			big.NewInt(0),
			uint64(1e6),
		)
		contract.Input = input
		contractAddr := contract.Address()
		txArgs := evmtypes.EvmTxArgs{
			ChainID:   s.App.EvmKeeper.ChainID(),
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
		msgEthereumTx.From = s.Address.String()
		err := msgEthereumTx.Sign(s.EthSigner, s.Signer)
		s.Require().NoError(err, "failed to sign Ethereum message")
		proposerAddress := s.Ctx.BlockHeader().ProposerAddress
		cfg, err := s.App.EvmKeeper.EVMConfig(
			s.Ctx, proposerAddress, s.App.EvmKeeper.ChainID(),
		)
		s.Require().NoError(err, "failed to instantiate EVM config")
		msg, err := msgEthereumTx.AsMessage(s.EthSigner, baseFee)
		s.Require().NoError(err, "failed to instantiate Ethereum message")
		evm := s.App.EvmKeeper.NewEVM(
			s.Ctx, msg, cfg, nil, s.StateDB,
		)
		params := s.App.EvmKeeper.GetParams(s.Ctx)
		activePrecompiles := params.GetActivePrecompilesAddrs()
		precompileMap := s.App.EvmKeeper.Precompiles(activePrecompiles...)
		err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
		s.Require().NoError(err, "invalid precompiles", activePrecompiles)
		evm.WithPrecompiles(precompileMap, activePrecompiles)
		bz, err := s.precompile.Run(evm, contract, true)
		s.Require().NoError(
			err, "expected no error when running the precompile",
		)
		s.Require().Equal(
			output, bz, "the return doesn't match the expected result",
		)
	})
}
