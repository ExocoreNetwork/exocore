package avs_test

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	avsManagerPrecompile "github.com/ExocoreNetwork/exocore/precompiles/avs"
	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	operatorKeeper "github.com/ExocoreNetwork/exocore/x/operator/keeper"
	"github.com/ExocoreNetwork/exocore/x/operator/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/vm"
)

type avsTestCases struct {
	name        string
	malleate    func() []interface{}
	postCheck   func(bz []byte)
	gas         uint64
	expErr      bool
	errContains string
}

var baseTestCases = []avsTestCases{
	{
		"fail - empty input args",
		func() []interface{} {
			return []interface{}{}
		},
		func(bz []byte) {},
		100000,
		true,
		"invalid number of arguments",
	},
	{
		"fail - invalid  address",
		func() []interface{} {
			return []interface{}{
				"invalid",
			}
		},
		func(bz []byte) {},
		100000,
		true,
		"invalid bech32 string",
	},
}

func (s *AVSManagerPrecompileSuite) TestGetOptedInOperatorAccAddrs() {
	method := s.precompile.Methods[avsManagerPrecompile.MethodGetOptinOperators]
	operatorAddress, avsAddr, slashContract := "exo18cggcpvwspnd5c6ny8wrqxpffj5zmhklprtnph", s.Address, "0xDF907c29719154eb9872f021d21CAE6E5025d7aB"

	operatorOptIn := func() {
		optedInfo := &types.OptedInfo{
			SlashContract: slashContract,
			// #nosec G701
			OptedInHeight:  uint64(s.Ctx.BlockHeight()),
			OptedOutHeight: types.DefaultOptedOutHeight,
		}
		err := s.App.OperatorKeeper.SetOptedInfo(s.Ctx, operatorAddress, avsAddr.String(), optedInfo)
		s.NoError(err)
	}
	testCases := []avsTestCases{
		{
			"fail - invalid avs address",
			func() []interface{} {
				return []interface{}{
					"invalid",
				}
			},
			func(bz []byte) {},
			100000,
			true,
			fmt.Sprintf(exocmn.ErrContractInputParaOrType, 0, "string", "0x0000000000000000000000000000000000000000"),
		},
		{
			"success - no operators",
			func() []interface{} {
				return []interface{}{
					s.Address,
				}
			},
			func(bz []byte) {
				var out []string
				err := s.precompile.UnpackIntoInterface(&out, avsManagerPrecompile.MethodGetOptinOperators, bz)
				s.Require().NoError(err, "failed to unpack output", err)
				s.Require().Equal(0, len(out))
			},
			100000,
			false,
			"",
		},
		{
			"success - existent operators",
			func() []interface{} {
				operatorOptIn()
				return []interface{}{
					s.Address,
				}
			},
			func(bz []byte) {
				var out []string
				err := s.precompile.UnpackIntoInterface(&out, avsManagerPrecompile.MethodGetOptinOperators, bz)
				s.Require().NoError(err, "failed to unpack output", err)
				s.Require().Equal(1, len(out))
				s.Require().Equal(operatorAddress, out[0])

			},
			100000,
			false,
			"",
		},
	}
	testCases = append(testCases, baseTestCases[0])

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			contract := vm.NewContract(vm.AccountRef(s.Address), s.precompile, big.NewInt(0), tc.gas)

			bz, err := s.precompile.GetOptedInOperatorAccAddrs(s.Ctx, contract, &method, tc.malleate())

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			} else {
				s.Require().NoError(err)
				s.Require().NotEmpty(bz)
				tc.postCheck(bz)
			}
		})
	}
}

func (suite *AVSManagerPrecompileSuite) TestAVSUSDValue() {
	method := s.precompile.Methods[avsManagerPrecompile.MethodGetAVSUSDValue]
	expectedUSDvalue := sdkmath.LegacyNewDec(0)

	setUp := func() {
		suite.prepare()
		// register the new token
		usdcAddr := common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
		usdcClientChainAsset := assetstype.AssetInfo{
			Name:             "USD coin",
			Symbol:           "USDC",
			Address:          usdcAddr.String(),
			Decimals:         6,
			TotalSupply:      sdkmath.NewInt(1e18),
			LayerZeroChainID: 101,
			MetaInfo:         "USDC",
		}
		err := suite.App.AssetsKeeper.SetStakingAssetInfo(
			suite.Ctx,
			&assetstype.StakingAssetInfo{
				AssetBasicInfo:     &usdcClientChainAsset,
				StakingTotalAmount: sdkmath.NewInt(0),
			},
		)
		suite.NoError(err)
		// register the new AVS
		suite.prepareAvs([]string{"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48_0x65", "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"})
		// opt in
		err = suite.App.OperatorKeeper.OptIn(suite.Ctx, suite.operatorAddr, suite.avsAddr)
		suite.NoError(err)
		usdtPrice, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
		suite.NoError(err)
		usdtValue := operatorKeeper.CalculateUSDValue(suite.delegationAmount, usdtPrice.Value, suite.assetDecimal, usdtPrice.Decimal)
		// deposit and delegate another asset to the operator
		suite.NoError(err)
		suite.prepareDeposit(usdcAddr, sdkmath.NewInt(1e8))
		usdcPrice, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
		suite.NoError(err)
		delegatedAmount := sdkmath.NewIntWithDecimal(8, 7)
		suite.prepareDelegation(true, usdcAddr, delegatedAmount)

		// updating the new voting power
		usdcValue := operatorKeeper.CalculateUSDValue(suite.delegationAmount, usdcPrice.Value, suite.assetDecimal, usdcPrice.Decimal)
		expectedUSDvalue = usdcValue.Add(usdtValue)
		suite.CommitAfter(time.Hour*1 + time.Nanosecond)
		suite.CommitAfter(time.Hour*1 + time.Nanosecond)
		suite.CommitAfter(time.Hour*1 + time.Nanosecond)
	}

	testCases := []avsTestCases{
		{
			"success - existent operators",
			func() []interface{} {
				setUp()
				return []interface{}{
					suite.avsAddr,
				}
			},
			func(bz []byte) {
				var out string
				err := s.precompile.UnpackIntoInterface(&out, avsManagerPrecompile.MethodGetAVSUSDValue, bz)
				s.Require().NoError(err, "failed to unpack output", err)
				s.Require().Equal(expectedUSDvalue.String(), out)

			},
			100000,
			false,
			"",
		},
	}
	testCases = append(testCases, baseTestCases[0])

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			contract := vm.NewContract(vm.AccountRef(s.Address), s.precompile, big.NewInt(0), tc.gas)

			bz, err := s.precompile.GetAVSUSDValue(s.Ctx, contract, &method, tc.malleate())

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			} else {
				s.Require().NoError(err)
				s.Require().NotEmpty(bz)
				tc.postCheck(bz)
			}
		})
	}
}

func (suite *AVSManagerPrecompileSuite) TestGetOperatorOptedUSDValue() {
	method := s.precompile.Methods[avsManagerPrecompile.MethodGetOperatorOptedUSDValue]
	expectedUSDvalue := sdkmath.LegacyNewDec(0)

	setUp := func() {
		suite.prepare()
		// register the new token
		usdcAddr := common.HexToAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
		usdcClientChainAsset := assetstype.AssetInfo{
			Name:             "USD coin",
			Symbol:           "USDC",
			Address:          usdcAddr.String(),
			Decimals:         6,
			TotalSupply:      sdkmath.NewInt(1e18),
			LayerZeroChainID: 101,
			MetaInfo:         "USDC",
		}
		err := suite.App.AssetsKeeper.SetStakingAssetInfo(
			suite.Ctx,
			&assetstype.StakingAssetInfo{
				AssetBasicInfo:     &usdcClientChainAsset,
				StakingTotalAmount: sdkmath.NewInt(0),
			},
		)
		suite.NoError(err)
		// register the new AVS
		suite.prepareAvs([]string{"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48_0x65", "0xdac17f958d2ee523a2206206994597c13d831ec7_0x65"})
		// opt in
		err = suite.App.OperatorKeeper.OptIn(suite.Ctx, suite.operatorAddr, suite.avsAddr)
		suite.NoError(err)
		usdtPrice, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
		suite.NoError(err)
		usdtValue := operatorKeeper.CalculateUSDValue(suite.delegationAmount, usdtPrice.Value, suite.assetDecimal, usdtPrice.Decimal)
		// deposit and delegate another asset to the operator
		suite.NoError(err)
		suite.prepareDeposit(usdcAddr, sdkmath.NewInt(1e8))
		usdcPrice, err := suite.App.OperatorKeeper.OracleInterface().GetSpecifiedAssetsPrice(suite.Ctx, suite.assetID)
		suite.NoError(err)
		delegatedAmount := sdkmath.NewIntWithDecimal(8, 7)
		suite.prepareDelegation(true, usdcAddr, delegatedAmount)

		// updating the new voting power
		usdcValue := operatorKeeper.CalculateUSDValue(suite.delegationAmount, usdcPrice.Value, suite.assetDecimal, usdcPrice.Decimal)
		expectedUSDvalue = usdcValue.Add(usdtValue)
		suite.CommitAfter(time.Hour*1 + time.Nanosecond)
		suite.CommitAfter(time.Hour*1 + time.Nanosecond)
		suite.CommitAfter(time.Hour*1 + time.Nanosecond)
	}

	testCases := []avsTestCases{
		{
			"success - existent operators",
			func() []interface{} {
				setUp()
				return []interface{}{
					suite.avsAddr,
					suite.operatorAddr.String(),
				}
			},
			func(bz []byte) {
				var out string
				err := s.precompile.UnpackIntoInterface(&out, avsManagerPrecompile.MethodGetOperatorOptedUSDValue, bz)
				s.Require().NoError(err, "failed to unpack output", err)
				s.Require().Equal(expectedUSDvalue.String(), out)

			},
			100000,
			false,
			"",
		},
	}
	testCases = append(testCases, baseTestCases[0])

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			contract := vm.NewContract(vm.AccountRef(s.Address), s.precompile, big.NewInt(0), tc.gas)

			bz, err := s.precompile.GetOperatorOptedUSDValue(s.Ctx, contract, &method, tc.malleate())

			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errContains)
			} else {
				s.Require().NoError(err)
				s.Require().NotEmpty(bz)
				tc.postCheck(bz)
			}
		})
	}
}
