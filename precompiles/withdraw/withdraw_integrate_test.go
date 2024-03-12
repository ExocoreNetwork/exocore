package withdraw_test

import (
	"math/big"

	"github.com/ExocoreNetwork/exocore/precompiles/testutil"
	"github.com/ExocoreNetwork/exocore/precompiles/testutil/contracts"
	"github.com/ExocoreNetwork/exocore/precompiles/withdraw"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ethereum/go-ethereum/common"
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

func (s *WithdrawPrecompileTestSuite) TestCallWithdrawFromEOA() {
	// withdraw params for test
	exocoreLzAppAddress := "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	exocoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	params := assetstype.Params{
		ExocoreLzAppAddress:    exocoreLzAppAddress,
		ExocoreLzAppEventTopic: exocoreLzAppEventTopic,
	}
	usdtAddress := paddingClientChainAddress(common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"), assetstype.GeneralClientChainAddrLength)
	clientChainLzID := 101
	stakerAddr := paddingClientChainAddress(s.Address.Bytes(), assetstype.GeneralClientChainAddrLength)
	opAmount := big.NewInt(100)
	assetAddr := usdtAddress
	method := "withdrawPrinciple"

	beforeEach := func() {
		s.SetupTest()
		// set the default call arguments
		defaultCallArgs = contracts.CallArgs{
			ContractAddr: s.precompile.Address(),
			ContractABI:  s.precompile.ABI,
			PrivKey:      s.PrivKey,
		}

		defaultLogCheck = testutil.LogCheckArgs{
			ABIEvents: s.precompile.ABI.Events,
		}
		passCheck = defaultLogCheck.WithExpPass(true)
	}

	prepareFunc := func(params *assetstype.Params, method string) contracts.CallArgs {
		err := s.App.AssetsKeeper.SetParams(s.Ctx, params)
		s.Require().NoError(err)
		defaultWithdrawArgs := defaultCallArgs.WithMethodName(method)
		return defaultWithdrawArgs.WithArgs(
			uint16(clientChainLzID),
			assetAddr,
			stakerAddr,
			opAmount)
	}

	beforeEach()
	setWithdrawArgs := prepareFunc(&params, method)
	_, _, err := contracts.CallContractAndCheckLogs(s.Ctx, s.App, setWithdrawArgs, passCheck)
	s.Require().ErrorContains(err, withdraw.ErrContractCaller)
}
