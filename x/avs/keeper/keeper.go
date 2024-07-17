package keeper

import (
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"

	assettypes "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
	//evmtypes "github.com/ExocoreNetwork/exocore/x/evm/keeper"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeKey       storetypes.StoreKey
		operatorKeeper delegationtypes.OperatorKeeper
		// other keepers
		assetsKeeper assettypes.Keeper
		epochsKeeper types.EpochsKeeper
		//evmKeeper    *evmtypes.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	operatorKeeper delegationtypes.OperatorKeeper,
	assetKeeper assettypes.Keeper,
	epochsKeeper types.EpochsKeeper,
	// evmKeeper *evmtypes.Keeper,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		operatorKeeper: operatorKeeper,
		assetsKeeper:   assetKeeper,
		epochsKeeper:   epochsKeeper,
		//evmKeeper:      evmKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) AVSInfoUpdate(ctx sdk.Context, params *AVSRegisterOrDeregisterParams) error {
	avsInfo, _ := k.GetAVSInfo(ctx, params.AvsAddress)
	action := params.Action
	epochIdentifier := params.EpochIdentifier
	if avsInfo != nil && avsInfo.Info.EpochIdentifier != "" {
		epochIdentifier = avsInfo.Info.EpochIdentifier
	}
	epoch, found := k.epochsKeeper.GetEpochInfo(ctx, epochIdentifier)
	if !found {
		return errorsmod.Wrap(types.ErrEpochNotFound, fmt.Sprintf("epoch info not found %s", epochIdentifier))
	}
	switch action {
	case RegisterAction:
		if avsInfo != nil {
			return errorsmod.Wrap(types.ErrAlreadyRegistered, fmt.Sprintf("the avsaddress is :%s", params.AvsAddress))
		}

		avs := &types.AVSInfo{
			Name:               params.AvsName,
			AvsAddress:         params.AvsAddress,
			RewardAddr:         params.RewardContractAddr,
			SlashAddr:          params.SlashContractAddr,
			AvsOwnerAddress:    params.AvsOwnerAddress,
			AssetId:            params.AssetID,
			MinSelfDelegation:  sdk.NewIntFromUint64(params.MinSelfDelegation),
			AvsUnbondingPeriod: uint32(params.UnbondingPeriod),
			EpochIdentifier:    epochIdentifier,
			OperatorAddress:    nil,
			StartingEpoch:      epoch.CurrentEpoch + 1, // Effective at CurrentEpoch+1, avoid immediate effects and ensure that the first epoch time of avs is equal to a normal identifier
		}
		return k.SetAVSInfo(ctx, avs)
	case DeRegisterAction:
		if avsInfo == nil {
			return errorsmod.Wrap(types.ErrUnregisterNonExistent, fmt.Sprintf("the avsaddress is :%s", params.AvsAddress))
		}
		// If avs DeRegisterAction check CallerAddress
		if !slices.Contains(avsInfo.Info.AvsOwnerAddress, params.CallerAddress) {
			return errorsmod.Wrap(types.ErrCallerAddressUnauthorized, fmt.Sprintf("this caller not qualified to deregister %s", params.CallerAddress))
		}

		// If avs DeRegisterAction check UnbondingPeriod
		if (epoch.CurrentEpoch - avsInfo.GetInfo().StartingEpoch) > int64(avsInfo.Info.AvsUnbondingPeriod) {
			return errorsmod.Wrap(types.ErrUnbondingPeriod, fmt.Sprintf("not qualified to deregister %s", avsInfo))
		}

		// If avs DeRegisterAction check avsname
		if avsInfo.Info.Name != params.AvsName {
			return errorsmod.Wrap(types.ErrAvsNameMismatch, fmt.Sprintf("Unregistered AVS name is incorrect %s", params.AvsName))
		}
		return k.DeleteAVSInfo(ctx, params.AvsAddress)
	case UpdateAction:
		if avsInfo == nil {
			return errorsmod.Wrap(types.ErrUnregisterNonExistent, fmt.Sprintf("the avsaddress is :%s", params.AvsAddress))
		}
		// If avs UpdateAction check UnbondingPeriod
		if int64(avsInfo.Info.AvsUnbondingPeriod) < (epoch.CurrentEpoch - avsInfo.GetInfo().StartingEpoch) {
			return errorsmod.Wrap(types.ErrUnbondingPeriod, fmt.Sprintf("not qualified to deregister %s", avsInfo))
		}
		// If avs UpdateAction check CallerAddress
		if !slices.Contains(avsInfo.Info.AvsOwnerAddress, params.CallerAddress) {
			return errorsmod.Wrap(types.ErrCallerAddressUnauthorized, fmt.Sprintf("this caller not qualified to update %s", params.CallerAddress))
		}
		avs := avsInfo.Info
		if params.AvsName != "" {
			avs.Name = params.AvsName
		}
		if params.MinSelfDelegation > 0 {
			avs.MinSelfDelegation = sdk.NewIntFromUint64(params.MinSelfDelegation)
		}
		if params.AvsOwnerAddress != nil {
			avs.AvsOwnerAddress = params.AvsOwnerAddress
		}
		if params.SlashContractAddr != "" {
			avs.SlashAddr = params.SlashContractAddr
		}
		if params.EpochIdentifier != "" {
			avs.EpochIdentifier = params.EpochIdentifier
		}
		if params.UnbondingPeriod > 0 {
			avs.AvsUnbondingPeriod = uint32(params.UnbondingPeriod)
		}
		if params.AssetID != nil {
			avs.AssetId = params.AssetID
		}
		avs.AvsAddress = params.AvsAddress
		avs.StartingEpoch = epoch.CurrentEpoch + 1

		return k.SetAVSInfo(ctx, avs)
	default:
		return errorsmod.Wrap(types.ErrInvalidAction, fmt.Sprintf("Invalid action: %d", action))
	}
}

func (k Keeper) AVSInfoUpdateWithOperator(ctx sdk.Context, params *OperatorOptParams) error {
	operatorAddress := params.OperatorAddress
	opAccAddr, err := sdk.AccAddressFromBech32(operatorAddress)
	if err != nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when parse acc address from Bech32,the addr is:%s", operatorAddress))
	}

	if !k.operatorKeeper.IsOperator(ctx, opAccAddr) {
		return errorsmod.Wrap(delegationtypes.ErrOperatorNotExist, fmt.Sprintf("AVSInfoUpdate: invalid operator address:%s", operatorAddress))
	}

	avsInfo, err := k.GetAVSInfo(ctx, params.AvsAddress)
	if err != nil || avsInfo == nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when get avs info,this avs address: %s", params.AvsAddress))
	}

	avs := avsInfo.GetInfo()
	operatorAddrList := avs.OperatorAddress

	switch params.Action {
	case RegisterAction:
		if slices.Contains(operatorAddrList, operatorAddress) {
			return errorsmod.Wrap(types.ErrAlreadyRegistered, fmt.Sprintf("Error: Already registered, operatorAddress %s", operatorAddress))
		}
		operatorAddrList = append(operatorAddrList, operatorAddress)
		avs.OperatorAddress = operatorAddrList
		return k.SetAVSInfo(ctx, avs)
	case DeRegisterAction:
		if !slices.Contains(operatorAddrList, operatorAddress) {
			return errorsmod.Wrap(types.ErrUnregisterNonExistent, fmt.Sprintf("No available operatorAddress to DeRegisterAction, operatorAddress: %s", operatorAddress))
		}
		epoch, found := k.epochsKeeper.GetEpochInfo(ctx, avs.EpochIdentifier)
		if !found {
			return errorsmod.Wrap(types.ErrEpochNotFound, fmt.Sprintf("epoch info not found %s", avs.EpochIdentifier))
		}
		// If avs UpdateAction check UnbondingPeriod
		if int64(avsInfo.Info.AvsUnbondingPeriod) < (epoch.CurrentEpoch - avsInfo.GetInfo().StartingEpoch) {
			return errorsmod.Wrap(types.ErrUnbondingPeriod, fmt.Sprintf("not qualified to deregister %s", operatorAddress))
		}
		avs.OperatorAddress = types.RemoveOperatorAddress(operatorAddrList, operatorAddress)
		return k.SetAVSInfo(ctx, avs)
	default:
		return errorsmod.Wrap(types.ErrInvalidAction, fmt.Sprintf("Invalid action: %d", params.Action))
	}
}

func (k Keeper) SetAVSInfo(ctx sdk.Context, avs *types.AVSInfo) (err error) {
	avsAddr, err := sdk.AccAddressFromBech32(avs.AvsAddress)
	if err != nil {
		return errorsmod.Wrap(err, "SetAVSInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)

	bz := k.cdc.MustMarshal(avs)
	store.Set(avsAddr, bz)
	return nil
}

func (k Keeper) GetAVSInfo(ctx sdk.Context, addr string) (*types.QueryAVSInfoResponse, error) {
	avsAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "GetAVSInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	if !store.Has(avsAddr) {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSInfo: key is %s", avsAddr))
	}

	value := store.Get(avsAddr)
	ret := types.AVSInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	res := &types.QueryAVSInfoResponse{
		Info: &ret,
	}
	return res, nil
}

func (k Keeper) DeleteAVSInfo(ctx sdk.Context, addr string) error {
	avsAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return errorsmod.Wrap(err, "AVSInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	if !store.Has(avsAddr) {
		return errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("AVSInfo didn't exist: key is %s", avsAddr))
	}
	store.Delete(avsAddr)
	return nil
}

// IterateAVSInfo iterate through avs
func (k Keeper) IterateAVSInfo(ctx sdk.Context, fn func(index int64, avsInfo types.AVSInfo) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		avs := types.AVSInfo{}
		k.cdc.MustUnmarshal(iterator.Value(), &avs)

		stop := fn(i, avs)

		if stop {
			break
		}
		i++
	}
}
