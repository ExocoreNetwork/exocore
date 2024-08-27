package types

import (
	epochstypes "github.com/ExocoreNetwork/exocore/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "coordinator"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_coordinator"

	// PortID is the default port id that the module binds to
	PortID = "coordinator"

	// SubscriberRewardsPool is the address that receives the rewards from the subscriber
	// chains. Technically, it is possible for the subscriber chain to send these rewards
	// directly to the FeeCollector, but this intermediate step allows the coordinator
	// module to ensure that the rewards are actually sent to us.
	SubscriberRewardsPool = "subscriber_rewards_pool"
)

const (
	// ParamsBytePrefix is the prefix for the coordinator module's parameters key
	ParamsBytePrefix byte = iota + 1
	// PendingSubscriberChainBytePrefix is the prefix for the coordinator module's pending subscriber chain key
	PendingSubscriberChainBytePrefix
	// ClientForChainBytePrefix is the prefix to store mapping from chain to client
	ClientForChainBytePrefix
	// SubSlashFractionDowntimeBytePrefix is the prefix to store the slashing fraction for downtime
	// against a particular chain.
	SubSlashFractionDowntimeBytePrefix
	// SubSlashFractionDoubleSignBytePrefix is the prefix to store the slashing fraction for double sign
	// against a particular chain.
	SubSlashFractionDoubleSignBytePrefix
	// SubDowntimeJailDurationBytePrefix is the prefix to store the downtime jail duration for a chain
	SubDowntimeJailDurationBytePrefix
	// SubscriberGenesisBytePrefix is the prefix for the subscriber genesis key
	SubscriberGenesisBytePrefix
	// InitTimeoutBytePrefix is the prefix for the init timeout key
	InitTimeoutBytePrefix
)

// AppendMany appends a variable number of byte slices together
func AppendMany(byteses ...[]byte) (out []byte) {
	for _, bytes := range byteses {
		out = append(out, bytes...)
	}
	return out
}

func ParamsKey() []byte {
	return []byte{ParamsBytePrefix}
}

// PendingSubscriberChainKey is the key used to store subscriber chains, which are scheduled
// to begin with the starting of the epoch with identifier and number. Since the data
// is stored alphabetically, this key structure is apt.
func PendingSubscriberChainKey(epochIdentifier string, epochNumber uint64) []byte {
	return AppendMany(
		[]byte{PendingSubscriberChainBytePrefix},
		[]byte(epochIdentifier),
		sdk.Uint64ToBigEndian(epochNumber),
	)
}

// ClientForChainKey returns the key under which the clientId for the given chainId is stored.
func ClientForChainKey(chainId string) []byte {
	return append([]byte{ClientForChainBytePrefix}, []byte(chainId)...)
}

// SubSlashFractionDowntimeKey returns the key under which the slashing fraction for downtime
// against a particular chain is stored.
func SubSlashFractionDowntimeKey(chainId string) []byte {
	return append([]byte{SubSlashFractionDowntimeBytePrefix}, []byte(chainId)...)
}

// SubSlashFractionDoubleSignKey returns the key under which the slashing fraction for double sign
// against a particular chain is stored.
func SubSlashFractionDoubleSignKey(chainId string) []byte {
	return append([]byte{SubSlashFractionDoubleSignBytePrefix}, []byte(chainId)...)
}

// SubDowntimeJailDurationKey returns the key under which the downtime jail duration for a chain
// is stored.
func SubDowntimeJailDurationKey(chainId string) []byte {
	return append([]byte{SubDowntimeJailDurationBytePrefix}, []byte(chainId)...)
}

// SubscriberGenesisKey returns the key under which the genesis state for a subscriber chain is stored.
func SubscriberGenesisKey(chainId string) []byte {
	return append([]byte{SubscriberGenesisBytePrefix}, []byte(chainId)...)
}

// InitTimeoutEpochKey returns the key under which the list of chains which will timeout (if not
// initialized by then) at the beginning of the epoch is stored.
func InitTimeoutEpochKey(epoch epochstypes.Epoch) []byte {
	return AppendMany(
		[]byte{InitTimeoutBytePrefix},
		[]byte(epoch.EpochIdentifier),
		sdk.Uint64ToBigEndian(epoch.EpochNumber),
	)
}
