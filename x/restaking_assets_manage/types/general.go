package types

const (
	CrossChainActionLength       = 1
	CrossChainOpAmountLength     = 32
	ClientChainLzIdIndexInTopics = 1

	ExoCoreOperatorAddrLength = 45
)

const (
	Deposit CrossChainOpType = iota
	WithdrawPrinciple
	WithDrawReward
	DelegationTo
	UnDelegationFrom
)

type CrossChainOpType uint8
