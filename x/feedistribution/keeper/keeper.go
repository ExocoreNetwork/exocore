package keeper

import (
	collections "cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"fmt"
	stakingkeeper "github.com/ExocoreNetwork/exocore/x/dogfood/keeper"
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		storeKey     storetypes.StoreKey
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority    string
		authKeeper   types.AccountKeeper
		bankKeeper   types.BankKeeper
		epochsKeeper types.EpochsKeeper
		poolKeeper   types.PoolKeeper

		feeCollectorName string
		// FeePool stores decimal tokens that cannot be yet distributed.
		FeePool       collections.Item[types.FeePool]
		StakingKeeper stakingkeeper.Keeper
		// ValidatorsAccumulatedCommission key: valAddr | value: ValidatorAccumulatedCommission
		ValidatorsAccumulatedCommission collections.Map[sdk.ValAddress, types.ValidatorAccumulatedCommission]
		// ValidatorCurrentRewards key: valAddr | value: ValidatorCurrentRewards
		ValidatorCurrentRewards collections.Map[sdk.ValAddress, types.ValidatorCurrentRewards]
		// ValidatorOutstandingRewards key: valAddr | value: ValidatorOustandingRewards
		ValidatorOutstandingRewards collections.Map[sdk.ValAddress, types.ValidatorOutstandingRewards]
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		storeKey:     storeKey,
		logger:       logger,

		bankKeeper: bankKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
