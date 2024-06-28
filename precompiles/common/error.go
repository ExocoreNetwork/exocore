package common

const (
	ErrContractInputParaOrType = "the contract input parameter type or value error,arg index:%d, expected type is:%s,value:%v"
	ErrContractCaller          = "the caller doesn't have the permission to call this function"

	ErrInvalidAddrLength = "invalid length of staker or asset addr, actualLength:%d,min:%d"

	ErrInputOperatorAddrLength = "mismatched length of the input operator address,actual is:%d,expected:%v"

	ErrInvalidInputList = "the length of input list is invalid, field:%s, actualLength:%d, expected:%v"

	ErrInvalidMetaInfoLength = "nil meta info or too long for chain or token,value:%s,actualLength:%d,max:%d"

	ErrInvalidNameLength = "nil name or too long for chain or token,value:%s,actualLength:%d,max:%d"
)
