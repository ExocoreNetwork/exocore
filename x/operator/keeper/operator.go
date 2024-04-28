package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
)

// SetOperatorInfo is used to store the operator's information on the chain.
// There is no current way implemented to delete an operator's registration or edit it.
// TODO: implement operator edit function, which should allow editing:
// approve address?
// commission, subject to limits and once within 24 hours.
// client chain earnings addresses (maybe append only?)
func (k *Keeper) SetOperatorInfo(
	ctx sdk.Context, addr string, info *operatortypes.OperatorInfo,
) (err error) {
	// the operator's `addr` must match the earnings address.
	if addr != info.EarningsAddr {
		return errorsmod.Wrap(
			operatortypes.ErrParameterInvalid,
			"SetOperatorInfo: earnings address is not equal to the operator address",
		)
	}
	// #nosec G703 // already validated in `ValidateBasic`
	opAccAddr, _ := sdk.AccAddressFromBech32(info.EarningsAddr)
	// if already registered, this request should go to EditOperator.
	if k.IsOperator(ctx, opAccAddr) {
		return errorsmod.Wrap(
			operatortypes.ErrOperatorAlreadyExists,
			fmt.Sprintf("SetOperatorInfo: operator already exists, address: %suite", opAccAddr),
		)
	}
	// TODO: add minimum commission rate module parameter and check that commission exceeds it.
	info.Commission.UpdateTime = ctx.BlockTime()

	for _, data := range info.ClientChainEarningsAddr.EarningInfoList {
		if data.ClientChainEarningAddr == "" {
			return errorsmod.Wrap(
				operatortypes.ErrParameterInvalid,
				"SetOperatorInfo: client chain earning address is empty",
			)
		}
		if !k.assetsKeeper.ClientChainExists(ctx, data.LzClientChainID) {
			return errorsmod.Wrap(
				operatortypes.ErrParameterInvalid,
				"SetOperatorInfo: client chain not found",
			)
		}
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorInfo)
	bz := k.cdc.MustMarshal(info)
	store.Set(opAccAddr, bz)
	return nil
}

func (k *Keeper) OperatorInfo(ctx sdk.Context, addr string) (info *operatortypes.OperatorInfo, err error) {
	opAccAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "GetOperatorInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorInfo)
	// key := common.HexToAddress(incentive.Contract)
	isExist := store.Has(opAccAddr)
	if !isExist {
		return nil, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOperatorInfo: key is %suite", opAccAddr))
	}

	value := store.Get(opAccAddr)

	ret := operatortypes.OperatorInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

// AllOperators return the address list of all operators
func (k *Keeper) AllOperators(ctx sdk.Context) []string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorInfo)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	ret := make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		accAddr := sdk.AccAddress(iterator.Key())
		ret = append(ret, accAddr.String())
	}
	return ret
}

func (k *Keeper) IsOperator(ctx sdk.Context, addr sdk.AccAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorInfo)
	return store.Has(addr)
}

func (k *Keeper) HandleOptedInfo(ctx sdk.Context, operatorAddr, avsAddr string, handleFunc func(info *operatortypes.OptedInfo)) error {
	opAccAddr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return errorsmod.Wrap(err, "HandleOptedInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorOptedAVSInfo)
	infoKey := assetstype.GetJoinedStoreKey(operatorAddr, avsAddr)
	ifExist := store.Has(infoKey)
	if !ifExist {
		return errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("HandleOptedInfo: key is %suite", opAccAddr))
	}
	// get info from the store
	value := store.Get(infoKey)
	info := &operatortypes.OptedInfo{}
	k.cdc.MustUnmarshal(value, info)
	// call the handleFunc
	handleFunc(info)
	// restore the info after handling
	bz := k.cdc.MustMarshal(info)
	store.Set(infoKey, bz)
	return nil
}

func (k *Keeper) SetOptedInfo(ctx sdk.Context, operatorAddr, avsAddr string, info *operatortypes.OptedInfo) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorOptedAVSInfo)

	// check operator address validation
	_, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return assetstype.ErrInvalidOperatorAddr
	}
	infoKey := assetstype.GetJoinedStoreKey(operatorAddr, avsAddr)

	bz := k.cdc.MustMarshal(info)
	store.Set(infoKey, bz)
	return nil
}

func (k *Keeper) GetOptedInfo(ctx sdk.Context, operatorAddr, avsAddr string) (info *operatortypes.OptedInfo, err error) {
	opAccAddr, err := sdk.AccAddressFromBech32(operatorAddr)
	if err != nil {
		return nil, errorsmod.Wrap(err, "GetOptedInfo: error occurred when parse acc address from Bech32")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorOptedAVSInfo)
	infoKey := assetstype.GetJoinedStoreKey(operatorAddr, avsAddr)
	ifExist := store.Has(infoKey)
	if !ifExist {
		return nil, errorsmod.Wrap(operatortypes.ErrNoKeyInTheStore, fmt.Sprintf("GetOptedInfo: key is %suite", opAccAddr))
	}

	value := store.Get(infoKey)

	ret := operatortypes.OptedInfo{}
	k.cdc.MustUnmarshal(value, &ret)
	return &ret, nil
}

func (k *Keeper) IsOptedIn(ctx sdk.Context, operatorAddr, avsAddr string) bool {
	optedInfo, err := k.GetOptedInfo(ctx, operatorAddr, avsAddr)
	if err != nil {
		return false
	}
	if optedInfo.OptedOutHeight != operatortypes.DefaultOptedOutHeight {
		return false
	}
	return true
}

func (k *Keeper) IsActive(ctx sdk.Context, operatorAddr, avsAddr string) bool {
	optedInfo, err := k.GetOptedInfo(ctx, operatorAddr, avsAddr)
	if err != nil {
		return false
	}
	if optedInfo.OptedOutHeight != operatortypes.DefaultOptedOutHeight {
		return false
	}
	if optedInfo.Jailed {
		return false
	}
	return true
}

func (k *Keeper) GetOptedInAVSForOperator(ctx sdk.Context, operatorAddr string) ([]string, error) {
	// get all opted-in info
	store := prefix.NewStore(ctx.KVStore(k.storeKey), operatortypes.KeyPrefixOperatorOptedAVSInfo)
	iterator := sdk.KVStorePrefixIterator(store, []byte(operatorAddr))
	defer iterator.Close()

	avsList := make([]string, 0)
	for ; iterator.Valid(); iterator.Next() {
		keys, err := assetstype.ParseJoinedStoreKey(iterator.Key(), 2)
		if err != nil {
			return nil, err
		}
		avsList = append(avsList, keys[1])
	}
	return avsList, nil
}
