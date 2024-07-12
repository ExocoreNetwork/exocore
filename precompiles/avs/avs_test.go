package avs_test

import (
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	"math/big"

	"github.com/ExocoreNetwork/exocore/app"
	"github.com/ExocoreNetwork/exocore/precompiles/avs"
	avskeeper "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

func (s *AVSManagerPrecompileSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			avs.MethodOperatorAction,
			s.precompile.Methods[avs.MethodOperatorAction].Name,
			true,
		},
		{
			avs.MethodAVSAction,
			s.precompile.Methods[avs.MethodAVSAction].Name,
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Require().Equal(s.precompile.IsTransaction(tc.method), tc.isTx)
		})
	}
}

func (s *AVSManagerPrecompileSuite) TestAVSManager() {
	avsName, operatorAddress, slashAddress := "avsTest", "exo18cggcpvwspnd5c6ny8wrqxpffj5zmhklprtnph", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"

	avsAction := avskeeper.RegisterAction
	from := s.Address
	avsOwnerAddress := []string{"0x3e108c058e8066DA635321Dc3018294cA82ddEdf", "0xDF907c29719154eb9872f021d21CAE6E5025d7aB", from.String()}
	assetID := []string{"11", "22", "33"}
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
		// prepare the call input for delegation test
		input, err := s.precompile.Pack(
			avs.MethodAVSAction,
			avsOwnerAddress,
			avsName,
			slashAddress,
			assetID,
			uint64(avsAction),
			uint64(10),
			uint64(7),
			epochIdentifier,
		)
		s.Require().NoError(err, "failed to pack input")
		return s.Address, input
	}
	successRet, err := s.precompile.Methods[avs.MethodAVSAction].Outputs.Pack(true)
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
			name: "pass for avs-manager",
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
			// setup basic test suite
			s.SetupTest()

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
