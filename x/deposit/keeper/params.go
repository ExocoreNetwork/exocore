package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	deposittype "github.com/exocore/x/deposit/types"
)

var ParamsKey = []byte("Params")

func (k Keeper) SetParams(ctx sdk.Context, params *deposittype.Params) error {
	// check if addr is evm address
	if !common.IsHexAddress(params.ExoCoreLzAppAddress) {
		return deposittype.ErrInvalidEvmAddressFormat
	}
	if len(common.FromHex(params.ExoCoreLzAppEventTopic)) != common.HashLength {
		return deposittype.ErrInvalidLzUaTopicIdLength
	}
	params.ExoCoreLzAppAddress = strings.ToLower(params.ExoCoreLzAppAddress)
	params.ExoCoreLzAppEventTopic = strings.ToLower(params.ExoCoreLzAppEventTopic)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), deposittype.KeyPrefixParams)
	bz := k.cdc.MustMarshal(params)
	store.Set(ParamsKey, bz)
	return nil
}

func (k Keeper) GetParams(ctx sdk.Context) (*deposittype.Params, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), deposittype.KeyPrefixParams)
	isExist := store.Has(ParamsKey)
	if !isExist {
		return nil, deposittype.ErrNoParamsKey
	}

	value := store.Get(ParamsKey)

	ret := &deposittype.Params{}
	k.cdc.MustUnmarshal(value, ret)
	return ret, nil
}

func (k Keeper) GetExoCoreLzAppAddress(ctx sdk.Context) (common.Address, error) {
	depositModuleParam, err := k.GetParams(ctx)
	if err != nil {
		return common.Address{}, err
	}
	return common.HexToAddress(depositModuleParam.ExoCoreLzAppAddress), nil
}
