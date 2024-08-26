package keeper

import (
	"context"

	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	// other keepers
	assetsKeeper   delegationtype.AssetsKeeper
	slashKeeper    delegationtype.SlashKeeper
	operatorKeeper delegationtype.OperatorKeeper
	oracleKeeper   delegationtype.OracleKeeper
	hooks          delegationtype.DelegationHooks
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	assetsKeeper delegationtype.AssetsKeeper,
	slashKeeper delegationtype.SlashKeeper,
	operatorKeeper delegationtype.OperatorKeeper,
	oracleKeeper delegationtype.OracleKeeper,
) Keeper {
	return Keeper{
		storeKey:       storeKey,
		cdc:            cdc,
		assetsKeeper:   assetsKeeper,
		slashKeeper:    slashKeeper,
		operatorKeeper: operatorKeeper,
		oracleKeeper:   oracleKeeper,
	}
}

// SetHooks stores the given hooks implementations.
// Note that the Keeper is changed into a pointer to prevent an ineffective assignment.
func (k *Keeper) SetHooks(hooks delegationtype.DelegationHooks) {
	if hooks == nil {
		panic("cannot set nil hooks")
	}
	if k.hooks != nil {
		panic("cannot set hooks twice")
	}
	k.hooks = hooks
}

func (k *Keeper) Hooks() delegationtype.DelegationHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set to prevent calling nil functions
		return delegationtype.MultiDelegationHooks{}
	}
	return k.hooks
}

// IDelegation interface will be implemented by delegation keeper
type IDelegation interface {
	// DelegateAssetToOperator handle the DelegateAssetToOperator txs from msg service
	DelegateAssetToOperator(ctx context.Context, delegation *delegationtype.MsgDelegation) (*delegationtype.DelegationResponse, error)
	// UndelegateAssetFromOperator handle the UndelegateAssetFromOperator txs from msg service
	UndelegateAssetFromOperator(ctx context.Context, delegation *delegationtype.MsgUndelegation) (*delegationtype.UndelegationResponse, error)

	GetSingleDelegationInfo(ctx sdk.Context, stakerID, assetID, operatorAddr string) (*delegationtype.DelegationAmounts, error)

	GetDelegationInfo(ctx sdk.Context, stakerID, assetID string) (*delegationtype.QueryDelegationInfoResponse, error)
}
