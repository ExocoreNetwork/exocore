package keeper

import (
	"context"

	delegationtype "github.com/ExocoreNetwork/exocore/x/delegation/types"
	depositkeeper "github.com/ExocoreNetwork/exocore/x/deposit/keeper"
	"github.com/ExocoreNetwork/exocore/x/restaking_assets_manage/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	//other keepers
	restakingStateKeeper      keeper.Keeper
	depositKeeper             depositkeeper.Keeper
	slashKeeper               delegationtype.ISlashKeeper
	expectedOperatorInterface delegationtype.ExpectedOperatorInterface
}

func NewKeeper(
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	restakingStateKeeper keeper.Keeper,
	depositKeeper depositkeeper.Keeper,
	slashKeeper delegationtype.ISlashKeeper,
	operatorKeeper delegationtype.ExpectedOperatorInterface,
) Keeper {
	return Keeper{
		storeKey:                  storeKey,
		cdc:                       cdc,
		restakingStateKeeper:      restakingStateKeeper,
		depositKeeper:             depositKeeper,
		slashKeeper:               slashKeeper,
		expectedOperatorInterface: operatorKeeper,
	}
}

// GetExoCoreLzAppAddress Get exoCoreLzAppAddr from deposit keeper,it will be used when check the caller of precompile contract.
// This function needs to be moved to `restaking_assets_manage` module,which will facilitate its use for the other modules
func (k Keeper) GetExoCoreLzAppAddress(ctx sdk.Context) (common.Address, error) {
	return k.depositKeeper.GetExoCoreLzAppAddress(ctx)
}

// IDelegation interface will be implemented by deposit keeper
type IDelegation interface {
	// PostTxProcessing automatically call PostTxProcessing to update delegation state after receiving delegation event tx from layerZero protocol
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error

	// DelegateAssetToOperator handle the DelegateAssetToOperator txs from msg service
	DelegateAssetToOperator(ctx context.Context, delegation *delegationtype.MsgDelegation) (*delegationtype.DelegationResponse, error)
	// UndelegateAssetFromOperator handle the UndelegateAssetFromOperator txs from msg service
	UndelegateAssetFromOperator(ctx context.Context, delegation *delegationtype.MsgUndelegation) (*delegationtype.UndelegationResponse, error)

	GetSingleDelegationInfo(ctx sdk.Context, stakerId, assetId, operatorAddr string) (*delegationtype.DelegationAmounts, error)

	GetDelegationInfo(ctx sdk.Context, stakerId, assetId string) (*delegationtype.QueryDelegationInfoResponse, error)
}
