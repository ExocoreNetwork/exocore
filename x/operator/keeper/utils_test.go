package keeper_test

import (
	"encoding/json"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"golang.org/x/exp/rand"
	"time"

	"github.com/ExocoreNetwork/exocore/testutil"
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
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	cmn "github.com/evmos/evmos/v14/precompiles/common"
	evmostypes "github.com/evmos/evmos/v14/types"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"
	inflationtypes "github.com/evmos/evmos/v14/x/inflation/types"
)

// SetupWithGenesisValSet initializes a new EvmosApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func (suite *KeeperTestSuite) SetupWithGenesisValSet(valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) {
	pruneOpts := pruningtypes.NewPruningOptionsFromString(pruningtypes.PruningOptionDefault)
	appI, genesisState := exocoreapp.SetupTestingApp(cmn.DefaultChainID, &pruneOpts, false)()
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
	suite.validators = validators

	// set validators and delegations
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
			ChainId:         cmn.DefaultChainID,
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
		cmn.DefaultChainID,
		sdk.ConsAddress(validators[0].GetOperator()),
		tmhash.Sum([]byte("app")),
		tmhash.Sum([]byte("validators")),
	)

	app.BeginBlock(abci.RequestBeginBlock{
		Header: header,
	})

	// need to create UncachedContext when retrieving historical state
	suite.ctx = app.BaseApp.NewUncachedContext(false, header)
	suite.app = app
}

func (suite *KeeperTestSuite) DoSetupTest() {
	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	suite.Require().NoError(err)

	privVal2 := mock.NewPV()
	pubKey2, err := privVal2.GetPubKey()
	suite.Require().NoError(err)

	// create validator set with two validators
	validator := tmtypes.NewValidator(pubKey, 1)
	validator2 := tmtypes.NewValidator(pubKey2, 2)
	suite.valSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{validator, validator2})
	signers := make(map[string]tmtypes.PrivValidator)
	signers[pubKey.Address().String()] = privVal
	signers[pubKey2.Address().String()] = privVal2

	// accAddress
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}
	rand.Read(pub.Key)
	suite.accAddress = sdk.AccAddress(pub.Address())

	// generate genesis account
	addr, priv := testutiltx.NewAddrKey()
	suite.privKey = priv
	suite.address = addr
	suite.signer = testutiltx.NewSigner(priv)

	baseAcc := authtypes.NewBaseAccount(priv.PubKey().Address().Bytes(), priv.PubKey(), 0, 0)

	acc := &evmostypes.EthAccount{
		BaseAccount: baseAcc,
		CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
	}

	amount := sdk.TokensFromConsensusPower(5, evmostypes.PowerReduction)

	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, amount)),
	}

	suite.SetupWithGenesisValSet(suite.valSet, []authtypes.GenesisAccount{acc}, balance)

	// Create StateDB
	suite.stateDB = statedb.New(suite.ctx, suite.app.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.ctx.HeaderHash().Bytes())))

	// bond denom
	stakingParams := suite.app.StakingKeeper.GetParams(suite.ctx)
	stakingParams.BondDenom = utils.BaseDenom
	suite.bondDenom = stakingParams.BondDenom
	err = suite.app.StakingKeeper.SetParams(suite.ctx, stakingParams)
	suite.Require().NoError(err)

	suite.ethSigner = ethtypes.LatestSignerForChainID(suite.app.EvmKeeper.ChainID())

	coins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(5000000000000000000)))
	inflCoins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(2000000000000000000)))
	distrCoins := sdk.NewCoins(sdk.NewCoin(utils.BaseDenom, sdk.NewInt(3000000000000000000)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, inflationtypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, inflationtypes.ModuleName, authtypes.FeeCollectorName, inflCoins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, inflationtypes.ModuleName, distrtypes.ModuleName, distrCoins)
	suite.Require().NoError(err)

}

// NextBlock commits the current block and sets up the next block.
func (suite *KeeperTestSuite) NextBlock() {
	var err error
	suite.ctx, err = testutil.CommitAndCreateNewCtx(suite.ctx, suite.app, time.Second, suite.valSet, true)
	suite.Require().NoError(err)
}
