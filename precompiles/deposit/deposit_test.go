package deposit_test

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/app"
	"github.com/ExocoreNetwork/exocore/precompiles/deposit"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	deposittype "github.com/ExocoreNetwork/exocore/x/deposit/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

func (s *DepositPrecompileSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			deposit.MethodDepositTo,
			s.precompile.Methods[deposit.MethodDepositTo].Name,
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

// TestRunDepositTo tests DepositTo method through calling Run function..
func (s *DepositPrecompileSuite) TestRunDepositTo() {
	// deposit params for test
	exocoreLzAppAddress := "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	exocoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	usdtAddress := paddingClientChainAddress(common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"), assetstype.GeneralClientChainAddrLength)
	usdcAddress := paddingClientChainAddress(common.FromHex("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), assetstype.GeneralClientChainAddrLength)
	clientChainLzID := 101
	stakerAddr := paddingClientChainAddress(s.Address.Bytes(), assetstype.GeneralClientChainAddrLength)
	opAmount := big.NewInt(100)
	assetAddr := usdtAddress
	commonMalleate := func() (common.Address, []byte) {
		// valAddr, err := sdk.ValAddressFromBech32(s.Validators[0].OperatorAddress)
		// s.Require().NoError(err)
		// val, _ := s.App.StakingKeeper.GetValidator(s.Ctx, valAddr)
		// coins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(1e18)))
		// s.App.DistrKeeper.AllocateTokensToValidator(s.Ctx, val, sdk.NewDecCoinsFromCoins(coins...))
		input, err := s.precompile.Pack(
			deposit.MethodDepositTo,
			uint16(clientChainLzID),
			assetAddr,
			stakerAddr,
			opAmount,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.Address, input
	}
	successRet, err := s.precompile.Methods[deposit.MethodDepositTo].Outputs.Pack(true, opAmount)
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
			name: "fail - depositTo transaction will fail because the exocoreLzAppAddress haven't been stored",
			malleate: func() (common.Address, []byte) {
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: assetstype.ErrNoParamsKey.Error(),
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
			name: "fail - depositTo transaction will fail because the staked asset hasn't been registered",
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
			errContains: deposittype.ErrDepositAssetNotExist.Error(),
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
				s.Require().Error(err, "expected error to be returned when running the precompile")
				s.Require().Nil(bz, "expected returned bytes to be nil")
				s.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}
