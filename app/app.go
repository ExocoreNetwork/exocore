package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/ExocoreNetwork/exocore/x/avstask"
	avsTaskKeeper "github.com/ExocoreNetwork/exocore/x/avstask/keeper"
	avsTaskTypes "github.com/ExocoreNetwork/exocore/x/avstask/types"

	"github.com/ExocoreNetwork/exocore/x/oracle"
	oracleKeeper "github.com/ExocoreNetwork/exocore/x/oracle/keeper"
	oracleTypes "github.com/ExocoreNetwork/exocore/x/oracle/types"

	"github.com/ExocoreNetwork/exocore/x/avs"
	"github.com/ExocoreNetwork/exocore/x/operator"
	operatorKeeper "github.com/ExocoreNetwork/exocore/x/operator/keeper"

	exoslash "github.com/ExocoreNetwork/exocore/x/slash"

	avsManagerKeeper "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	avsManagerTypes "github.com/ExocoreNetwork/exocore/x/avs/types"
	slashKeeper "github.com/ExocoreNetwork/exocore/x/slash/keeper"
	exoslashTypes "github.com/ExocoreNetwork/exocore/x/slash/types"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	// for EIP-1559 fee handling
	ethante "github.com/evmos/evmos/v14/app/ante/evm"
	// for encoding and decoding of EIP-712 messages
	"github.com/evmos/evmos/v14/ethereum/eip712"

	"github.com/ExocoreNetwork/exocore/x/assets"
	assetsKeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	assetsTypes "github.com/ExocoreNetwork/exocore/x/assets/types"

	"github.com/ExocoreNetwork/exocore/x/delegation"
	delegationKeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	delegationTypes "github.com/ExocoreNetwork/exocore/x/delegation/types"

	"github.com/ExocoreNetwork/exocore/x/deposit"
	depositKeeper "github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	depositTypes "github.com/ExocoreNetwork/exocore/x/deposit/types"

	operatorTypes "github.com/ExocoreNetwork/exocore/x/operator/types"

	"github.com/ExocoreNetwork/exocore/x/reward"
	rewardKeeper "github.com/ExocoreNetwork/exocore/x/reward/keeper"
	rewardTypes "github.com/ExocoreNetwork/exocore/x/reward/types"

	"github.com/ExocoreNetwork/exocore/x/withdraw"
	withdrawKeeper "github.com/ExocoreNetwork/exocore/x/withdraw/keeper"
	withdrawTypes "github.com/ExocoreNetwork/exocore/x/withdraw/types"

	// increases or decreases block gas limit based on usage
	"github.com/evmos/evmos/v14/x/feemarket"
	feemarketkeeper "github.com/evmos/evmos/v14/x/feemarket/keeper"
	feemarkettypes "github.com/evmos/evmos/v14/x/feemarket/types"

	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"

	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"

	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/evmos/evmos/v14/app/ante"

	"github.com/evmos/evmos/v14/precompiles/common"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"

	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctestingtypes "github.com/cosmos/ibc-go/v7/testing/types"

	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	// this module allows the transfer of ERC20 tokens over IBC. for such transfers to occur,
	// they must be enabled in the ERC20 keeper.
	transfer "github.com/evmos/evmos/v14/x/ibc/transfer"
	transferkeeper "github.com/evmos/evmos/v14/x/ibc/transfer/keeper"

	"cosmossdk.io/simapp"
	simappparams "cosmossdk.io/simapp/params"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	srvflags "github.com/evmos/evmos/v14/server/flags"

	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"

	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	"github.com/cosmos/cosmos-sdk/x/consensus"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"

	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/ExocoreNetwork/exocore/x/evm"
	evmkeeper "github.com/ExocoreNetwork/exocore/x/evm/keeper"
	evmtypes "github.com/evmos/evmos/v14/x/evm/types"

	"github.com/evmos/evmos/v14/encoding"
	evmostypes "github.com/evmos/evmos/v14/types"

	// The recovery module is an IBC middleware for helping users recover funds that they sent
	// to the Cosmos secp256k1 address instead of the Ethereum ethsecp256k1 address. It only
	// works for authorized chains.
	"github.com/evmos/evmos/v14/x/recovery"
	recoverykeeper "github.com/evmos/evmos/v14/x/recovery/keeper"
	recoverytypes "github.com/evmos/evmos/v14/x/recovery/types"

	"github.com/evmos/evmos/v14/x/epochs"
	epochskeeper "github.com/evmos/evmos/v14/x/epochs/keeper"
	epochstypes "github.com/evmos/evmos/v14/x/epochs/types"

	"github.com/evmos/evmos/v14/x/erc20"
	erc20keeper "github.com/evmos/evmos/v14/x/erc20/keeper"
	erc20types "github.com/evmos/evmos/v14/x/erc20/types"

	// unnamed import of statik for swagger UI support
	_ "github.com/evmos/evmos/v14/client/docs/statik"

	// Force-load the tracer engines to trigger registration due to Go-Ethereum v1.10.15 changes
	_ "github.com/ethereum/go-ethereum/eth/tracers/js"
	_ "github.com/ethereum/go-ethereum/eth/tracers/native"
)

// Name defines the application binary name
const Name = "exocored"

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+Name)

	// manually update the power reduction by replacing micro (u) -> atto (a) evmos
	sdk.DefaultPowerReduction = evmostypes.PowerReduction
	// modify fee market parameter defaults through global
	feemarkettypes.DefaultMinGasPrice = MainnetMinGasPrices
	feemarkettypes.DefaultMinGasMultiplier = MainnetMinGasMultiplier
	// modify default min commission to 5%
	stakingtypes.DefaultMinCommissionRate = sdk.NewDecWithPrec(5, 2)
}

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler, upgradeclient.LegacyProposalHandler, upgradeclient.LegacyCancelProposalHandler,
				ibcclientclient.UpdateClientProposalHandler, ibcclientclient.UpgradeProposalHandler,
				// Evmos proposal types
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		ica.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{AppModuleBasic: &ibctransfer.AppModuleBasic{}},
		evm.AppModuleBasic{},
		feemarket.AppModuleBasic{},
		// evmos modules
		erc20.AppModuleBasic{},
		epochs.AppModuleBasic{},
		recovery.AppModuleBasic{},
		consensus.AppModuleBasic{},
		// Exocore modules
		assets.AppModuleBasic{},
		deposit.AppModuleBasic{},
		operator.AppModuleBasic{},
		delegation.AppModuleBasic{},
		withdraw.AppModuleBasic{},
		reward.AppModuleBasic{},
		exoslash.AppModuleBasic{},
		avs.AppModuleBasic{},
		avstask.AppModuleBasic{},
		oracle.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		icatypes.ModuleName:            nil,
		evmtypes.ModuleName: {
			authtypes.Minter,
			authtypes.Burner,
		}, // used for secure addition and subtraction of balance using module account
		erc20types.ModuleName: {authtypes.Minter, authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{}
)

var (
	_ servertypes.Application = (*ExocoreApp)(nil)
	_ ibctesting.TestingApp   = (*ExocoreApp)(nil)
)

// ExocoreApp implements an extended ABCI application. It is an application
// that may process transactions through Ethereum's EVM running atop of
// Tendermint consensus.
type ExocoreApp struct {
	*baseapp.BaseApp

	// encoding
	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	CrisisKeeper          crisiskeeper.Keeper
	UpgradeKeeper         upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	ICAHostKeeper         icahostkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        transferkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	// Ethermint keepers
	EvmKeeper       *evmkeeper.Keeper
	FeeMarketKeeper feemarketkeeper.Keeper

	// Evmos keepers
	Erc20Keeper    erc20keeper.Keeper
	EpochsKeeper   epochskeeper.Keeper
	RecoveryKeeper *recoverykeeper.Keeper

	// exocore assets module keepers
	AssetsKeeper     assetsKeeper.Keeper
	DepositKeeper    depositKeeper.Keeper
	DelegationKeeper delegationKeeper.Keeper
	WithdrawKeeper   withdrawKeeper.Keeper
	RewardKeeper     rewardKeeper.Keeper
	OperatorKeeper   operatorKeeper.Keeper
	ExoSlashKeeper   slashKeeper.Keeper
	AVSManagerKeeper avsManagerKeeper.Keeper
	TaskKeeper       avsTaskKeeper.Keeper
	OracleKeeper     oracleKeeper.Keeper

	// the module manager
	mm *module.Manager

	// the configurator
	configurator module.Configurator

	// simulation manager
	sm *module.SimulationManager

	tpsCounter *tpsCounter
}

// SimulationManager implements runtime.AppI
func (*ExocoreApp) SimulationManager() *module.SimulationManager {
	panic("unimplemented")
}

// NewExocoreApp is the constructor for new Exocore
func NewExocoreApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig simappparams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *ExocoreApp {
	appCodec := encodingConfig.Codec
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	eip712.SetEncodingConfig(encodingConfig)

	// Setup Mempool and Proposal Handlers
	baseAppOptions = append(baseAppOptions, func(app *baseapp.BaseApp) {
		mempool := mempool.NoOpMempool{}
		app.SetMempool(mempool)
		handler := baseapp.NewDefaultProposalHandler(mempool, app)
		app.SetPrepareProposal(handler.PrepareProposalHandler())
		app.SetProcessProposal(handler.ProcessProposalHandler())
	})
	// NOTE we use custom transaction decoder that supports the sdk.Tx interface instead of
	// sdk.StdTx
	bApp := baseapp.NewBaseApp(
		Name,
		logger,
		db,
		encodingConfig.TxConfig.TxDecoder(),
		baseAppOptions...,
	)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		// SDK keys
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, capabilitytypes.StoreKey, consensusparamtypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey, crisistypes.StoreKey,
		// ibc keys
		ibcexported.StoreKey, ibctransfertypes.StoreKey,
		// ica keys
		icahosttypes.StoreKey,
		// ethermint keys
		evmtypes.StoreKey, feemarkettypes.StoreKey,
		// evmos keys
		erc20types.StoreKey,
		epochstypes.StoreKey,
		recoverytypes.StoreKey,
		// exoCore module keys
		assetsTypes.StoreKey,
		delegationTypes.StoreKey,
		depositTypes.StoreKey,
		withdrawTypes.StoreKey,
		rewardTypes.StoreKey,
		exoslashTypes.StoreKey,
		operatorTypes.StoreKey,
		avsManagerTypes.StoreKey,
		avsTaskTypes.StoreKey,
		oracleTypes.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey, evmtypes.TransientKey, feemarkettypes.TransientKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, oracleTypes.MemStoreKey)

	// load state streaming if enabled
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, logger, keys); err != nil {
		fmt.Printf("failed to load state streaming: %s", err)
		os.Exit(1)
	}

	app := &ExocoreApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	// init params keeper and subspaces
	app.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	// get authority address
	authAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// set the BaseApp's parameter store
	app.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(
		appCodec,
		keys[consensusparamtypes.StoreKey],
		authAddr,
	)
	bApp.SetParamStore(&app.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)

	// Applications that wish to enforce statically created ScopedKeepers should call `Seal`
	// after creating
	// their scoped modules in `NewApp` with `ScopeToModule`
	app.CapabilityKeeper.Seal()
	// add account keepers
	accountK := authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey],
		evmostypes.ProtoAccount, maccPerms,
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		authAddr,
	)
	app.AccountKeeper = accountK

	// add bank keepers
	bankK := bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], accountK, app.BlockedAddrs(), authAddr,
	)
	app.BankKeeper = bankK

	// add keepers
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], accountK, bankK, authAddr,
	)
	app.StakingKeeper = *stakingKeeper
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, app.LegacyAmino(), keys[slashingtypes.StoreKey], stakingKeeper, authAddr,
	)
	app.CrisisKeeper = *crisiskeeper.NewKeeper(
		appCodec, keys[crisistypes.StoreKey], invCheckPeriod, bankK, authtypes.FeeCollectorName, authAddr,
	)
	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], accountK)
	app.UpgradeKeeper = *upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp, authAddr)

	app.AuthzKeeper = authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		appCodec,
		app.MsgServiceRouter(),
		app.AccountKeeper,
	)

	tracer := cast.ToString(appOpts.Get(srvflags.EVMTracer))

	// Create Ethermint keepers
	app.FeeMarketKeeper = feemarketkeeper.NewKeeper(
		appCodec, authtypes.NewModuleAddress(govtypes.ModuleName),
		keys[feemarkettypes.StoreKey],
		tkeys[feemarkettypes.TransientKey],
		app.GetSubspace(feemarkettypes.ModuleName),
	)

	evmKeeper := evmkeeper.NewKeeper(
		appCodec,
		keys[evmtypes.StoreKey],
		tkeys[evmtypes.TransientKey],
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		stakingKeeper,
		app.FeeMarketKeeper,
		tracer,
		app.GetSubspace(evmtypes.ModuleName),
	)

	app.EvmKeeper = evmKeeper

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		stakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(&app.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper)).
		AddRoute(erc20types.RouterKey, erc20.NewErc20ProposalHandler(&app.Erc20Keeper))
	govConfig := govtypes.DefaultConfig()
	/*
		Example of setting gov params:
		govConfig.MaxMetadataLen = 10000
	*/
	govKeeper := govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], accountK, bankK,
		stakingKeeper, app.MsgServiceRouter(), govConfig, authAddr,
	)

	// Set legacy router for backwards compatibility with gov v1beta1
	govKeeper.SetLegacyRouter(govRouter)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	// NOTE: Slashing must be created before calling the Hooks method to avoid
	// returning a Keeper without its table generated
	stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			app.SlashingKeeper.Hooks(),
		),
	)

	app.StakingKeeper = *stakingKeeper

	app.RecoveryKeeper = recoverykeeper.NewKeeper(
		keys[recoverytypes.StoreKey],
		appCodec,
		authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.TransferKeeper,
	)

	app.Erc20Keeper = erc20keeper.NewKeeper(
		keys[erc20types.StoreKey], appCodec, authtypes.NewModuleAddress(govtypes.ModuleName),
		app.AccountKeeper, app.BankKeeper, app.EvmKeeper, app.StakingKeeper, app.RecoveryKeeper,
	)

	app.TransferKeeper = transferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.RecoveryKeeper, // ICS4 Wrapper: recovery IBC middleware
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
		app.Erc20Keeper, // Add ERC20 Keeper for ERC20 transfers
	)

	// set exoCore staking keepers
	app.AssetsKeeper = assetsKeeper.NewKeeper(keys[assetsTypes.StoreKey], appCodec)
	app.DepositKeeper = depositKeeper.NewKeeper(
		keys[depositTypes.StoreKey],
		appCodec,
		app.AssetsKeeper,
	)
	app.OperatorKeeper = operatorKeeper.NewKeeper(
		keys[operatorTypes.StoreKey],
		appCodec,
		app.AssetsKeeper,
		operatorTypes.MockOracle{},
		operatorTypes.MockAvs{},
		delegationTypes.VirtualSlashKeeper{},
	)
	// todo: need to replace the virtual keepers with actual keepers after they have been
	// implemented
	app.DelegationKeeper = delegationKeeper.NewKeeper(
		keys[depositTypes.StoreKey],
		appCodec,
		app.AssetsKeeper,
		delegationTypes.VirtualSlashKeeper{},
		&app.OperatorKeeper,
	)
	app.OperatorKeeper.RegisterExpectDelegationInterface(&app.DelegationKeeper)

	app.WithdrawKeeper = *withdrawKeeper.NewKeeper(appCodec, keys[withdrawTypes.StoreKey], app.AssetsKeeper, app.DepositKeeper)
	app.RewardKeeper = *rewardKeeper.NewKeeper(appCodec, keys[rewardTypes.StoreKey], app.AssetsKeeper)
	app.ExoSlashKeeper = slashKeeper.NewKeeper(appCodec, keys[exoslashTypes.StoreKey], app.AssetsKeeper)
	app.AVSManagerKeeper = *avsManagerKeeper.NewKeeper(appCodec, keys[avsManagerTypes.StoreKey], &app.OperatorKeeper, app.AssetsKeeper)
	app.TaskKeeper = avsTaskKeeper.NewKeeper(appCodec, keys[avsTaskTypes.StoreKey], app.AVSManagerKeeper)
	app.OracleKeeper = oracleKeeper.NewKeeper(appCodec, keys[oracleTypes.StoreKey], memKeys[oracleTypes.MemStoreKey], app.GetSubspace(oracleTypes.ModuleName), app.StakingKeeper)
	// We call this after setting the hooks to ensure that the hooks are set on the keeper
	evmKeeper.WithPrecompiles(
		evmkeeper.AvailablePrecompiles(
			*stakingKeeper,
			app.AuthzKeeper,
			app.TransferKeeper,
			app.IBCKeeper.ChannelKeeper,
			app.DepositKeeper,
			app.DelegationKeeper,
			app.AssetsKeeper,
			app.WithdrawKeeper,
			app.ExoSlashKeeper,
			app.RewardKeeper,
			app.AVSManagerKeeper,
			app.TaskKeeper,
		),
	)
	epochsKeeper := epochskeeper.NewKeeper(appCodec, keys[epochstypes.StoreKey])
	app.EpochsKeeper = *epochsKeeper.SetHooks(
		epochskeeper.NewMultiEpochHooks(
		// insert epoch hooks receivers here
		),
	)

	app.GovKeeper = *govKeeper.SetHooks(
		govtypes.NewMultiGovHooks(),
	)

	app.EvmKeeper = app.EvmKeeper.SetHooks(
		evmkeeper.NewMultiEvmHooks(
			app.Erc20Keeper.Hooks(),
		),
	)

	// NOTE: app.Erc20Keeper is already initialized elsewhere

	// Set the ICS4 wrappers for custom module middlewares
	app.RecoveryKeeper.SetICS4Wrapper(app.IBCKeeper.ChannelKeeper)

	// Override the ICS20 app module
	transferModule := transfer.NewAppModule(app.TransferKeeper)

	// Create the app.ICAHostKeeper
	app.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, app.keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.RecoveryKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)

	// create host IBC module
	icaHostIBCModule := icahost.NewIBCModule(app.ICAHostKeeper)

	/*
		Create Transfer Stack

		transfer stack contains (from bottom to top):
			- ERC-20 Middleware
		 	- Recovery Middleware
			- IBC Transfer

		SendPacket, since it is originating from the application to core IBC:
		 	transferKeeper.SendPacket -> recovery.SendPacket -> erc20.SendPacket -> channel.SendPacket

		RecvPacket, message that originates from core IBC and goes down to app, the flow is the other way
			channel.RecvPacket -> erc20.OnRecvPacket -> recovery.OnRecvPacket -> transfer.OnRecvPacket
	*/

	// create IBC module from top to bottom of stack
	var transferStack porttypes.IBCModule

	transferStack = transfer.NewIBCModule(app.TransferKeeper)
	transferStack = recovery.NewIBCMiddleware(*app.RecoveryKeeper, transferStack)
	transferStack = erc20.NewIBCMiddleware(app.Erc20Keeper, transferStack)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.
		AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibctransfertypes.ModuleName, transferStack)

	app.IBCKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		// SDK app modules
		genutil.NewAppModule(
			accountK, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(
			appCodec,
			accountK,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),
		bank.NewAppModule(appCodec, bankK, accountK, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		crisis.NewAppModule(
			&app.CrisisKeeper,
			skipGenesisInvariants,
			app.GetSubspace(crisistypes.ModuleName),
		),
		gov.NewAppModule(
			appCodec,
			govKeeper,
			accountK,
			bankK,
			app.GetSubspace(govtypes.ModuleName),
		),
		slashing.NewAppModule(
			appCodec,
			app.SlashingKeeper,
			accountK,
			bankK,
			app.StakingKeeper,
			app.GetSubspace(slashingtypes.ModuleName),
		),
		staking.NewAppModule(
			appCodec,
			&app.StakingKeeper,
			accountK,
			bankK,
			app.GetSubspace(stakingtypes.ModuleName),
		),
		upgrade.NewAppModule(&app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		params.NewAppModule(app.ParamsKeeper),
		feegrantmodule.NewAppModule(
			appCodec,
			accountK,
			bankK,
			app.FeeGrantKeeper,
			app.interfaceRegistry,
		),
		authzmodule.NewAppModule(
			appCodec,
			app.AuthzKeeper,
			app.AccountKeeper,
			app.BankKeeper,
			app.interfaceRegistry,
		),
		consensus.NewAppModule(appCodec, app.ConsensusParamsKeeper),

		// ibc modules
		ibc.NewAppModule(app.IBCKeeper),
		ica.NewAppModule(nil, &app.ICAHostKeeper),
		transferModule,
		// Ethermint app modules
		evm.NewAppModule(
			app.EvmKeeper,
			app.AccountKeeper,
			app.GetSubspace(evmtypes.ModuleName),
		),
		feemarket.NewAppModule(app.FeeMarketKeeper, app.GetSubspace(feemarkettypes.ModuleName)),
		// Evmos app modules
		erc20.NewAppModule(app.Erc20Keeper, app.AccountKeeper,
			app.GetSubspace(erc20types.ModuleName)),
		epochs.NewAppModule(appCodec, app.EpochsKeeper),
		recovery.NewAppModule(*app.RecoveryKeeper,
			app.GetSubspace(recoverytypes.ModuleName)),
		// exoCore app modules
		assets.NewAppModule(appCodec, app.AssetsKeeper),
		deposit.NewAppModule(appCodec, app.DepositKeeper),
		operator.NewAppModule(appCodec, app.OperatorKeeper),
		delegation.NewAppModule(appCodec, app.DelegationKeeper),
		withdraw.NewAppModule(appCodec, app.WithdrawKeeper),
		reward.NewAppModule(appCodec, app.RewardKeeper),
		exoslash.NewAppModule(appCodec, app.ExoSlashKeeper),
		avs.NewAppModule(appCodec, app.AVSManagerKeeper),
		avstask.NewAppModule(appCodec, app.TaskKeeper),
		oracle.NewAppModule(appCodec, app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
	)

	// During begin block slashing happens after reward.BeginBlocker so that
	// there is nothing left over in the validator fee pool, to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: upgrade module must go first to handle software upgrades.
	// NOTE: staking module is required if HistoricalEntries param > 0.
	// NOTE: capability module's beginblocker must come before any modules using capabilities
	// (e.g. IBC)
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		// Note: epochs' begin should be "real" start of epochs, we keep epochs beginblock at
		// the beginning
		epochstypes.ModuleName,
		feemarkettypes.ModuleName,
		evmtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibcexported.ModuleName,
		// no-op modules
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		govtypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		erc20types.ModuleName,
		recoverytypes.ModuleName,
		consensusparamtypes.ModuleName,
		// ExoCore modules
		assetsTypes.ModuleName,
		depositTypes.ModuleName,
		operatorTypes.ModuleName,
		delegationTypes.ModuleName,
		withdrawTypes.ModuleName,
		rewardTypes.ModuleName,
		exoslashTypes.ModuleName,
		avsManagerTypes.ModuleName,
		avsTaskTypes.ModuleName,
		oracleTypes.ModuleName,
	)

	// NOTE: fee market module must go last in order to retrieve the block gas used.
	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		evmtypes.ModuleName,
		feemarkettypes.ModuleName,
		// Note: epochs' endblock should be "real" end of epochs, we keep epochs endblock at the
		// end
		epochstypes.ModuleName,
		// no-op modules
		ibcexported.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		// Evmos modules
		erc20types.ModuleName,
		recoverytypes.ModuleName,
		consensusparamtypes.ModuleName,
		// ExoCore modules
		assetsTypes.ModuleName,
		depositTypes.ModuleName,
		operatorTypes.ModuleName,
		delegationTypes.ModuleName,
		withdrawTypes.ModuleName,
		rewardTypes.ModuleName,
		exoslashTypes.ModuleName,
		avsManagerTypes.ModuleName,
		avsTaskTypes.ModuleName,
		oracleTypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		// SDK modules
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		ibcexported.ModuleName,
		// Ethermint modules
		evmtypes.ModuleName,
		// NOTE: feemarket module needs to be initialized before genutil module:
		// gentx transactions use MinGasPriceDecorator.AnteHandle
		feemarkettypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		icatypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		// ExoCore modules
		assetsTypes.ModuleName,
		depositTypes.ModuleName,
		operatorTypes.ModuleName,
		delegationTypes.ModuleName,
		withdrawTypes.ModuleName,
		rewardTypes.ModuleName,
		exoslashTypes.ModuleName,
		// no particular order required
		oracleTypes.ModuleName,
		// Evmos modules
		erc20types.ModuleName,
		epochstypes.ModuleName,
		recoverytypes.ModuleName,
		// NOTE: crisis module must go at the end to check for invariants on each module
		crisistypes.ModuleName,
		consensusparamtypes.ModuleName,
		avsManagerTypes.ModuleName,
		avsTaskTypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.configurator = module.NewConfigurator(
		app.appCodec,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
	)
	app.mm.RegisterServices(app.configurator)

	// add test gRPC service for testing gRPC queries in isolation
	// testdata.RegisterTestServiceServer(app.GRPCQueryRouter(), testdata.TestServiceImpl{})

	// create the simulation manager and define the order of the modules for deterministic
	// simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	overrideModules := map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(
			app.appCodec,
			app.AccountKeeper,
			authsims.RandomGenesisAccounts,
			app.GetSubspace(authtypes.ModuleName),
		),
	}
	app.sm = module.NewSimulationManagerFromAppModules(app.mm.Modules, overrideModules)

	autocliv1.RegisterQueryServer(
		app.GRPCQueryRouter(),
		runtimeservices.NewAutoCLIQueryService(app.mm.Modules),
	)

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	maxGasWanted := cast.ToUint64(appOpts.Get(srvflags.EVMMaxTxGasWanted))

	app.setAnteHandler(encodingConfig.TxConfig, maxGasWanted)
	app.setPostHandler()
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			logger.Error("error on loading last version", "err", err)
			os.Exit(1)
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	// Finally start the tpsCounter.
	app.tpsCounter = newTPSCounter(logger)
	go func() {
		// Unfortunately golangci-lint is so pedantic
		// so we have to ignore this error explicitly.
		_ = app.tpsCounter.start(context.Background())
	}()

	return app
}

// Name returns the name of the App
func (app *ExocoreApp) Name() string { return app.BaseApp.Name() }

func (app *ExocoreApp) setAnteHandler(txConfig client.TxConfig, maxGasWanted uint64) {
	options := ante.HandlerOptions{
		Cdc:                    app.appCodec,
		AccountKeeper:          app.AccountKeeper,
		BankKeeper:             app.BankKeeper,
		ExtensionOptionChecker: evmostypes.HasDynamicFeeExtensionOption,
		StakingKeeper:          app.StakingKeeper,
		FeegrantKeeper:         app.FeeGrantKeeper,
		DistributionKeeper:     app.RewardKeeper,
		IBCKeeper:              app.IBCKeeper,
		SignModeHandler:        txConfig.SignModeHandler(),
		SigGasConsumer:         ante.SigVerificationGasConsumer,
		MaxTxGasWanted:         maxGasWanted,
		FeeMarketKeeper:        app.FeeMarketKeeper,
		EvmKeeper:              app.EvmKeeper,
		TxFeeChecker:           ethante.NewDynamicFeeChecker(app.EvmKeeper),
	}

	if err := options.Validate(); err != nil {
		panic(err)
	}

	app.SetAnteHandler(ante.NewAnteHandler(options))
}

func (app *ExocoreApp) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(
		posthandler.HandlerOptions{},
	)
	if err != nil {
		panic(err)
	}

	app.SetPostHandler(postHandler)
}

// BeginBlocker runs the Tendermint ABCI BeginBlock logic. It executes state changes at the
// beginning of the new block for every registered module. If there is a registered fork at the
// current height,
// BeginBlocker will schedule the upgrade plan and perform the state migration (if any).
func (app *ExocoreApp) BeginBlocker(
	ctx sdk.Context,
	req abci.RequestBeginBlock,
) abci.ResponseBeginBlock {
	// Perform any scheduled forks before executing the modules logic
	app.ScheduleForkUpgrade(ctx)
	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker updates every end block
func (app *ExocoreApp) EndBlocker(
	ctx sdk.Context,
	req abci.RequestEndBlock,
) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// The DeliverTx method is intentionally decomposed to calculate the transactions per second.
func (app *ExocoreApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	defer func() {
		// TODO: Record the count along with the code and or reason so as to display
		// in the transactions per second live dashboards.
		if res.IsErr() {
			app.tpsCounter.incrementFailure()
		} else {
			app.tpsCounter.incrementSuccess()
		}
	}()
	return app.BaseApp.DeliverTx(req)
}

// InitChainer updates at chain initialization
func (app *ExocoreApp) InitChainer(
	ctx sdk.Context,
	req abci.RequestInitChain,
) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads state at a particular height
func (app *ExocoreApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *ExocoreApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)

	accs := make([]string, 0, len(maccPerms))
	for k := range maccPerms {
		accs = append(accs, k)
	}
	sort.Strings(accs)

	for _, acc := range accs {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns Evmos's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ExocoreApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns Evmos's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ExocoreApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Evmos's InterfaceRegistry
func (app *ExocoreApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ExocoreApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ExocoreApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *ExocoreApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ExocoreApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *ExocoreApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx

	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	node.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

func (app *ExocoreApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(
		app.BaseApp.GRPCQueryRouter(),
		clientCtx,
		app.BaseApp.Simulate,
		app.interfaceRegistry,
	)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *ExocoreApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// RegisterNodeService registers the node gRPC service on the provided
// application gRPC query router.
func (app *ExocoreApp) RegisterNodeService(clientCtx client.Context) {
	node.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}

// IBC Go TestingApp functions

// GetBaseApp implements the TestingApp interface.
func (app *ExocoreApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

// GetStakingKeeper implements the TestingApp interface.
func (app *ExocoreApp) GetStakingKeeper() ibctestingtypes.StakingKeeper {
	return app.StakingKeeper
}

// GetStakingKeeperSDK implements the TestingApp interface.
func (app *ExocoreApp) GetStakingKeeperSDK() stakingkeeper.Keeper {
	return app.StakingKeeper
}

// GetIBCKeeper implements the TestingApp interface.
func (app *ExocoreApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *ExocoreApp) GetTxConfig() client.TxConfig {
	cfg := encoding.MakeConfig(ModuleBasics)
	return cfg.TxConfig
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(_ client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(
	appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	// SDK subspaces
	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).
		// nolint: staticcheck
		WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	// ethermint subspaces
	paramsKeeper.Subspace(evmtypes.ModuleName).WithKeyTable(evmtypes.ParamKeyTable()) //nolint:staticcheck
	paramsKeeper.Subspace(oracleTypes.ModuleName).WithKeyTable(oracleTypes.ParamKeyTable())
	return paramsKeeper
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *ExocoreApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)

	accs := make([]string, 0, len(maccPerms))
	for k := range maccPerms {
		accs = append(accs, k)
	}
	sort.Strings(accs)

	for _, acc := range accs {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	for _, precompile := range common.DefaultPrecompilesBech32 {
		blockedAddrs[precompile] = true
	}

	return blockedAddrs
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *ExocoreApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}
