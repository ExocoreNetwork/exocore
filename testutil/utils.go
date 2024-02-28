package testutil

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/evmos/evmos/v14/testutil"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"

	testutiltx "github.com/ExocoreNetwork/exocore/testutil/tx"

	exocoreapp "github.com/ExocoreNetwork/exocore/app"
	"github.com/ExocoreNetwork/exocore/utils"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	tmtypes "github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmostypes "github.com/evmos/evmos/v14/types"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
)

type BaseTestSuite struct {
	suite.Suite

	Ctx        sdk.Context
	App        *exocoreapp.ExocoreApp
	Address    common.Address
	AccAddress sdk.AccAddress

	Validators     []stakingtypes.Validator
	ValSet         *tmtypes.ValidatorSet
	EthSigner      ethtypes.Signer
	PrivKey        cryptotypes.PrivKey
	Signer         keyring.Signer
	BondDenom      string
	StateDB        *statedb.StateDB
	QueryClientEVM evmtypes.QueryClient

	// needed by test
	/*	operatorAddr          sdk.AccAddress
		avsAddr               string
		assetID               string
		stakerID              string
		assetAddr             common.Address
		assetDecimal          uint32
		clientChainLzID       uint64
		depositAmount         sdkmath.Int
		delegationAmount      sdkmath.Int
		updatedAmountForOptIn sdkmath.Int*/
}

func (suite *BaseTestSuite) SetupTest() {
	suite.DoSetupTest()
}

// SetupWithGenesisValSet initializes a new EvmosApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func (suite *BaseTestSuite) SetupWithGenesisValSet(valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) {
	appI, genesisState := exocoreapp.SetupTestingApp(utils.DefaultChainID, false)()
	app, ok := appI.(*exocoreapp.ExocoreApp)
	suite.Require().True(ok)

	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.TokensFromConsensusPower(1, evmostypes.PowerReduction)

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		suite.Require().NoError(err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		suite.Require().NoError(err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))
	}
	suite.Validators = validators

	// set Validators and delegations
	stakingParams := stakingtypes.DefaultParams()
	// set bond demon to be aevmos
	stakingParams.BondDenom = utils.BaseDenom
	stakingGenesis := stakingtypes.NewGenesisState(stakingParams, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalBondAmt := bondAmt.Add(bondAmt)
	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin(utils.BaseDenom, totalBondAmt))...)
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(utils.BaseDenom, totalBondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	suite.Require().NoError(err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         utils.DefaultChainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: exocoreapp.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	app.Commit()

	// instantiate new header
	header := testutil.NewHeader(
		2,
		time.Now().UTC(),
		utils.DefaultChainID,
		sdk.ConsAddress(validators[0].GetOperator()),
		tmhash.Sum([]byte("App")),
		tmhash.Sum([]byte("Validators")),
	)

	app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// need to create UncachedContext when retrieving historical state
	suite.Ctx = app.BaseApp.NewUncachedContext(false, header)
	suite.App = app
}

func (suite *BaseTestSuite) DoSetupTest() {
	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	suite.Require().NoError(err)

	privVal2 := mock.NewPV()
	pubKey2, err := privVal2.GetPubKey()
	suite.Require().NoError(err)

	// create validator set with two Validators
	validator := tmtypes.NewValidator(pubKey, 1)
	validator2 := tmtypes.NewValidator(pubKey2, 2)
	suite.ValSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator, validator2})
	signers := make(map[string]tmtypes.PrivValidator)
	signers[pubKey.Address().String()] = privVal
	signers[pubKey2.Address().String()] = privVal2

	// create AccAddress for test
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}
	_, err = rand.Read(pub.Key)
	suite.Require().NoError(err)
	suite.AccAddress = sdk.AccAddress(pub.Address())

	// generate genesis account
	addr, priv := testutiltx.NewAddrKey()
	suite.PrivKey = priv
	suite.Address = addr
	suite.Signer = testutiltx.NewSigner(priv)
	baseAcc := authtypes.NewBaseAccount(priv.PubKey().Address().Bytes(), priv.PubKey(), 0, 0)
	acc := &evmostypes.EthAccount{
		BaseAccount: baseAcc,
		CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
	}
	// set amount for genesis account
	amount := sdk.TokensFromConsensusPower(5, evmostypes.PowerReduction)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, amount)),
	}

	// Initialize an ExocoreApp for test
	suite.SetupWithGenesisValSet(suite.ValSet, []authtypes.GenesisAccount{acc}, balance)

	// Create StateDB
	suite.StateDB = statedb.New(suite.Ctx, suite.App.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.Ctx.HeaderHash().Bytes())))

	// bond denom
	/*	stakingParams := suite.App.StakingKeeper.GetParams(suite.Ctx)
		stakingParams.BondDenom = utils.BaseDenom
		suite.BondDenom = stakingParams.BondDenom
		err = suite.App.StakingKeeper.SetParams(suite.Ctx, stakingParams)
		suite.Require().NoError(err)*/

	suite.BondDenom = utils.BaseDenom
	suite.EthSigner = ethtypes.LatestSignerForChainID(suite.App.EvmKeeper.ChainID())

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelperEvm, suite.App.EvmKeeper)
	suite.QueryClientEVM = evmtypes.NewQueryClient(queryHelperEvm)

	/*	coins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(5000000000000000000)))
		inflCoins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(2000000000000000000)))
		distrCoins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(3000000000000000000)))
		err = suite.App.BankKeeper.MintCoins(suite.Ctx, inflationtypes.ModuleName, coins)
		suite.Require().NoError(err)
		err = suite.App.BankKeeper.SendCoinsFromModuleToModule(suite.Ctx, inflationtypes.ModuleName, authtypes.FeeCollectorName, inflCoins)
		suite.Require().NoError(err)
		err = suite.App.BankKeeper.SendCoinsFromModuleToModule(suite.Ctx, inflationtypes.ModuleName, distrtypes.ModuleName, distrCoins)
		suite.Require().NoError(err)*/
}

// DeployContract deploys a contract that calls the deposit precompile's methods for testing purposes.
func (suite *BaseTestSuite) DeployContract(contract evmtypes.CompiledContract) (addr common.Address, err error) {
	addr, err = DeployContract(
		suite.Ctx,
		suite.App,
		suite.PrivKey,
		suite.QueryClientEVM,
		contract,
	)
	return
}

// NextBlock commits the current block and sets up the next block.
func (suite *BaseTestSuite) NextBlock() {
	var err error
	suite.Ctx, err = CommitAndCreateNewCtx(suite.Ctx, suite.App, time.Second, suite.ValSet)
	suite.Require().NoError(err)
}
