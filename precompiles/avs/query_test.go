package avs_test

import (
	"fmt"
	"math/big"

	avsManagerPrecompile "github.com/ExocoreNetwork/exocore/precompiles/avs"
	exocmn "github.com/ExocoreNetwork/exocore/precompiles/common"
	"github.com/ExocoreNetwork/exocore/x/operator/types"

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
