package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName defines the module name
	ModuleName = "feedistribution"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for distribution
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey            = "mem_feedistribute"
	ProtocolPoolModuleName = "protocolpool"
)

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

const (
	// EpochIdentifier defines the epoch identifier for fee distribution module
	prefixParams = iota + 1
	prefixEpochIdentifier
)

var (
	KeyPrefixParams          = KeyPrefix(prefixParams)
	KeyPrefixEpochIdentifier = KeyPrefix(prefixEpochIdentifier)
)
var (
	EventTypeCommission         = "commission"
	EventTypeSetWithdrawAddress = "set_withdraw_address"
	EventTypeRewards            = "rewards"
	EventTypeWithdrawRewards    = "withdraw_rewards"
	EventTypeWithdrawCommission = "withdraw_commission"
	EventTypeProposerReward     = "proposer_reward"

	AttributeKeyWithdrawAddress = "withdraw_address"
	AttributeKeyValidator       = "validator"
	AttributeKeyDelegator       = "delegator"
)

func KeyPrefix(p uint64) []byte {
	return []byte(p)
}
