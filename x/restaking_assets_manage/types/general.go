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

	ClientChainLzIdIndexInTopics = 1

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
	DelegationTo
	UnDelegationFrom
)

func GetStakeIDAndAssetId(clientChainLzId uint64, stakerAddress []byte, assetsAddress []byte) (stakeId string, assetId string) {
	clientChainLzIdStr := hexutil.EncodeUint64(clientChainLzId)
	stakeId = strings.Join([]string{hexutil.Encode(stakerAddress), clientChainLzIdStr}, "_")
	assetId = strings.Join([]string{hexutil.Encode(assetsAddress), clientChainLzIdStr}, "_")
	return
}
