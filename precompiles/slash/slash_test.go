// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package slash_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
	"github.com/exocore/app"
	"github.com/exocore/precompiles/slash"
	"github.com/exocore/precompiles/withdraw"
	"github.com/exocore/utils"
	"github.com/exocore/x/restaking_assets_manage/types"
	types1 "github.com/exocore/x/withdraw/types"
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
			slash.MethodSlash,
			s.precompile.Methods[slash.MethodSlash].Name,
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

// TestRun tests the precompiles Run method withdraw.
func (s *PrecompileTestSuite) TestRunWithdraw() {
	//withdraw params for test
	exoCoreLzAppAddress := "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	exoCoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	usdtAddress := paddingClientChainAddress(common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"), types.GeneralClientChainAddrLength)
	usdcAddress := paddingClientChainAddress(common.FromHex("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), types.GeneralClientChainAddrLength)
	clientChainLzId := 101
	withdrawAddr := paddingClientChainAddress(s.address.Bytes(), types.GeneralClientChainAddrLength)
	opAmount := big.NewInt(100)
	assetAddr := usdtAddress
	commonMalleate := func() (common.Address, []byte) {
		valAddr, err := sdk.ValAddressFromBech32(s.validators[0].OperatorAddress)
		s.Require().NoError(err)
		val, _ := s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
		coins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(1e18)))
		s.app.DistrKeeper.AllocateTokensToValidator(s.ctx, val, sdk.NewDecCoinsFromCoins(coins...))
		input, err := s.precompile.Pack(
			withdraw.MethodWithdraw,
			uint16(clientChainLzId),
			assetAddr,
			withdrawAddr,
			opAmount,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.address, input
	}
	successRet, err := s.precompile.Methods[withdraw.MethodWithdraw].Outputs.Pack(true, opAmount)
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
			name: "fail - withdraw transaction will fail because the exoCoreLzAppAddress haven't been stored",
			malleate: func() (common.Address, []byte) {
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: types1.ErrNoParamsKey.Error(),
		},
		{
			name: "fail - withdraw transaction will fail because the contract caller isn't the exoCoreLzAppAddr",
			malleate: func() (common.Address, []byte) {
				withdrawModuleParam := &types1.Params{
					ExoCoreLzAppAddress:    exoCoreLzAppAddress,
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.WithdrawKeeper.SetParams(s.ctx, withdrawModuleParam)
				s.Require().NoError(err)
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: strings.Split(withdraw.ErrContractCaller, ",")[0],
		},
		{
			name: "fail - withdraw transaction will fail because the staked asset hasn't been registered",
			malleate: func() (common.Address, []byte) {
				withdrawModuleParam := &types1.Params{
					ExoCoreLzAppAddress:    s.address.String(),
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.WithdrawKeeper.SetParams(s.ctx, withdrawModuleParam)
				s.Require().NoError(err)
				assetAddr = usdcAddress
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     false,
			errContains: types1.ErrWithdrawAssetNotExist.Error(),
		},
		{
			name: "pass - withdraw transaction",
			malleate: func() (common.Address, []byte) {
				withdrawModuleParam := &types1.Params{
					ExoCoreLzAppAddress:    s.address.String(),
					ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
				}
				err := s.app.WithdrawKeeper.SetParams(s.ctx, withdrawModuleParam)
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
