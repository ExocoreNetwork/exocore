package types

import (
	"github.com/ethereum/go-ethereum/common"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "avs"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey   = "mem_avs"
	prefixAVSInfo = iota + 1
	prefixAVSOperatorInfo
)

// ModuleAddress is the native module address for EVM
var (
	ModuleAddress            common.Address
	KeyPrefixAVSInfo         = []byte{prefixAVSInfo}
	KeyPrefixAVSOperatorInfo = []byte{prefixAVSOperatorInfo}
)

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}