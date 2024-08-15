package avs_test

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/app"
	"github.com/ExocoreNetwork/exocore/precompiles/avs"
	"github.com/ExocoreNetwork/exocore/x/avs/types"
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/cometbft/cometbft/libs/rand"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	utiltx "github.com/evmos/evmos/v14/testutil/tx"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

func (suite *AVSManagerPrecompileSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			avs.MethodRegisterAVS,
			suite.precompile.Methods[avs.MethodRegisterAVS].Name,
			true,
		},
		{
			avs.MethodDeregisterAVS,
			suite.precompile.Methods[avs.MethodDeregisterAVS].Name,
			true,
		},
		{
			avs.MethodUpdateAVS,
			suite.precompile.Methods[avs.MethodUpdateAVS].Name,
			true,
		},
		{
			avs.MethodRegisterOperatorToAVS,
			suite.precompile.Methods[avs.MethodRegisterOperatorToAVS].Name,
			true,
		},
		{
			avs.MethodDeregisterOperatorFromAVS,
			suite.precompile.Methods[avs.MethodDeregisterOperatorFromAVS].Name,
			true,
		},
		{
			avs.MethodSubmitProof,
			suite.precompile.Methods[avs.MethodSubmitProof].Name,
			true,
		},
		{
			avs.MethodCreateAVSTask,
			suite.precompile.Methods[avs.MethodCreateAVSTask].Name,
			true,
		},
		{
			avs.MethodRegisterBLSPublicKey,
			suite.precompile.Methods[avs.MethodRegisterBLSPublicKey].Name,
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.Require().Equal(suite.precompile.IsTransaction(tc.method), tc.isTx)
		})
	}
}

func (suite *AVSManagerPrecompileSuite) TestRegisterAVS() {
	avsName, operatorAddress, slashAddress, rewardAddress := "avsTest", "exo18cggcpvwspnd5c6ny8wrqxpffj5zmhklprtnph", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"
	avsOwnerAddress := []string{
		sdk.AccAddress(suite.Address.Bytes()).String(),
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()).String(),
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()).String(),
	}
	assetID := []string{"11", "22", "33"}
	minStakeAmount, taskAddr := uint64(3), "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"
	avsUnbondingPeriod, minSelfDelegation := uint64(3), uint64(3)
	epochIdentifier := epochstypes.DayEpochID
	params := []uint64{2, 3, 4, 4}

	registerOperator := func() {
		registerReq := &operatortypes.RegisterOperatorReq{
			FromAddress: operatorAddress,
			Info: &operatortypes.OperatorInfo{
				EarningsAddr: operatorAddress,
			},
		}
		_, err := suite.OperatorMsgServer.RegisterOperator(sdk.WrapSDKContext(suite.Ctx), registerReq)
		suite.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		input, err := suite.precompile.Pack(
			avs.MethodRegisterAVS,
			avsName,
			minStakeAmount,
			common.HexToAddress(taskAddr),
			common.HexToAddress(slashAddress),
			common.HexToAddress(rewardAddress),
			avsOwnerAddress,
			assetID,
			avsUnbondingPeriod,
			minSelfDelegation,
			epochIdentifier,
			params,
		)
		suite.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}

	successRet, err := suite.precompile.Methods[avs.MethodRegisterAVS].Outputs.Pack(true)
	suite.Require().NoError(err)

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
		suite.Run(tc.name, func() {

			baseFee := suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), suite.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   suite.App.EvmKeeper.ChainID(),
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

			msgEthereumTx.From = suite.Address.String()
			err := msgEthereumTx.Sign(suite.EthSigner, suite.Signer)
			suite.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := suite.Ctx.BlockHeader().ProposerAddress
			cfg, err := suite.App.EvmKeeper.EVMConfig(suite.Ctx, proposerAddress, suite.App.EvmKeeper.ChainID())
			suite.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(suite.EthSigner, baseFee)
			suite.Require().NoError(err, "failed to instantiate Ethereum message")

			// Instantiate EVM
			evm := suite.App.EvmKeeper.NewEVM(
				suite.Ctx, msg, cfg, nil, suite.StateDB,
			)

			params := suite.App.EvmKeeper.GetParams(suite.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := suite.App.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			suite.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := suite.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				suite.Require().NoError(err, "expected no error when running the precompile")
				suite.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				suite.Require().Error(err, "expected error to be returned when running the precompile")
				suite.Require().Nil(bz, "expected returned bytes to be nil")
				suite.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}

func (suite *AVSManagerPrecompileSuite) TestDeregisterAVS() {
	avsName := "avsTest"
	commonMalleate := func() (common.Address, []byte) {
		// prepare the call input for delegation test
		input, err := suite.precompile.Pack(
			avs.MethodDeregisterAVS,
			avsName,
		)
		suite.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}
	successRet, err := suite.precompile.Methods[avs.MethodDeregisterAVS].Outputs.Pack(true)
	suite.Require().NoError(err)

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
				suite.TestRegisterAVS()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     true,
			returnBytes: successRet,
		},
	}

	for _, tc := range testcases {
		tc := tc
		suite.Run(tc.name, func() {
			baseFee := suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), suite.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   suite.App.EvmKeeper.ChainID(),
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

			msgEthereumTx.From = suite.Address.String()
			err := msgEthereumTx.Sign(suite.EthSigner, suite.Signer)
			suite.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := suite.Ctx.BlockHeader().ProposerAddress
			cfg, err := suite.App.EvmKeeper.EVMConfig(suite.Ctx, proposerAddress, suite.App.EvmKeeper.ChainID())
			suite.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(suite.EthSigner, baseFee)
			suite.Require().NoError(err, "failed to instantiate Ethereum message")

			// Instantiate EVM
			evm := suite.App.EvmKeeper.NewEVM(
				suite.Ctx, msg, cfg, nil, suite.StateDB,
			)

			params := suite.App.EvmKeeper.GetParams(suite.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := suite.App.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			suite.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := suite.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				suite.Require().NoError(err, "expected no error when running the precompile")
				suite.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				suite.Require().Error(err, "expected error to be returned when running the precompile")
				suite.Require().Nil(bz, "expected returned bytes to be nil")
				suite.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}

func (suite *AVSManagerPrecompileSuite) TestUpdateAVS() {
	avsName, slashAddress, rewardAddress := "avsTest", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"
	avsOwnerAddress := []string{
		sdk.AccAddress(suite.Address.Bytes()).String(),
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()).String(),
		sdk.AccAddress(utiltx.GenerateAddress().Bytes()).String(),
	}
	assetID := []string{"11", "22", "33"}
	minStakeAmount, taskAddr := uint64(3), "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"
	avsUnbondingPeriod, minSelfDelegation := uint64(3), uint64(3)
	epochIdentifier := epochstypes.DayEpochID
	params := []uint64{2, 3, 4, 4}
	commonMalleate := func() (common.Address, []byte) {
		input, err := suite.precompile.Pack(
			avs.MethodUpdateAVS,
			avsName,
			minStakeAmount,
			common.HexToAddress(taskAddr),
			common.HexToAddress(slashAddress),
			common.HexToAddress(rewardAddress),
			avsOwnerAddress,
			assetID,
			avsUnbondingPeriod,
			minSelfDelegation,
			epochIdentifier,
			params,
		)
		suite.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}

	successRet, err := suite.precompile.Methods[avs.MethodUpdateAVS].Outputs.Pack(true)
	suite.Require().NoError(err)

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
				suite.TestRegisterAVS()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     true,
			returnBytes: successRet,
		},
	}

	for _, tc := range testcases {
		tc := tc
		suite.Run(tc.name, func() {
			baseFee := suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), suite.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   suite.App.EvmKeeper.ChainID(),
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

			msgEthereumTx.From = suite.Address.String()
			err := msgEthereumTx.Sign(suite.EthSigner, suite.Signer)
			suite.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := suite.Ctx.BlockHeader().ProposerAddress
			cfg, err := suite.App.EvmKeeper.EVMConfig(suite.Ctx, proposerAddress, suite.App.EvmKeeper.ChainID())
			suite.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(suite.EthSigner, baseFee)
			suite.Require().NoError(err, "failed to instantiate Ethereum message")

			// Instantiate EVM
			evm := suite.App.EvmKeeper.NewEVM(
				suite.Ctx, msg, cfg, nil, suite.StateDB,
			)

			params := suite.App.EvmKeeper.GetParams(suite.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := suite.App.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			suite.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := suite.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				suite.Require().NoError(err, "expected no error when running the precompile")
				suite.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				suite.Require().Error(err, "expected error to be returned when running the precompile")
				suite.Require().Nil(bz, "expected returned bytes to be nil")
				suite.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}

func (suite *AVSManagerPrecompileSuite) TestRegisterOperatorToAVS() {
	// from := s.Address
	operatorAddress := sdk.AccAddress(suite.Address.Bytes()).String()

	registerOperator := func() {
		registerReq := &operatortypes.RegisterOperatorReq{
			FromAddress: operatorAddress,
			Info: &operatortypes.OperatorInfo{
				EarningsAddr: operatorAddress,
			},
		}
		_, err := suite.OperatorMsgServer.RegisterOperator(sdk.WrapSDKContext(suite.Ctx), registerReq)
		suite.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		input, err := suite.precompile.Pack(
			avs.MethodRegisterOperatorToAVS,
		)
		suite.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}
	successRet, err := suite.precompile.Methods[avs.MethodRegisterAVS].Outputs.Pack(true)
	suite.Require().NoError(err)

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
				suite.TestRegisterAVS()
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
		suite.Run(tc.name, func() {
			baseFee := suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx)

			// malleate testcase
			caller, input := tc.malleate()
			contract := vm.NewPrecompile(vm.AccountRef(caller), suite.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input
			contract.CallerAddress = caller

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   suite.App.EvmKeeper.ChainID(),
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

			msgEthereumTx.From = suite.Address.String()
			err := msgEthereumTx.Sign(suite.EthSigner, suite.Signer)
			suite.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := suite.Ctx.BlockHeader().ProposerAddress
			cfg, err := suite.App.EvmKeeper.EVMConfig(suite.Ctx, proposerAddress, suite.App.EvmKeeper.ChainID())
			suite.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(suite.EthSigner, baseFee)
			suite.Require().NoError(err, "failed to instantiate Ethereum message")

			// Instantiate EVM
			evm := suite.App.EvmKeeper.NewEVM(
				suite.Ctx, msg, cfg, nil, suite.StateDB,
			)

			params := suite.App.EvmKeeper.GetParams(suite.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := suite.App.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			suite.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := suite.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				suite.Require().NoError(err, "expected no error when running the precompile")
				suite.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				suite.Require().Error(err, "expected error to be returned when running the precompile")
				suite.Require().Nil(bz, "expected returned bytes to be nil")
				suite.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}

func (suite *AVSManagerPrecompileSuite) TestDeregisterOperatorFromAVS() {
	// from := s.Address
	// operatorAddress, err := util.ProcessAddress(from.String())

	// registerOperator := func() {
	// 	registerReq := &operatortypes.RegisterOperatorReq{
	// 		FromAddress: operatorAddress,
	// 		Info: &operatortypes.OperatorInfo{
	// 			EarningsAddr: operatorAddress,
	// 		},
	// 	}
	// 	_, err := s.OperatorMsgServer.RegisterOperator(sdk.WrapSDKContext(s.Ctx), registerReq)
	// 	s.NoError(err)
	// }
	commonMalleate := func() (common.Address, []byte) {
		input, err := suite.precompile.Pack(
			avs.MethodDeregisterOperatorFromAVS,
		)
		suite.Require().NoError(err, "failed to pack input")
		return common.HexToAddress("0x3e108c058e8066DA635321Dc3018294cA82ddEdf"), input
	}
	successRet, err := suite.precompile.Methods[avs.MethodDeregisterOperatorFromAVS].Outputs.Pack(true)
	suite.Require().NoError(err)

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
				suite.TestRegisterOperatorToAVS()
				// registerOperator()
				return commonMalleate()
			},
			readOnly:    false,
			expPass:     true,
			returnBytes: successRet,
		},
	}

	for _, tc := range testcases {
		tc := tc
		suite.Run(tc.name, func() {
			baseFee := suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx)

			// malleate testcase
			caller, input := tc.malleate()
			contract := vm.NewPrecompile(vm.AccountRef(caller), suite.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input
			contract.CallerAddress = caller

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   suite.App.EvmKeeper.ChainID(),
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

			msgEthereumTx.From = suite.Address.String()
			err := msgEthereumTx.Sign(suite.EthSigner, suite.Signer)
			suite.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := suite.Ctx.BlockHeader().ProposerAddress
			cfg, err := suite.App.EvmKeeper.EVMConfig(suite.Ctx, proposerAddress, suite.App.EvmKeeper.ChainID())
			suite.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(suite.EthSigner, baseFee)
			suite.Require().NoError(err, "failed to instantiate Ethereum message")

			// Instantiate EVM
			evm := suite.App.EvmKeeper.NewEVM(
				suite.Ctx, msg, cfg, nil, suite.StateDB,
			)

			params := suite.App.EvmKeeper.GetParams(suite.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := suite.App.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			suite.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := suite.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				suite.Require().NoError(err, "expected no error when running the precompile")
				suite.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				suite.Require().Error(err, "expected error to be returned when running the precompile")
				suite.Require().Nil(bz, "expected returned bytes to be nil")
				suite.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}

// TestRun tests the precompiles Run method reg avstask.
func (suite *AVSManagerPrecompileSuite) TestRunRegTaskInfo() {
	taskAddr := utiltx.GenerateAddress()
	registerAVS := func() {
		avsName := "avsTest"
		avsOwnerAddress := []string{
			sdk.AccAddress(suite.Address.Bytes()).String(),
			"exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkjr",
			"exo13h6xg79g82e2g2vhjwg7j4r2z2hlncelwutkj2",
		}
		assetID := []string{"11", "22", "33"}
		avsInfo := &types.AVSInfo{
			Name:                avsName,
			AvsAddress:          utiltx.GenerateAddress().String(),
			SlashAddr:           utiltx.GenerateAddress().String(),
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
			TaskAddr:            taskAddr.String(),
		}

		err := suite.App.AVSManagerKeeper.SetAVSInfo(suite.Ctx, avsInfo)
		suite.NoError(err)
	}
	commonMalleate := func() (common.Address, []byte) {
		input, err := suite.precompile.Pack(
			avs.MethodCreateAVSTask,
			"test-avstask",
			rand.Bytes(3),
			"3",
			uint64(3),
			uint64(3),
			uint64(3),
		)
		suite.Require().NoError(err, "failed to pack input")
		return suite.Address, input
	}
	successRet, err := suite.precompile.Methods[avs.MethodCreateAVSTask].Outputs.Pack(true)
	suite.Require().NoError(err)
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
				suite.Require().NoError(err)
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
		suite.Run(tc.name, func() {
			baseFee := suite.App.FeeMarketKeeper.GetBaseFee(suite.Ctx)

			// malleate testcase
			caller, input := tc.malleate()

			contract := vm.NewPrecompile(vm.AccountRef(caller), suite.precompile, big.NewInt(0), uint64(1e6))
			contract.Input = input
			contract.CallerAddress = taskAddr

			contractAddr := contract.Address()
			// Build and sign Ethereum transaction
			txArgs := evmtypes.EvmTxArgs{
				ChainID:   suite.App.EvmKeeper.ChainID(),
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

			msgEthereumTx.From = suite.Address.String()
			err := msgEthereumTx.Sign(suite.EthSigner, suite.Signer)
			suite.Require().NoError(err, "failed to sign Ethereum message")

			// Instantiate config
			proposerAddress := suite.Ctx.BlockHeader().ProposerAddress
			cfg, err := suite.App.EvmKeeper.EVMConfig(suite.Ctx, proposerAddress, suite.App.EvmKeeper.ChainID())
			suite.Require().NoError(err, "failed to instantiate EVM config")

			msg, err := msgEthereumTx.AsMessage(suite.EthSigner, baseFee)
			suite.Require().NoError(err, "failed to instantiate Ethereum message")

			// Create StateDB
			suite.StateDB = statedb.New(suite.Ctx, suite.App.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.Ctx.HeaderHash().Bytes())))
			// Instantiate EVM
			evm := suite.App.EvmKeeper.NewEVM(
				suite.Ctx, msg, cfg, nil, suite.StateDB,
			)
			params := suite.App.EvmKeeper.GetParams(suite.Ctx)
			activePrecompiles := params.GetActivePrecompilesAddrs()
			precompileMap := suite.App.EvmKeeper.Precompiles(activePrecompiles...)
			err = vm.ValidatePrecompiles(precompileMap, activePrecompiles)
			suite.Require().NoError(err, "invalid precompiles", activePrecompiles)
			evm.WithPrecompiles(precompileMap, activePrecompiles)

			// Run precompiled contract
			bz, err := suite.precompile.Run(evm, contract, tc.readOnly)

			// Check results
			if tc.expPass {
				suite.Require().NoError(err, "expected no error when running the precompile")
				suite.Require().Equal(tc.returnBytes, bz, "the return doesn't match the expected result")
			} else {
				suite.Require().Error(err, "expected error to be returned when running the precompile")
				suite.Require().Nil(bz, "expected returned bytes to be nil")
				suite.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}
