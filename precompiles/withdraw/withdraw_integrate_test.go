package withdraw_test

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/exocore/precompiles/testutil"
	"github.com/exocore/precompiles/testutil/contracts"
	"github.com/exocore/precompiles/withdraw"
	deposittype "github.com/exocore/x/deposit/types"
	"github.com/exocore/x/restaking_assets_manage/types"
)

// General variables used for integration tests
var (
	// defaultCallArgs  are the default arguments for calling the smart contract
	//
	// NOTE: this has to be populated in a BeforeEach block because the contractAddr would otherwise be a nil address.
	defaultCallArgs contracts.CallArgs

	// defaultLogCheck instantiates a log check arguments struct with the precompile ABI events populated.
	defaultLogCheck testutil.LogCheckArgs

	// passCheck defines the arguments to check if the precompile returns no error
	passCheck testutil.LogCheckArgs
)

func (s *PrecompileTestSuite) TestCallWithdrawFromEOA() {
	// withdraw params for test
	exoCoreLzAppAddress := "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	exoCoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	params := deposittype.Params{
		ExoCoreLzAppAddress:    exoCoreLzAppAddress,
		ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
	}
	usdtAddress := paddingClientChainAddress(common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"), types.GeneralClientChainAddrLength)
	clientChainLzId := 101
	stakerAddr := paddingClientChainAddress(s.address.Bytes(), types.GeneralClientChainAddrLength)
	opAmount := big.NewInt(100)
	assetAddr := usdtAddress
	method := "withdrawPrinciple"

	beforeEach := func() {
		s.SetupTest()
		// set the default call arguments
		defaultCallArgs = contracts.CallArgs{
			ContractAddr: s.precompile.Address(),
			ContractABI:  s.precompile.ABI,
			PrivKey:      s.privKey,
		}

		defaultLogCheck = testutil.LogCheckArgs{
			ABIEvents: s.precompile.ABI.Events,
		}
		passCheck = defaultLogCheck.WithExpPass(true)
	}

	prepareFunc := func(params *deposittype.Params, method string) contracts.CallArgs {
		err := s.app.DepositKeeper.SetParams(s.ctx, params)
		s.Require().NoError(err)
		defaultWithdrawArgs := defaultCallArgs.WithMethodName(method)
		return defaultWithdrawArgs.WithArgs(
			uint16(clientChainLzId),
			assetAddr,
			stakerAddr,
			opAmount)
	}

	beforeEach()
	setWithdrawArgs := prepareFunc(&params, method)
	_, _, err := contracts.CallContractAndCheckLogs(s.ctx, s.app, setWithdrawArgs, passCheck)
	s.Require().ErrorContains(err, strings.Split(withdraw.ErrContractCaller, ",")[0])
}
