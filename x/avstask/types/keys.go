package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "task"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

const (
	prefixAVSTaskInfo = iota + 1
	prefixOperatePub
)

var (
	// KeyPrefixAVSTaskInfo key-value: taskAddr->AVSTaskInfo
	KeyPrefixAVSTaskInfo = []byte{prefixAVSTaskInfo}
	KeyPrefixOperatePub  = []byte{prefixOperatePub}
)

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}
