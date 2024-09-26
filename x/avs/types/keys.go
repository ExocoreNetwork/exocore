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
)

const (
	prefixAVSInfo = iota + 5
	prefixParams
	prefixAVSTaskInfo
	prefixOperatePub
	prefixAVSAddressToChainID
	LatestTaskNum
	TaskResult
	TaskChallengeResult
)

// ModuleAddress is the native module address for EVM
var (
	ModuleAddress        common.Address
	KeyPrefixAVSInfo     = []byte{prefixAVSInfo}
	KeyPrefixParams      = []byte{prefixParams}
	ParamsKey            = []byte("Params")
	KeyPrefixAVSTaskInfo = []byte{prefixAVSTaskInfo}
	KeyPrefixOperatePub  = []byte{prefixOperatePub}
	// KeyPrefixAVSAddressToChainID is used to store the reverse lookup from AVS address to chainID.
	KeyPrefixAVSAddressToChainID = []byte{prefixAVSAddressToChainID}
	KeyPrefixLatestTaskNum       = []byte{LatestTaskNum}
	KeyPrefixTaskResult          = []byte{TaskResult}
	KeyPrefixTaskChallengeResult = []byte{TaskChallengeResult}
)

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}
