package types

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strings"
)

const (
	CrossChainActionLength       = 1
	CrossChainOpAmountLength     = 32
	GeneralAssetsAddrLength      = 32
	GeneralClientChainAddrLength = 32

	ClientChainLzIdIndexInTopics = 0
	LzNonceIndexInTopics         = 2

	ExoCoreOperatorAddrLength = 45
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

func GetStakeIDAndAssetId(clientChainLzId uint64, stakerAddress []byte, assetsAddress []byte) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(clientChainLzId)
	if stakerAddress != nil {
		stakeId = strings.Join([]string{hexutil.Encode(stakerAddress), clientChainLzIdStr}, "_")
	}

	if assetsAddress != nil {
		assetId = strings.Join([]string{hexutil.Encode(assetsAddress), clientChainLzIdStr}, "_")
	}
	return
}

func GetStakeIDAndAssetIdFromStr(clientChainLzId uint64, stakerAddress string, assetsAddress string) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(clientChainLzId)
	if stakerAddress != "" {
		stakeId = strings.Join([]string{strings.ToLower(stakerAddress), clientChainLzIdStr}, "_")
	}

	if assetsAddress != "" {
		assetId = strings.Join([]string{strings.ToLower(assetsAddress), clientChainLzIdStr}, "_")
	}
	return
}
