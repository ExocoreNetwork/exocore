package keeper

import (
	context "context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/delegation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/minio/sha256-simd"
)

var _ types.MsgServer = &Keeper{}

// DelegateAssetToOperator todo: Delegation and Undelegation from exoCore chain directly will be implemented in future.At the moment,they are executed from client chain
func (k *Keeper) DelegateAssetToOperator(goCtx context.Context, msg *types.MsgDelegation) (*types.DelegationResponse, error) {
	// TODO: currently we only support delegation with native token by invoking service
	ctx := sdk.UnwrapSDKContext(goCtx)
	fromAddr := sdk.MustAccAddressFromBech32(msg.BaseInfo.FromAddress)
	nonce, err := k.accountKeeper.GetSequence(ctx, fromAddr)
	if err != nil {
		return nil, err
	}
	stakerAddr := sdk.MustAccAddressFromBech32(msg.BaseInfo.FromAddress).Bytes()
	txBytes := ctx.TxBytes()
	txHash := sha256.Sum256(txBytes)
	combined := fmt.Sprintf("%s-%d", txHash, nonce)
	uniqueHash := sha256.Sum256([]byte(combined))
	for operatorAddrStr, value := range msg.BaseInfo.PerOperatorAmounts {
		operatorAddr := sdk.MustAccAddressFromBech32(operatorAddrStr)
		inputParams := types.NewDelegationOrUndelegationParams(
			assetstypes.NativeChainLzID,
			assetstypes.DelegateTo,
			[]byte(assetstypes.NativeAssetAddr),
			operatorAddr,
			stakerAddr,
			value.Amount,
			nonce,
			uniqueHash,
		)
		ctx := sdk.UnwrapSDKContext(goCtx)
		err := k.DelegateTo(ctx, inputParams)
		return &types.DelegationResponse{}, err
	}

	return nil, errorsmod.Wrap(types.ErrNotSupportYet, "func:DelegateAssetToOperator")
}

func (k *Keeper) UndelegateAssetFromOperator(_ context.Context, _ *types.MsgUndelegation) (*types.UndelegationResponse, error) {
	return nil, errorsmod.Wrap(types.ErrNotSupportYet, "func:UndelegateAssetFromOperator")
}
