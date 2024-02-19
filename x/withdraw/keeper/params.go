package keeper

import (
	"strings"

	paramstypes "github.com/ExocoreNetwork/exocore/x/deposit/types"
	types "github.com/ExocoreNetwork/exocore/x/withdraw/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (*paramstypes.Params, error) {
	// store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixParams)
	// ifExist := store.Has(types.ParamsKey)
	// if !ifExist {
	// 	return nil, types.ErrNoParamsKey
	// }

	// value := store.Get(types.ParamsKey)

	// ret := &types.Params{}
	// k.cdc.MustUnmarshal(value, ret)
	// return ret, nil
	// Uify the way to obtain Params from deposit keeper
	return k.depositKeeper.GetParams(ctx)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	// check if addr is evm address
	if !common.IsHexAddress(params.ExoCoreLzAppAddress) {
		return types.ErrInvalidEvmAddressFormat
	}
	if len(common.FromHex(params.ExoCoreLzAppEventTopic)) != common.HashLength {
		return types.ErrInvalidLzUaTopicIDLength
	}
	params.ExoCoreLzAppAddress = strings.ToLower(params.ExoCoreLzAppAddress)
	params.ExoCoreLzAppEventTopic = strings.ToLower(params.ExoCoreLzAppEventTopic)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixParams)
	bz := k.cdc.MustMarshal(params)
	store.Set(types.ParamsKey, bz)
	return nil
}
