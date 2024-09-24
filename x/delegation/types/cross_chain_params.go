package types

import (
	sdkmath "cosmossdk.io/math"
	assetstype "github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

type DelegationOrUndelegationParams struct {
	ClientChainID   uint64
	Action          assetstype.CrossChainOpType
	AssetsAddress   []byte
	OperatorAddress sdk.AccAddress
	StakerAddress   []byte
	OpAmount        sdkmath.Int
	LzNonce         uint64
	TxHash          common.Hash
	// todo: The operator approved signature might be needed here in future
}

func NewDelegationOrUndelegationParams(
	clientChainID uint64,
	action assetstype.CrossChainOpType,
	assetsAddress []byte,
	operatorAddress sdk.AccAddress,
	stakerAddress []byte,
	opAmount sdkmath.Int,
	lzNonce uint64,
	txHash common.Hash,
) *DelegationOrUndelegationParams {
	return &DelegationOrUndelegationParams{
		ClientChainID:   clientChainID,
		Action:          action,
		AssetsAddress:   assetsAddress,
		OperatorAddress: operatorAddress,
		StakerAddress:   stakerAddress,
		OpAmount:        opAmount,
		LzNonce:         lzNonce,
		TxHash:          txHash,
	}
}
