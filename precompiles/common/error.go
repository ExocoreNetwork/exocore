package common

const (
	ErrContractInputParaOrType = "the contract input parameter type or value error,arg index:%d, expected type is:%s,value:%v"
	ErrContractCaller          = "the caller doesn't have the permission to call this function"

	ErrInvalidAddrLength = "the length of input client chain or asset addr doesn't match,input:%d,expected:%v"

	ErrInputOperatorAddrLength = "mismatched length of the input operator address,actual is:%d,expected:%v"

	ErrInvalidInputList = "the length of input list is invalid, actual:%d, expected:%v"

	ErrInvalidMetaInfoLength = "nil meta info or too long for chain,input:%d,max:%d"
)
