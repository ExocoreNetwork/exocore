package keeper

import (
	"testing"

	distrkeeper "github.com/ExocoreNetwork/exocore/x/feedistribution/keeper"
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/stretchr/testify/require"
)

func FeedistributeKeeper(t testing.TB) (distrkeeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	encCfg := moduletestutil.MakeTestEncodingConfig(auth.AppModule{}, bank.AppModule{}, distribution.AppModule{})
	key := sdk.NewKVStoreKey(banktypes.StoreKey)
	blockedAddresses := map[string]bool{
		accountkeeper.AccountKeeper{}.GetAuthority(): false,
	}
	authority := authtypes.NewModuleAddress(types.ModuleName)
	bankkeeper := bankkeeper.NewBaseKeeper(
		encCfg.Codec, key, accountkeeper.AccountKeeper{},
		blockedAddresses, authority.String(),
	)
	k := distrkeeper.NewKeeper(
		cdc,
		log.NewNopLogger(),
		authority.String(),
		storeKey,
		bankkeeper,
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
