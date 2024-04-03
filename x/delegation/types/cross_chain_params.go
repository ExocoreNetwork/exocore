package types

import (
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type DelegationOrUndelegationParams struct {
	ClientChainLzID uint64
	Action          assetstype.CrossChainOpType
	AssetsAddress   []byte
	OperatorAddress sdk.AccAddress
	StakerAddress   []byte
	OpAmount        sdkmath.Int
	LzNonce         uint64
	TxHash          common.Hash
	// todo: The operator approved signature might be needed here in future
}
