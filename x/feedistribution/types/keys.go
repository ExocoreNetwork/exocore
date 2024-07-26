package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
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
	prefixParams          = "feedistributionPrefixParams"
	prefixEpochIdentifier = "feedistrEpochPrefixEpochIdentifier"
	prefixFeePool         = "feePoolKey"
)

var (
	KeyPrefixParams                      = KeyPrefix(prefixParams)
	KeyPrefixEpochIdentifier             = KeyPrefix(prefixEpochIdentifier)
	FeePoolKey                           = KeyPrefix(prefixFeePool)
	ValidatorAccumulatedCommissionPrefix = []byte{0x00} // key for accumulated validator commission
	ValidatorCurrentRewardsPrefix        = []byte{0x01} // key for current validator rewards
	ValidatorOutstandingRewardsPrefix    = []byte{0x02} // key for outstanding rewards
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

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// GetValidatorAccumulatedCommissionKey creates the key for a validator's current commission.
func GetValidatorAccumulatedCommissionKey(v sdk.ValAddress) []byte {
	return append(ValidatorAccumulatedCommissionPrefix, address.MustLengthPrefix(v.Bytes())...)
}

// GetValidatorCurrentRewardsKey creates the key for a validator's current rewards.
func GetValidatorCurrentRewardsKey(v sdk.ValAddress) []byte {
	return append(ValidatorCurrentRewardsPrefix, address.MustLengthPrefix(v.Bytes())...)
}

// GetValidatorOutstandingRewardsKey creates the outstanding rewards key for a validator.
func GetValidatorOutstandingRewardsKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorOutstandingRewardsPrefix, address.MustLengthPrefix(valAddr.Bytes())...)
}
