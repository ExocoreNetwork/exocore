package clientchains_test

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/app"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

func (s *ClientChainsPrecompileSuite) TestIsTransaction() {
	testCases := []struct {
		name   string
		method string
		isTx   bool
	}{
		{
			"non existant method",
			"HelloFakeMethod",
			false,
		},
		{
			"actual method",
			"getClientChains",
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

func (s *ClientChainsPrecompileSuite) TestGetClientChains() {
	input, err := s.precompile.Pack("getClientChains")
	s.Require().NoError(err, "failed to pack input")
	output, err := s.precompile.Methods["getClientChains"].Outputs.Pack(true, []uint16{101})
	s.Require().NoError(err, "failed to pack output")
	testcases := []struct {
		name        string
		malleate    func() []byte
		readOnly    bool
		expPass     bool
		errContains string
		returnBytes []byte
	}{
		{
			name: "get client chains",
			malleate: func() []byte {
				return input
			},
			readOnly:    true,
			expPass:     true,
			returnBytes: output,
		},
	}
	for _, tc := range testcases {
		tc := tc
		s.Run(tc.name, func() {
			s.SetupTest()
			baseFee := s.App.FeeMarketKeeper.GetBaseFee(s.Ctx)
			contract := vm.NewPrecompile(
				vm.AccountRef(s.Address),
				s.precompile,
				big.NewInt(0),
				uint64(1e6),
			)
			contract.Input = tc.malleate()
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
			bz, err := s.precompile.Run(evm, contract, tc.readOnly)
			if tc.expPass {
				s.Require().NoError(
					err, "expected no error when running the precompile",
				)
				s.Require().Equal(
					tc.returnBytes, bz, "the return doesn't match the expected result",
				)
			} else {
				s.Require().Error(
					err, "expected error to be returned when running the precompile",
				)
				s.Require().Nil(
					bz, "expected returned bytes to be nil",
				)
				s.Require().ErrorContains(err, tc.errContains)
			}
		})
	}
}
