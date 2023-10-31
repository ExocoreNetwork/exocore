package types

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
