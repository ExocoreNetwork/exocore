package delegation

const (
	ErrContractInputParaOrType = "the contract input parameter type or value error,arg index:%d, type is:%s,value:%v"
	ErrContractCaller          = "the caller doesn't have the permission to call this function,caller:%s,need:%s"
	ErrCtxTxHash               = "ctx TxHash type error or is nil,type is:%v,value:%v"

	ErrInputOperatorAddrLength = "mismatched length of the input operator address,actual is:%d,expect:%v"
)
