package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/exocore/x/deposit/types"
	"strings"
)

var ParamsKey = []byte("Params")

func (k Keeper) SetParams(ctx sdk.Context, params *types2.Params) error {
	//check if addr is evm address
	if !common.IsHexAddress(params.ExoCoreLzAppAddress) {
		return types2.ErrInvalidEvmAddressFormat
	}
	if len(common.FromHex(params.ExoCoreLzAppEventTopic)) != common.HashLength {
		return types2.ErrInvalidLzUaTopicIdLength
	}
	params.ExoCoreLzAppAddress = strings.ToLower(params.ExoCoreLzAppAddress)
	params.ExoCoreLzAppEventTopic = strings.ToLower(params.ExoCoreLzAppEventTopic)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixParams)
	//key := common.HexToAddress(incentive.Contract)
	bz := k.cdc.MustMarshal(params)
	store.Set(ParamsKey, bz)
	return nil
}

func (k Keeper) GetParams(ctx sdk.Context) (*types2.Params, error) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types2.KeyPrefixParams)
	ifExist := store.Has(ParamsKey)
	if !ifExist {
		return nil, types2.ErrNoParamsKey
	}

	value := store.Get(ParamsKey)

	ret := &types2.Params{}
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
