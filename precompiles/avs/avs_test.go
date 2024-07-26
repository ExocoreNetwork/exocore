package avs_test

import (
	"github.com/ExocoreNetwork/exocore/app"
	"github.com/ExocoreNetwork/exocore/precompiles/avs"
	util "github.com/ExocoreNetwork/exocore/utils"
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
	"math/big"
)

func (s *AVSManagerPrecompileSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			avs.MethodRegisterAVS,
			s.precompile.Methods[avs.MethodRegisterAVS].Name,
			true,
		},
		{
			avs.MethodDeregisterAVS,
			s.precompile.Methods[avs.MethodDeregisterAVS].Name,
			true,
		},
		{
			avs.MethodUpdateAVS,
			s.precompile.Methods[avs.MethodUpdateAVS].Name,
			true,
		},
		{
			avs.MethodRegisterOperatorToAVS,
			s.precompile.Methods[avs.MethodRegisterOperatorToAVS].Name,
			true,
		},
		{
			avs.MethodDeregisterOperatorFromAVS,
			s.precompile.Methods[avs.MethodDeregisterOperatorFromAVS].Name,
			true,
		},
		{
			avs.MethodSubmitProof,
			s.precompile.Methods[avs.MethodSubmitProof].Name,
			true,
		},
		{
			avs.MethodCreateAVSTask,
			s.precompile.Methods[avs.MethodCreateAVSTask].Name,
			true,
		},
		{
			avs.MethodRegisterBLSPublicKey,
			s.precompile.Methods[avs.MethodRegisterBLSPublicKey].Name,
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().Equal(s.precompile.IsTransaction(tc.method), tc.isTx)
		})
	}
}

func (s *AVSManagerPrecompileSuite) TestRegisterAVS() {
	avsName, operatorAddress, slashAddress, rewardAddress := "avsTest", "exo18cggcpvwspnd5c6ny8wrqxpffj5zmhklprtnph", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"
	avsOwnerAddress := []string{"0x3e108c058e8066DA635321Dc3018294cA82ddEdf", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", s.Address.String()}
	assetID := []string{"11", "22", "33"}
	minStakeAmount, taskAddr, miniptInOperators, minTotalStakeAmount, avsReward, avsSlash := uint64(3), "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", uint64(3), uint64(3), uint64(3), uint64(3)
	avsUnbondingPeriod, minSelfDelegation := uint64(3), uint64(3)
	epochIdentifier := epochstypes.DayEpochID
	registerOperator := func() {
		registerReq := &operatortypes.RegisterOperatorReq{
			FromAddress: operatorAddress,
			Info: &operatortypes.OperatorInfo{
				EarningsAddr: operatorAddress,
			},
		}
		_, err := s.App.OperatorKeeper.RegisterOperator(s.Ctx, registerReq)
		s.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		input, err := s.precompile.Pack(
			avs.MethodRegisterAVS,
			avsName,
			minStakeAmount,
			taskAddr,
			slashAddress,
			rewardAddress,
			avsOwnerAddress,
			assetID,
			avsUnbondingPeriod,
			minSelfDelegation,
			epochIdentifier,
			miniptInOperators,
			minTotalStakeAmount,
			avsReward,
			avsSlash,
		)
		s.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}

	successRet, err := s.precompile.Methods[avs.MethodRegisterAVS].Outputs.Pack(true)
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
			name: "pass for avs-registered",
			malleate: func() (common.Address, []byte) {
				registerOperator()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     true,
			returnBytes: successRet,
		},
	}

	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {

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

func (s *AVSManagerPrecompileSuite) TestDeregisterAVS() {
	avsName := "avsTest"
	commonMalleate := func() (common.Address, []byte) {
		// prepare the call input for delegation test
		input, err := s.precompile.Pack(
			avs.MethodDeregisterAVS,
			avsName,
		)
		s.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}
	successRet, err := s.precompile.Methods[avs.MethodDeregisterAVS].Outputs.Pack(true)
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
			name: "pass for avs-deregister",
			malleate: func() (common.Address, []byte) {
				s.TestRegisterAVS()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     true,
			returnBytes: successRet,
		},
	}

	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
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

func (s *AVSManagerPrecompileSuite) TestUpdateAVS() {
	avsName, slashAddress, rewardAddress := "avsTest", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"
	avsOwnerAddress := []string{"0x3e108c058e8066DA635321Dc3018294cA82ddEdf", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", s.Address.String()}
	assetID := []string{"11", "22", "33"}
	minStakeAmount, taskAddr, minOptInOperators, minTotalStakeAmount, avsReward, avsSlash := uint64(3), "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", uint64(3), uint64(3), uint64(3), uint64(3)
	avsUnbondingPeriod, minSelfDelegation := uint64(3), uint64(3)
	epochIdentifier := epochstypes.DayEpochID
	commonMalleate := func() (common.Address, []byte) {
		input, err := s.precompile.Pack(
			avs.MethodUpdateAVS,
			avsName,
			minStakeAmount,
			taskAddr,
			slashAddress,
			rewardAddress,
			avsOwnerAddress,
			assetID,
			avsUnbondingPeriod,
			minSelfDelegation,
			epochIdentifier,
			minOptInOperators,
			minTotalStakeAmount,
			avsReward,
			avsSlash,
		)
		s.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}

	successRet, err := s.precompile.Methods[avs.MethodUpdateAVS].Outputs.Pack(true)
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
			name: "pass for avs-update",
			malleate: func() (common.Address, []byte) {
				s.TestRegisterAVS()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     true,
			returnBytes: successRet,
		},
	}

	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
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

func (s *AVSManagerPrecompileSuite) TestRegisterOperatorToAVS() {
	from := s.Address
	operatorAddress, err := util.ProcessAddress(s.Address.String())

	registerOperator := func() {
		registerReq := &operatortypes.RegisterOperatorReq{
			FromAddress: operatorAddress,
			Info: &operatortypes.OperatorInfo{
				EarningsAddr: operatorAddress,
			},
		}
		_, err := s.App.OperatorKeeper.RegisterOperator(s.Ctx, registerReq)
		s.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		input, err := s.precompile.Pack(
			avs.MethodRegisterOperatorToAVS,
		)
		s.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}
	successRet, err := s.precompile.Methods[avs.MethodRegisterAVS].Outputs.Pack(true)
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
			name: "pass for operator opt-in avs",
			malleate: func() (common.Address, []byte) {
				s.TestRegisterAVS()
				registerOperator()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     true,
			returnBytes: successRet,
		},
	}

	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
			baseFee := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx)

			// malleate testcase
			caller, input := tc.malleate()
			contract := vm.NewPrecompile(vm.AccountRef(caller), s.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input
			contract.CallerAddress = from

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

func (s *AVSManagerPrecompileSuite) TestDeregisterOperatorFromAVS() {
	from := s.Address
	operatorAddress, err := util.ProcessAddress(from.String())

	registerOperator := func() {
		registerReq := &operatortypes.RegisterOperatorReq{
			FromAddress: operatorAddress,
			Info: &operatortypes.OperatorInfo{
				EarningsAddr: operatorAddress,
			},
		}
		_, err := s.App.OperatorKeeper.RegisterOperator(s.Ctx, registerReq)
		s.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		input, err := s.precompile.Pack(
			avs.MethodDeregisterOperatorFromAVS,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.Address, input
	}
	successRet, err := s.precompile.Methods[avs.MethodDeregisterOperatorFromAVS].Outputs.Pack(true)
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
			name: "pass for operator opt-out avs",
			malleate: func() (common.Address, []byte) {
				//s.TestRegisterOperatorToAVS()
				registerOperator()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     true,
			returnBytes: successRet,
		},
	}

	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
			baseFee := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx)

			// malleate testcase
			caller, input := tc.malleate()
			contract := vm.NewPrecompile(vm.AccountRef(caller), s.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input
			contract.CallerAddress = from

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

// TestRun tests the precompiles Run method reg avstask.
func (s *AVSManagerPrecompileSuite) TestRunRegTaskinfo() {
	registerAVS := func() {
		avsName, avsAddres, slashAddress := "avsTest", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutash"
		addr, _ := util.ProcessAddress(s.Address.String())
		avsOwnerAddress := []string{"exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr", addr, "exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj2"}
		assetID := []string{"11", "22", "33"}
		avs := &types.AVSInfo{
			Name:                avsName,
			AvsAddress:          avsAddres,
			SlashAddr:           slashAddress,
			AvsOwnerAddress:     avsOwnerAddress,
			AssetIDs:            assetID,
			AvsUnbondingPeriod:  7,
			MinSelfDelegation:   10,
			EpochIdentifier:     epochstypes.DayEpochID,
			StartingEpoch:       1,
			MinOptInOperators:   100,
			MinTotalStakeAmount: 1000,
			AvsSlash:            sdk.MustNewDecFromStr("0.001"),
			AvsReward:           sdk.MustNewDecFromStr("0.002"),
			TaskAddr:            addr,
		}

		err := s.App.AVSManagerKeeper.SetAVSInfo(s.Ctx, avs)
		s.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		input, err := s.precompile.Pack(
			avs.MethodCreateAVSTask,
			"test-avstask",
			rand.Bytes(3),
			"3",
			uint64(3),
			uint64(3),
			uint64(3),
		)
		s.Require().NoError(err, "failed to pack input")
		return s.Address, input
	}
	successRet, err := s.precompile.Methods[avs.MethodCreateAVSTask].Outputs.Pack(true)
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
			name: "pass - avstask via pre-compiles",
			malleate: func() (common.Address, []byte) {
				s.Require().NoError(err)
				registerAVS()
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
