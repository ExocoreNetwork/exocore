package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueueOperation adds an operation to the consensus public key queue. If a similar operation
// already exists, the queue is not modified and QueueResultExists is returned. If a reverse
// operation already exists (removal + addition, or addition + removal), the old operation is
// dropped from the queue and QueueResultRemoved is returned. In the case that the operation is
// added to the queue, QueueResultSuccess is returned.
func (k Keeper) QueueOperation(
	ctx sdk.Context, addr sdk.AccAddress,
	key tmprotocrypto.PublicKey, operationType types.OperationType,
) types.QueueResultType {
	if operationType == types.KeyOpUnspecified {
		// should never happen
		panic("invalid operation type")
	}
	currentQueue := k.GetQueuedOperations(ctx)
	indexToDelete := len(currentQueue)
	for i, operation := range currentQueue {
		if operation.PubKey.Equal(key) {
			if operation.OperationType == operationType {
				return types.QueueResultExists
			}
			// reverse operation exists, remove it
			indexToDelete = i
			break
		}
	}
	ret := types.QueueResultSuccess
	if indexToDelete > len(currentQueue) {
		currentQueue = append(currentQueue[:indexToDelete], currentQueue[indexToDelete+1:]...)
		ret = types.QueueResultRemoved
	} else {
		operation := types.Operation{OperationType: operationType, OperatorAddress: addr, PubKey: key}
		currentQueue = append(currentQueue, operation)
	}
	operations := types.Operations{List: currentQueue}
	k.setQueuedOperations(ctx, operations)
	return ret
}
