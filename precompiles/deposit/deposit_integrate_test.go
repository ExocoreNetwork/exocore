package deposit_test

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/exocore/precompiles/deposit"
	"github.com/exocore/precompiles/testutil"
	"github.com/exocore/precompiles/testutil/contracts"
	types3 "github.com/exocore/x/deposit/types"
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

func (s *PrecompileTestSuite) TestCallDepositToFromEOA() {
	// deposit params for test
	exoCoreLzAppAddress := "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	exoCoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	depositParams := types3.Params{
		ExoCoreLzAppAddress:    exoCoreLzAppAddress,
		ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
	}
	usdtAddress := paddingClientChainAddress(common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"), types.GeneralClientChainAddrLength)
	clientChainLzId := 101
	stakerAddr := paddingClientChainAddress(s.address.Bytes(), types.GeneralClientChainAddrLength)
	opAmount := big.NewInt(100)
	assetAddr := usdtAddress
	method := "depositTo"

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

	prepareFunc := func(params *types3.Params, method string) contracts.CallArgs {
		err := s.app.DepositKeeper.SetParams(s.ctx, params)
		s.Require().NoError(err)
		defaultDepositArgs := defaultCallArgs.WithMethodName(method)
		return defaultDepositArgs.WithArgs(
			uint16(clientChainLzId),
			assetAddr,
			stakerAddr,
			opAmount)
	}

	// test caller error
	beforeEach()
	setDepositToArgs := prepareFunc(&depositParams, method)
	_, _, err := contracts.CallContractAndCheckLogs(s.ctx, s.app, setDepositToArgs, passCheck)
	s.Require().ErrorContains(err, strings.Split(deposit.ErrContractCaller, ",")[0])

	// test success
	beforeEach()
	depositParams.ExoCoreLzAppAddress = s.address.String()
	setDepositToArgs = prepareFunc(&depositParams, method)
	_, ethRes, err := contracts.CallContractAndCheckLogs(s.ctx, s.app, setDepositToArgs, passCheck)
	successRet, err := s.precompile.Methods[deposit.MethodDepositTo].Outputs.Pack(true, opAmount)
	s.Require().NoError(err)
	s.Require().Equal(successRet, ethRes.Ret)
}

func (s *PrecompileTestSuite) TestCallDepositToFromContract() {
	// deposit params for test
	exoCoreLzAppAddress := "0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD"
	exoCoreLzAppEventTopic := "0xc6a377bfc4eb120024a8ac08eef205be16b817020812c73223e81d1bdb9708ec"
	depositParams := types3.Params{
		ExoCoreLzAppAddress:    exoCoreLzAppAddress,
		ExoCoreLzAppEventTopic: exoCoreLzAppEventTopic,
	}
	usdtAddress := paddingClientChainAddress(common.FromHex("0xdAC17F958D2ee523a2206206994597C13D831ec7"), types.GeneralClientChainAddrLength)
	clientChainLzId := 101
	stakerAddr := paddingClientChainAddress(s.address.Bytes(), types.GeneralClientChainAddrLength)
	opAmount := big.NewInt(100)
	assetAddr := usdtAddress

	// contractAddr is the address of the smart contract that will be deployed
	var contractAddr common.Address
	var err error

	// deploy the caller contract
	s.SetupTest()
	contractAddr, err = s.DeployContract(contracts.DepositCallerContract)
	s.Require().NoError(err)
	// NextBlock the smart contract
	s.NextBlock()
	// check contract was correctly deployed
	cAcc := s.app.EvmKeeper.GetAccount(s.ctx, contractAddr)
	s.Require().NotNil(cAcc)
	s.Require().True(cAcc.IsContract())

	beforeEach := func() {
		s.SetupTest()
		// populate default call args
		defaultCallArgs = contracts.CallArgs{
			ContractAddr: contractAddr,
			ContractABI:  contracts.DepositCallerContract.ABI,
			PrivKey:      s.privKey,
		}

		// default log check arguments
		defaultLogCheck = testutil.LogCheckArgs{ABIEvents: contracts.DepositCallerContract.ABI.Events}
		passCheck = defaultLogCheck.WithExpPass(true)
	}

	prepareFunc := func(params *types3.Params, method string) contracts.CallArgs {
		err := s.app.DepositKeeper.SetParams(s.ctx, params)
		s.Require().NoError(err)
		defaultDepositArgs := defaultCallArgs.WithMethodName(method)
		return defaultDepositArgs.WithArgs(
			uint16(clientChainLzId),
			assetAddr,
			stakerAddr,
			opAmount)
	}

	// testDepositTo
	beforeEach()
	depositParams.ExoCoreLzAppAddress = contractAddr.String()
	setDepositToArgs := prepareFunc(&depositParams, "testDepositTo")
	_, _, err = contracts.CallContractAndCheckLogs(s.ctx, s.app, setDepositToArgs, passCheck)
	s.Require().NoError(err)
	//todo: need to find why the ethRet is nil when called by contract
	/*	successRet, err := contracts.DepositCallerContract.ABI.Methods["testDepositTo"].Outputs.Pack(true, opAmount)
		s.Require().NoError(err)
		s.Require().Equal(successRet, ethRes.Ret)*/

	// testCallDepositToAndEmitEvent
	beforeEach()
	setDepositToArgs = prepareFunc(&depositParams, "testCallDepositToAndEmitEvent")
	// todo: need to check why can't get the ethereum log
	// eventCheck := passCheck.WithExpEvents("callDepositToResult")
	_, _, err = contracts.CallContractAndCheckLogs(s.ctx, s.app, setDepositToArgs, passCheck)
	s.Require().NoError(err)
	/*	successRet, err = contracts.DepositCallerContract.ABI.Methods["testCallDepositToAndEmitEvent"].Outputs.Pack(true, opAmount)
		s.Require().NoError(err)
		s.Require().Equal(successRet, ethRes.Ret)*/

	// testCallDepositToWithTryCatch
	beforeEach()
	depositParams.ExoCoreLzAppAddress = exoCoreLzAppAddress
	setDepositToArgs = prepareFunc(&depositParams, "testCallDepositToWithTryCatch")
	// eventCheck = passCheck.WithExpEvents("ErrorOccurred")
	// todo: need to check the ethereum log
	_, _, err = contracts.CallContractAndCheckLogs(s.ctx, s.app, setDepositToArgs, passCheck)
	s.Require().NoError(err)
	/*	successRet, err = contracts.DepositCallerContract.ABI.Methods["testCallDepositToWithTryCatch"].Outputs.Pack(false, big.NewInt(0))
		s.Require().NoError(err)
		s.Require().Equal(successRet, ethRes.Ret)*/
}
