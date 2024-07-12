package keeper

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (*types.Params, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixParams)
	ifExist := store.Has(types.ParamsKey)
	if !ifExist {
		return nil, fmt.Errorf("params %s not found", types.KeyPrefixParams)
	}

	value := store.Get(types.ParamsKey)

	ret := &types.Params{}
	k.cdc.MustUnmarshal(value, ret)
	return ret, nil
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixParams)
	bz := k.cdc.MustMarshal(params)
	store.Set(types.ParamsKey, bz)
	return nil
}

type AVSRegisterOrDeregisterParams struct {
	AvsName         string
	AvsAddress      string
	Action          uint64
	AvsOwnerAddress []string
	AssetID         []string

	MinSelfDelegation  uint64
	UnbondingPeriod    uint64
	RewardContractAddr string
	SlashContractAddr  string
	OperatorAddress    []string
	EpochIdentifier    string
	CallerAddress      string
}
type OperatorOptParams struct {
	Name            string
	BlsPublicKey    string
	IsRegistered    bool
	Action          uint64
	OperatorAddress string
	Status          string
	AvsAddress      string
}

const (
	RegisterAction   = 1
	DeRegisterAction = 2
	UpdateAction     = 3
)
