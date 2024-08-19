package keeper

import (
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/common"

	errorsmod "cosmossdk.io/errors"

	delegationtypes "github.com/ExocoreNetwork/exocore/x/delegation/types"
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
		operatorKeeper types.OperatorKeeper
		// other keepers
		assetsKeeper types.AssetsKeeper
		epochsKeeper types.EpochsKeeper
		evmKeeper    types.EVMKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	operatorKeeper types.OperatorKeeper,
	assetKeeper types.AssetsKeeper,
	epochsKeeper types.EpochsKeeper,
	evmKeeper types.EVMKeeper,
) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		operatorKeeper: operatorKeeper,
		assetsKeeper:   assetKeeper,
		epochsKeeper:   epochsKeeper,
		evmKeeper:      evmKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetOperatorKeeper returns the operatorKeeper from the Keeper struct.
func (k Keeper) GetOperatorKeeper() types.OperatorKeeper {
	return k.operatorKeeper
}

func (k Keeper) AVSInfoUpdate(ctx sdk.Context, params *types.AVSRegisterOrDeregisterParams) error {
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
		startingEpoch := uint64(epoch.CurrentEpoch + 1)
		if params.ChainID == types.ChainIDWithoutRevision(ctx.ChainID()) {
			// TODO: handle this better
			startingEpoch = uint64(epoch.CurrentEpoch)
		}
		// The caller must ensure that the relevant addresses are set to be the appropriate format (hex/bech32)
		avs := &types.AVSInfo{
			Name:                params.AvsName,
			AvsAddress:          params.AvsAddress,
			RewardAddr:          params.RewardContractAddr,
			SlashAddr:           params.SlashContractAddr,
			AvsOwnerAddress:     params.AvsOwnerAddress,
			AssetIDs:            params.AssetID,
			MinSelfDelegation:   params.MinSelfDelegation,
			AvsUnbondingPeriod:  params.UnbondingPeriod,
			EpochIdentifier:     epochIdentifier,
			StartingEpoch:       startingEpoch,
			MinOptInOperators:   params.MinOptInOperators,
			TaskAddr:            params.TaskAddr,
			MinStakeAmount:      params.MinStakeAmount, // Effective at CurrentEpoch+1, avoid immediate effects and ensure that the first epoch time of avs is equal to a normal identifier
			MinTotalStakeAmount: params.MinTotalStakeAmount,
			AvsSlash:            sdk.NewDecWithPrec(int64(params.AvsSlash), 2),
			AvsReward:           sdk.NewDecWithPrec(int64(params.AvsReward), 2),
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
		if epoch.CurrentEpoch-int64(avsInfo.GetInfo().StartingEpoch) > int64(avsInfo.Info.AvsUnbondingPeriod) {
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
		if int64(avsInfo.Info.AvsUnbondingPeriod) < (epoch.CurrentEpoch - int64(avsInfo.GetInfo().StartingEpoch)) {
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
		if params.MinStakeAmount > 0 {
			avs.MinStakeAmount = params.MinStakeAmount
		}
		if params.TaskAddr != "" {
			avs.TaskAddr = params.TaskAddr
		}
		if params.SlashContractAddr != "" {
			avs.SlashAddr = params.SlashContractAddr
		}
		if params.RewardContractAddr != "" {
			avs.RewardAddr = params.RewardContractAddr
		}
		if params.AvsOwnerAddress != nil {
			avs.AvsOwnerAddress = params.AvsOwnerAddress
		}
		if params.AssetID != nil {
			avs.AssetIDs = params.AssetID
		}

		if params.UnbondingPeriod > 0 {
			avs.AvsUnbondingPeriod = params.UnbondingPeriod
		}
		if params.MinSelfDelegation > 0 {
			avs.MinSelfDelegation = params.MinSelfDelegation
		}
		if params.EpochIdentifier != "" {
			avs.EpochIdentifier = params.EpochIdentifier
		}

		if params.MinOptInOperators > 0 {
			avs.MinOptInOperators = params.MinOptInOperators
		}
		if params.MinTotalStakeAmount > 0 {
			avs.MinTotalStakeAmount = params.MinTotalStakeAmount
		}
		if params.AvsSlash > 0 {
			avs.AvsSlash = sdk.NewDecWithPrec(int64(params.AvsSlash), 2)
		}
		if params.AvsReward > 0 {
			avs.AvsReward = sdk.NewDecWithPrec(int64(params.AvsReward), 2)
		}
		avs.AvsAddress = params.AvsAddress
		avs.StartingEpoch = uint64(epoch.CurrentEpoch + 1)

		return k.SetAVSInfo(ctx, avs)
	default:
		return errorsmod.Wrap(types.ErrInvalidAction, fmt.Sprintf("Invalid action: %d", action))
	}
}

func (k Keeper) CreateAVSTask(ctx sdk.Context, params *TaskParams) error {
	avsInfo := k.GetAVSInfoByTaskAddress(ctx, params.TaskContractAddress)
	if avsInfo.AvsAddress == "" {
		return errorsmod.Wrap(types.ErrUnregisterNonExistent, fmt.Sprintf("the taskaddr is :%s", params.TaskContractAddress))
	}
	// If avs CreateAVSTask check CallerAddress
	if !slices.Contains(avsInfo.AvsOwnerAddress, params.CallerAddress) {
		return errorsmod.Wrap(types.ErrCallerAddressUnauthorized, fmt.Sprintf("this caller not qualified to CreateAVSTask %s", params.CallerAddress))
	}

	epoch, found := k.epochsKeeper.GetEpochInfo(ctx, avsInfo.EpochIdentifier)
	if !found {
		return errorsmod.Wrap(types.ErrEpochNotFound, fmt.Sprintf("epoch info not found %s", avsInfo.EpochIdentifier))
	}

	if k.IsExistTask(ctx, params.TaskID, params.TaskContractAddress) {
		return errorsmod.Wrap(types.ErrAlreadyExists, fmt.Sprintf("the task is :%s", params.TaskID))
	}
	task := &types.TaskInfo{
		Name:                params.TaskName,
		Hash:                params.Hash,
		TaskContractAddress: params.TaskContractAddress,
		TaskId:              params.TaskID,
		TaskChallengePeriod: params.TaskChallengePeriod,
		ThresholdPercentage: params.ThresholdPercentage,
		TaskResponsePeriod:  params.TaskResponsePeriod,
		StartingEpoch:       uint64(epoch.CurrentEpoch + 1),
	}
	return k.SetTaskInfo(ctx, task)
}

func (k Keeper) RegisterBLSPublicKey(ctx sdk.Context, params *BlsParams) error {
	// TODO:check bls signature to prevent rogue key attacks
	// params.pubkeyRegistrationSignature == key.sig(params.pubkeyRegistrationMessageHash)
	if k.IsExistPubKey(ctx, params.Operator) {
		return errorsmod.Wrap(types.ErrAlreadyExists, fmt.Sprintf("the operator is :%s", params.Operator))
	}
	bls := &types.BlsPubKeyInfo{
		Name:     params.Name,
		Operator: params.Operator,
		PubKey:   params.PubKey,
	}
	return k.SetOperatorPubKey(ctx, bls)
}

func (k Keeper) GetOptInOperators(ctx sdk.Context, avsAddr string) ([]string, error) {
	return k.operatorKeeper.GetOptedInOperatorListByAVS(ctx, avsAddr)
}

func (k Keeper) OperatorOptAction(ctx sdk.Context, params *OperatorOptParams) error {
	operatorAddress := params.OperatorAddress
	opAccAddr, err := sdk.AccAddressFromBech32(operatorAddress)
	if err != nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when parse acc address from Bech32,the addr is:%s", operatorAddress))
	}

	if !k.operatorKeeper.IsOperator(ctx, opAccAddr) {
		return errorsmod.Wrap(delegationtypes.ErrOperatorNotExist, fmt.Sprintf("AVSInfoUpdate: invalid operator address:%s", operatorAddress))
	}

	f, err := k.IsAVS(ctx, params.AvsAddress)
	if err != nil {
		return errorsmod.Wrap(err, fmt.Sprintf("error occurred when get avs info,this avs address: %s", params.AvsAddress))
	}
	if !f {
		return fmt.Errorf("avs does not exist,this avs address: %s", params.AvsAddress)
	}

	switch params.Action {
	case RegisterAction:
		return k.operatorKeeper.OptIn(ctx, opAccAddr, params.AvsAddress)
	case DeRegisterAction:
		return k.operatorKeeper.OptOut(ctx, opAccAddr, params.AvsAddress)
	default:
		return errorsmod.Wrap(types.ErrInvalidAction, fmt.Sprintf("Invalid action: %d", params.Action))
	}
}

// SetAVSInfo sets the avs info. The caller must ensure that avs.AvsAddress is hex.
func (k Keeper) SetAVSInfo(ctx sdk.Context, avs *types.AVSInfo) (err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	bz := k.cdc.MustMarshal(avs)
	store.Set(common.HexToAddress(avs.AvsAddress).Bytes(), bz)
	return nil
}

func (k Keeper) GetAVSInfo(ctx sdk.Context, addr string) (*types.QueryAVSInfoResponse, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	value := store.Get(common.HexToAddress(addr).Bytes())
	if value == nil {
		return nil, errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("GetAVSInfo: key is %s", addr))
	}
	ret := types.AVSInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	res := &types.QueryAVSInfoResponse{
		Info: &ret,
	}
	return res, nil
}

func (k *Keeper) IsAVS(ctx sdk.Context, addr string) (bool, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	return store.Has(common.HexToAddress(addr).Bytes()), nil
}

// IsAVSByChainID queries whether an AVS exists by chainID.
// It returns the AVS address if it exists.
func (k Keeper) IsAVSByChainID(ctx sdk.Context, chainID string) (bool, common.Address) {
	avsAddr := types.GenerateAVSAddr(chainID)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	return store.Has(avsAddr.Bytes()), avsAddr
}

func (k Keeper) DeleteAVSInfo(ctx sdk.Context, addr string) error {
	hexAddr := common.HexToAddress(addr)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAVSInfo)
	if !store.Has(hexAddr.Bytes()) {
		return errorsmod.Wrap(types.ErrNoKeyInTheStore, fmt.Sprintf("AVSInfo didn't exist: key is %s", addr))
	}
	store.Delete(hexAddr[:])
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
