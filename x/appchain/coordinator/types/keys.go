package types

import (
	"github.com/ExocoreNetwork/exocore/utils"
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
	// PortBytePrefix is the prefix for the port key
	PortBytePrefix
	// ChannelForChainBytePrefix is the prefix for the channel for chain key
	ChannelForChainBytePrefix
	// ChainForChannelBytePrefix is the prefix for the chain for channel key
	ChainForChannelBytePrefix
	// ChainInitTimeoutBytePrefix is the prefix for the chain init timeout key
	ChainInitTimeoutBytePrefix
	// InitChainHeightBytePrefix is the prefix for the init chain height key
	InitChainHeightBytePrefix
	// HeightToChainVscIDBytePrefix is the prefix for the height to chain id + vsc id key
	HeightToChainVscIDBytePrefix
	// SlashAcksBytePrefix is the prefix for the slashing acks key
	SlashAcksBytePrefix
)

// ParamsKey returns the key under which the coordinator module's parameters are stored.
func ParamsKey() []byte {
	return []byte{ParamsBytePrefix}
}

// PendingSubscriberChainKey is the key used to store subscriber chains, which are scheduled
// to begin with the starting of the epoch with identifier and number. Since the data
// is stored alphabetically, this key structure is apt.
func PendingSubscriberChainKey(epochIdentifier string, epochNumber uint64) []byte {
	return utils.AppendMany(
		[]byte{PendingSubscriberChainBytePrefix},
		[]byte(epochIdentifier),
		sdk.Uint64ToBigEndian(epochNumber),
	)
}

// ClientForChainKey returns the key under which the clientId for the given chainId is stored.
func ClientForChainKey(chainID string) []byte {
	return append([]byte{ClientForChainBytePrefix}, []byte(chainID)...)
}

// SubSlashFractionDowntimeKey returns the key under which the slashing fraction for downtime
// against a particular chain is stored.
func SubSlashFractionDowntimeKey(chainID string) []byte {
	return append([]byte{SubSlashFractionDowntimeBytePrefix}, []byte(chainID)...)
}

// SubSlashFractionDoubleSignKey returns the key under which the slashing fraction for double sign
// against a particular chain is stored.
func SubSlashFractionDoubleSignKey(chainID string) []byte {
	return append([]byte{SubSlashFractionDoubleSignBytePrefix}, []byte(chainID)...)
}

// SubDowntimeJailDurationKey returns the key under which the downtime jail duration for a chain
// is stored.
func SubDowntimeJailDurationKey(chainID string) []byte {
	return append([]byte{SubDowntimeJailDurationBytePrefix}, []byte(chainID)...)
}

// SubscriberGenesisKey returns the key under which the genesis state for a subscriber chain is stored.
func SubscriberGenesisKey(chainID string) []byte {
	return append([]byte{SubscriberGenesisBytePrefix}, []byte(chainID)...)
}

// InitTimeoutEpochKey returns the key under which the list of chains which will timeout (if not
// initialized by then) at the beginning of the epoch is stored.
func InitTimeoutEpochKey(epoch epochstypes.Epoch) []byte {
	return utils.AppendMany(
		[]byte{InitTimeoutBytePrefix},
		[]byte(epoch.EpochIdentifier),
		sdk.Uint64ToBigEndian(epoch.EpochNumber),
	)
}

// PortKey returns the key for the port (hello Harry Potter!)
func PortKey() []byte {
	return []byte{PortBytePrefix}
}

// ChannelForChainKey returns the key under which the ibc channel id
// for the given chainId is stored.
func ChannelForChainKey(chainID string) []byte {
	return append([]byte{ChannelForChainBytePrefix}, []byte(chainID)...)
}

// ChainForChannelKey returns the key under which the chainId
// for the given channelId is stored.
func ChainForChannelKey(channelID string) []byte {
	return append([]byte{ChainForChannelBytePrefix}, []byte(channelID)...)
}

// ChainInitTimeoutKey returns the key for the chain init timeout
func ChainInitTimeoutKey(chainID string) []byte {
	return append([]byte{ChainInitTimeoutBytePrefix}, []byte(chainID)...)
}

// InitChainHeightKey returns the key for the init chain height
func InitChainHeightKey(chainID string) []byte {
	return append([]byte{InitChainHeightBytePrefix}, []byte(chainID)...)
}

// HeightToChainVscIDKey returns the key for the height to chain id + vsc id
func HeightToChainVscIDKey(chainID string, vscID uint64) []byte {
	return utils.AppendMany(
		[]byte{HeightToChainVscIDBytePrefix},
		[]byte(chainID),
		sdk.Uint64ToBigEndian(vscID),
	)
}

// SlashAcksKey returns the key for the slashing acks
func SlashAcksKey(chainID string) []byte {
	return append(
		[]byte{SlashAcksBytePrefix},
		[]byte(chainID)...,
	)
}
