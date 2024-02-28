package types

import (
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	CrossChainActionLength       = 1
	CrossChainOpAmountLength     = 32
	GeneralAssetsAddrLength      = 32
	GeneralClientChainAddrLength = 32

	ClientChainLzIDIndexInTopics = 0
	LzNonceIndexInTopics         = 2

	ExoCoreOperatorAddrLength = 42
)

type GeneralAssetsAddr [32]byte

type GeneralClientChainAddr [32]byte

type CrossChainOpType uint8

type WithdrawerAddress [32]byte

const (
	Deposit CrossChainOpType = iota
	WithdrawPrinciple
	WithDrawReward
	DelegateTo
	UndelegateFrom
	Slash
)

// GetStakeIDAndAssetID stakerID = stakerAddress+'_'+clientChainLzID,assetID = assetAddress+'_'+clientChainLzID
func GetStakeIDAndAssetID(clientChainLzID uint64, stakerAddress []byte, assetsAddress []byte) (stakeID string, assetID string) {
	clientChainLzIDStr := hexutil.EncodeUint64(clientChainLzID)
	if stakerAddress != nil {
		stakeID = strings.Join([]string{hexutil.Encode(stakerAddress), clientChainLzIDStr}, "_")
	}

	if assetsAddress != nil {
		assetID = strings.Join([]string{hexutil.Encode(assetsAddress), clientChainLzIDStr}, "_")
	}
	return
}

// GetStakeIDAndAssetIDFromStr stakerID = stakerAddress+'_'+clientChainLzID,assetID = assetAddress+'_'+clientChainLzID
func GetStakeIDAndAssetIDFromStr(clientChainLzID uint64, stakerAddress string, assetsAddress string) (stakeID string, assetID string) {
	clientChainLzIDStr := hexutil.EncodeUint64(clientChainLzID)
	if stakerAddress != "" {
		stakeID = strings.Join([]string{strings.ToLower(stakerAddress), clientChainLzIDStr}, "_")
	}

	if assetsAddress != "" {
		assetID = strings.Join([]string{strings.ToLower(assetsAddress), clientChainLzIDStr}, "_")
	}
	return
}
