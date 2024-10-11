package keeper

import (
	"bytes"
	"fmt"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingkeeper "github.com/ExocoreNetwork/exocore/x/dogfood/keeper"
	"github.com/ExocoreNetwork/exocore/x/feedistribution/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
)

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		logger   log.Logger
		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority    string
		authKeeper   types.AccountKeeper
		bankKeeper   types.BankKeeper
		epochsKeeper types.EpochsKeeper

		feeCollectorName string

		StakingKeeper stakingkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	logger log.Logger,
	feeCollectorName, authority string,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	stakingkeeper stakingkeeper.Keeper,
	epochKeeper types.EpochsKeeper,
) Keeper {
	// ensure distribution module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	k := &Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		logger:           logger,
		authority:        authority,
		authKeeper:       accountKeeper,
		bankKeeper:       bankKeeper,
		epochsKeeper:     epochKeeper,
		feeCollectorName: feeCollectorName,
		StakingKeeper:    stakingkeeper,
	}

	return *k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// set the global fee pool distribution info
func (k Keeper) SetFeePool(ctx sdk.Context, feePool *types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(feePool)
	store.Set(types.FeePoolKey, b)
}

// get the global fee pool distribution info
func (k Keeper) GetFeePool(ctx sdk.Context) (feePool *types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.FeePoolKey)
	if b == nil {
		feePool := &types.FeePool{}
		store := ctx.KVStore(k.storeKey)
		b := k.cdc.MustMarshal(feePool)
		store.Set(types.FeePoolKey, b)
		return feePool
	}
	fp := &types.FeePool{}
	k.cdc.MustUnmarshal(b, fp)
	return fp
}

// get accumulated commission for a validator
func (k Keeper) GetValidatorAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress) (commission types.ValidatorAccumulatedCommission) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorAccumulatedCommissionKey(val))
	if b == nil {
		return types.ValidatorAccumulatedCommission{}
	}
	k.cdc.MustUnmarshal(b, &commission)
	return
}

// GetAllValidatorData returns a slice containing all accumulated commissions for validators.
func (k Keeper) GetAllValidatorData(ctx sdk.Context) (map[string]interface{}, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	commissions := make([]types.ValidatorAccumulatedCommissions, 0)
	currentList := make([]types.ValidatorCurrentRewardsList, 0)
	outList := make([]types.ValidatorOutstandingRewardsList, 0)
	stakerList := make([]types.StakerOutstandingRewardsList, 0)
	validatorData := map[string]interface{}{
		"ValidatorAccumulatedCommissions": commissions,
		"ValidatorCurrentRewardsList":     currentList,
		"ValidatorOutstandingRewardsList": outList,
		"StakerOutstandingRewardsList":    stakerList,
	}

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()

		switch {
		case bytes.HasPrefix(key, types.GetValidatorAccumulatedCommissionKey(sdk.ValAddress{})):
			var commission types.ValidatorAccumulatedCommission
			err := k.cdc.Unmarshal(value, &commission)
			if err != nil {
				return nil, err
			}
			commissions = append(commissions, types.ValidatorAccumulatedCommissions{
				ValAddr:    string(key),
				Commission: &commission,
			})
		case bytes.HasPrefix(key, types.GetValidatorCurrentRewardsKey(sdk.ValAddress{})):
			var rewards types.ValidatorCurrentRewards
			err := k.cdc.Unmarshal(value, &rewards)
			if err != nil {
				return nil, err
			}
			currentList = append(currentList, types.ValidatorCurrentRewardsList{
				ValAddr:        string(key),
				CurrentRewards: &rewards,
			})

		case bytes.HasPrefix(key, types.GetValidatorOutstandingRewardsKey(sdk.ValAddress{})):
			var outstandingRewards types.ValidatorOutstandingRewards
			err := k.cdc.Unmarshal(value, &outstandingRewards)
			if err != nil {
				return nil, err
			}
			outList = append(outList, types.ValidatorOutstandingRewardsList{
				ValAddr:            string(key),
				OutstandingRewards: &outstandingRewards,
			})

		case bytes.HasPrefix(key, types.GetStakerOutstandingRewardsKey("")):
			var stakerRewards types.StakerOutstandingRewards
			err := k.cdc.Unmarshal(value, &stakerRewards)
			if err != nil {
				return nil, err
			}
			stakerList = append(stakerList, types.StakerOutstandingRewardsList{
				ValAddr:                  string(key),
				StakerOutstandingRewards: &stakerRewards,
			})

		default:
			continue
		}
	}

	return validatorData, nil
}

// set accumulated commission for a validator
func (k Keeper) SetValidatorAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress, commission types.ValidatorAccumulatedCommission) {
	var bz []byte

	store := ctx.KVStore(k.storeKey)
	if commission.Commission.IsZero() {
		bz = k.cdc.MustMarshal(&types.ValidatorAccumulatedCommission{})
	} else {
		bz = k.cdc.MustMarshal(&commission)
	}

	store.Set(types.GetValidatorAccumulatedCommissionKey(val), bz)
}

// get current rewards for a validator
func (k Keeper) GetValidatorCurrentRewards(ctx sdk.Context, val sdk.ValAddress) (rewards types.ValidatorCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorCurrentRewardsKey(val))
	k.cdc.MustUnmarshal(b, &rewards)
	return
}

// set current rewards for a validator
func (k Keeper) SetValidatorCurrentRewards(ctx sdk.Context, val sdk.ValAddress, rewards types.ValidatorCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetValidatorCurrentRewardsKey(val), b)
}

// get validator outstanding rewards
func (k Keeper) GetValidatorOutstandingRewards(ctx sdk.Context, val sdk.ValAddress) (rewards types.ValidatorOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorOutstandingRewardsKey(val))
	k.cdc.MustUnmarshal(bz, &rewards)
	return
}

// set validator outstanding rewards
func (k Keeper) SetValidatorOutstandingRewards(ctx sdk.Context, val sdk.ValAddress, rewards types.ValidatorOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetValidatorOutstandingRewardsKey(val), b)
}

// set the reward to delegator
func (k Keeper) SetStakerRewards(ctx sdk.Context, stakerAddress string, rewards types.StakerOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&rewards)
	store.Set(types.GetStakerOutstandingRewardsKey(stakerAddress), b)
}

// get the reward of delegator
func (k Keeper) GetStakerRewards(ctx sdk.Context, stakerAddress string) (rewards types.StakerOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetStakerOutstandingRewardsKey(stakerAddress))
	k.cdc.MustUnmarshal(bz, &rewards)
	return
}
