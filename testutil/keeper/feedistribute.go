package keeper

import (
	"testing"

	stakingkeeper "github.com/ExocoreNetwork/exocore/x/dogfood/keeper"
	epochskeeper "github.com/ExocoreNetwork/exocore/x/epochs/keeper"
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
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtestutil "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	"github.com/golang/mock/gomock"
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
	distrAcc := authtypes.NewEmptyModuleAccount(types.ModuleName)
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	authority := authtypes.NewModuleAddress(types.ModuleName)
	ctrl := gomock.NewController(t)
	accountKeeper := distrtestutil.NewMockAccountKeeper(ctrl)
	accountKeeper.EXPECT().GetModuleAddress("distribution").Return(distrAcc.GetAddress())
	bankKeeper := distrtestutil.NewMockBankKeeper(ctrl)

	k := distrkeeper.NewKeeper(
		cdc,
		log.NewNopLogger(),
		"fee_collector",
		authority.String(),
		storeKey,
		bankKeeper,
		accountKeeper,
		stakingkeeper.Keeper{},
		epochskeeper.Keeper{},
	)

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}
