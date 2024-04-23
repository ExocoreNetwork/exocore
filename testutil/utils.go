package testutil

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/evmos/evmos/v14/testutil"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/rand"

	testutiltx "github.com/ExocoreNetwork/exocore/testutil/tx"

	exocoreapp "github.com/ExocoreNetwork/exocore/app"
	"github.com/ExocoreNetwork/exocore/utils"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	dogfoodtypes "github.com/ExocoreNetwork/exocore/x/dogfood/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/tmhash"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmostypes "github.com/evmos/evmos/v14/types"
	"github.com/evmos/evmos/v14/x/evm/statedb"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

type BaseTestSuite struct {
	suite.Suite

	Ctx        sdk.Context
	App        *exocoreapp.ExocoreApp
	Address    common.Address
	AccAddress sdk.AccAddress

	PrivKey   cryptotypes.PrivKey
	Signer    keyring.Signer
	EthSigner ethtypes.Signer

	// construct genesis state from this info
	// x/assets
	ClientChains []assetstypes.ClientChainInfo
	Assets       []assetstypes.AssetInfo
	// for tracking validator across blocks
	ValSet *tmtypes.ValidatorSet

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
func (suite *BaseTestSuite) SetupWithGenesisValSet(genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) {
	pruneOpts := pruningtypes.NewPruningOptionsFromString(pruningtypes.PruningOptionDefault)
	appI, genesisState := exocoreapp.SetupTestingApp(utils.DefaultChainID, &pruneOpts, false)()
	app, ok := appI.(*exocoreapp.ExocoreApp)
	suite.Require().True(ok)

	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	// x/operator initialization - address only
	operator1 := sdk.AccAddress(testutiltx.GenerateAddress().Bytes())
	operator2 := sdk.AccAddress(testutiltx.GenerateAddress().Bytes())
	stakerID1, _ := assetstypes.GetStakeIDAndAssetIDFromStr(
		suite.ClientChains[0].LayerZeroChainID,
		common.Address(operator1.Bytes()).String(), "",
	)
	stakerID2, _ := assetstypes.GetStakeIDAndAssetIDFromStr(
		suite.ClientChains[0].LayerZeroChainID,
		common.Address(operator2.Bytes()).String(), "",
	)
	_, assetID := assetstypes.GetStakeIDAndAssetIDFromStr(
		suite.ClientChains[0].LayerZeroChainID,
		"", suite.Assets[0].Address,
	)
	// x/assets initialization - deposits (client chains and tokens are from caller)
	depositAmount := sdk.TokensFromConsensusPower(1, evmostypes.PowerReduction)
	depositsByStaker := []assetstypes.DepositsByStaker{
		{
			StakerID: stakerID1,
			Deposits: []assetstypes.DepositByAsset{
				{
					AssetID: assetID,
					Info: assetstypes.StakerAssetInfo{
						TotalDepositAmount:  depositAmount,
						WithdrawableAmount:  depositAmount,
						WaitUnbondingAmount: sdk.ZeroInt(),
					},
				},
			},
		},
		{
			StakerID: stakerID2,
			Deposits: []assetstypes.DepositByAsset{
				{
					AssetID: assetID,
					Info: assetstypes.StakerAssetInfo{
						TotalDepositAmount:  depositAmount,
						WithdrawableAmount:  depositAmount,
						WaitUnbondingAmount: sdk.ZeroInt(),
					},
				},
			},
		},
	}
	assetsGenesis := assetstypes.NewGenesis(
		assetstypes.DefaultParams(),
		suite.ClientChains, []assetstypes.StakingAssetInfo{
			{
				AssetBasicInfo: &suite.Assets[0],
				// required to be 0, since deposits are handled after token init.
				StakingTotalAmount: sdk.ZeroInt(),
			},
		}, depositsByStaker,
	)
	genesisState[assetstypes.ModuleName] = app.AppCodec().MustMarshalJSON(assetsGenesis)

	// operator registration
	operatorInfos := []operatortypes.OperatorInfo{
		operatortypes.OperatorInfo{
			EarningsAddr:     operator1.String(),
			OperatorMetaInfo: "operator1",
			Commission:       stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		},
		operatortypes.OperatorInfo{
			EarningsAddr:     operator2.String(),
			OperatorMetaInfo: "operator2",
			Commission:       stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		},
	}
	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	suite.Require().NoError(err)
	privVal2 := mock.NewPV()
	pubKey2, err := privVal2.GetPubKey()
	suite.Require().NoError(err)
	// operator consensus keys
	consensusKeyRecords := []operatortypes.OperatorConsKeyRecord{
		operatortypes.OperatorConsKeyRecord{
			OperatorAddress: operatorInfos[0].EarningsAddr,
			Chains: []operatortypes.ChainDetails{
				{
					ChainID:      utils.DefaultChainID,
					ConsensusKey: hexutil.Encode(pubKey.Bytes()),
				},
			},
		},
		operatortypes.OperatorConsKeyRecord{
			OperatorAddress: operatorInfos[1].EarningsAddr,
			Chains: []operatortypes.ChainDetails{
				{
					ChainID:      utils.DefaultChainID,
					ConsensusKey: hexutil.Encode(pubKey2.Bytes()),
				},
			},
		},
	}
	operatorGenesis := operatortypes.NewGenesisState(operatorInfos, consensusKeyRecords)
	genesisState[operatortypes.ModuleName] = app.AppCodec().MustMarshalJSON(operatorGenesis)

	// x/delegation
	delegationsByStaker := []delegationtypes.DelegationsByStaker{
		{
			StakerID: stakerID1,
			Delegations: []delegationtypes.DelegatedSingleAssetInfo{
				{
					AssetID: assetID,
					PerOperatorAmounts: []delegationtypes.KeyValue{
						{
							Key: operator1.String(),
							Value: &delegationtypes.ValueField{
								Amount: depositAmount,
							},
						},
					},
				},
			},
		},
		{
			StakerID: stakerID2,
			Delegations: []delegationtypes.DelegatedSingleAssetInfo{
				{
					AssetID: assetID,
					PerOperatorAmounts: []delegationtypes.KeyValue{
						{
							Key: operator2.String(),
							Value: &delegationtypes.ValueField{
								Amount: depositAmount,
							},
						},
					},
				},
			},
		},
	}
	delegationGenesis := delegationtypes.NewGenesis(delegationsByStaker)
	genesisState[delegationtypes.ModuleName] = app.AppCodec().MustMarshalJSON(delegationGenesis)

	dogfoodGenesis := dogfoodtypes.NewGenesis(
		dogfoodtypes.DefaultParams(), []dogfoodtypes.GenesisValidator{
			{
				PublicKey: hexutil.Encode(pubKey.Bytes()),
				Power:     1,
			},
			{
				PublicKey: hexutil.Encode(pubKey.Bytes()),
				Power:     1,
			},
		},
	)
	genesisState[dogfoodtypes.ModuleName] = app.AppCodec().MustMarshalJSON(dogfoodGenesis)

	suite.ValSet = tmtypes.NewValidatorSet([]*tmtypes.Validator{
		tmtypes.NewValidator(pubKey, 1),
		tmtypes.NewValidator(pubKey2, 1),
	})

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultParams(), balances, totalSupply,
		[]banktypes.Metadata{}, []banktypes.SendEnabled{},
	)
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
	convKey, err := cryptocodec.FromTmPubKeyInterface(pubKey)
	suite.Require().NoError(err)
	header := testutil.NewHeader(
		2,
		time.Now().UTC(),
		utils.DefaultChainID,
		sdk.GetConsAddress(convKey),
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
	// create AccAddress for test
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}
	_, err := rand.Read(pub.Key)
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

	// Exocore modules genesis
	// x/assets
	suite.ClientChains = []assetstypes.ClientChainInfo{
		assetstypes.ClientChainInfo{
			Name:               "ethereum",
			MetaInfo:           "ethereum blockchain",
			ChainId:            1,
			FinalizationBlocks: 10,
			LayerZeroChainID:   101,
			AddressLength:      20,
		},
	}
	suite.Assets = []assetstypes.AssetInfo{
		assetstypes.AssetInfo{
			Name:             "Tether USD",
			Symbol:           "USDT",
			Address:          "0xdAC17F958D2ee523a2206206994597C13D831ec7",
			Decimals:         6,
			LayerZeroChainID: suite.ClientChains[0].LayerZeroChainID,
			MetaInfo:         "Tether USD token",
		},
	}
	{
		totalSupply, _ := sdk.NewIntFromString("40022689732746729")
		suite.Assets[0].TotalSupply = totalSupply
	}

	// Initialize an ExocoreApp for test
	suite.SetupWithGenesisValSet(
		[]authtypes.GenesisAccount{acc}, balance,
	)

	// Create StateDB
	suite.StateDB = statedb.New(suite.Ctx, suite.App.EvmKeeper, statedb.NewEmptyTxConfig(common.BytesToHash(suite.Ctx.HeaderHash().Bytes())))

	suite.EthSigner = ethtypes.LatestSignerForChainID(suite.App.EvmKeeper.ChainID())

	queryHelperEvm := baseapp.NewQueryServerTestHelper(suite.Ctx, suite.App.InterfaceRegistry())
	evmtypes.RegisterQueryServer(queryHelperEvm, suite.App.EvmKeeper)
	suite.QueryClientEVM = evmtypes.NewQueryClient(queryHelperEvm)
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
	suite.Ctx, err = CommitAndCreateNewCtx(suite.Ctx, suite.App, time.Second, suite.ValSet, true)
	suite.Require().NoError(err)
}
