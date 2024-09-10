package types

import (
	"bytes"
	fmt "fmt"

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
	// SubscriberValidatorBytePrefix is the prefix for the subscriber validator key
	SubscriberValidatorBytePrefix
	// MaxValidatorsBytePrefix is the prefix for the max validators key
	MaxValidatorsBytePrefix
	// VscIDForChainBytePrefix is the prefix to go from chainID to vscID
	VscIDForChainBytePrefix
	// ChainIDToVscPacketsBytePrefix is the prefix for the vsc packets key for a chainID
	ChainIDToVscPacketsBytePrefix
	// VscTimeoutBytePrefix is the prefix for the vsc timeout key
	VscTimeoutBytePrefix
	// ConsKeysToPruneBytePrefix is the prefix for the consensus keys to prune key
	ConsKeysToPruneBytePrefix
	// MaturityVscIDForChainIDConsAddrBytePrefix is the prefix for the vsc id for chain cons addr key
	MaturityVscIDForChainIDConsAddrBytePrefix
	// UndelegationsToReleaseBytePrefix is the prefix for the undelegations to release key
	UndelegationsToReleaseBytePrefix
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
	return utils.AppendMany( // safe to do, since...
		[]byte{InitTimeoutBytePrefix},            // size 1
		[]byte(epoch.EpochIdentifier),            // size unknown
		sdk.Uint64ToBigEndian(epoch.EpochNumber), // size 8
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

// SubscriberValidatorKey returns the key for the subscriber validator
// It is used to store the validator object for the subscriber chain, indexed by
// prefix + len(chainID) + chainID + consensusAddr
func SubscriberValidatorKey(chainID string, consensusAddr []byte) []byte {
	return utils.AppendMany(
		[]byte{SubscriberValidatorBytePrefix},
		utils.ChainIDWithLenKey(chainID),
		consensusAddr,
	)
}

// MaxValidatorsKey returns the key for the max validators
func MaxValidatorsKey(chainID string) []byte {
	return append([]byte{MaxValidatorsBytePrefix}, []byte(chainID)...)
}

// VscIDForChainKey returns the key for the vsc id to chain
func VscIDForChainKey(chainID string) []byte {
	return append([]byte{VscIDForChainBytePrefix}, []byte(chainID)...)
}

// ChainIDToVscPacketsKey returns the key for the vsc packets for a chain
func ChainIDToVscPacketsKey(chainID string) []byte {
	return append([]byte{ChainIDToVscPacketsBytePrefix}, []byte(chainID)...)
}

// VscTimeoutKey returns the key for the vsc timeout
func VscTimeoutKey(chainID string, vscID uint64) []byte {
	return utils.AppendMany(
		[]byte{VscTimeoutBytePrefix},
		[]byte(chainID),
		sdk.Uint64ToBigEndian(vscID),
	)
}

// ParseVscTimeoutKey parses the chainID and vscID from the key of the format
// prefix + chainID + vscID
func ParseVscTimeoutKey(bz []byte) (chainID string, vscID uint64, err error) {
	return ParseChainIDAndUintIDKey(VscTimeoutBytePrefix, bz)
}

// ParseChainIDAndUintIDKey returns the chain ID and uint ID for a ChainIdAndUintId key
func ParseChainIDAndUintIDKey(prefix byte, bz []byte) (string, uint64, error) {
	expectedPrefix := []byte{prefix}
	prefixL := len(expectedPrefix)
	if len(bz) < prefixL+8 { // for uint64
		return "", 0, fmt.Errorf("invalid key length; expected at least %d bytes, got: %d", prefixL+8, len(bz))
	}
	if prefix := bz[:prefixL]; !bytes.Equal(prefix, expectedPrefix) {
		return "", 0, fmt.Errorf("invalid prefix; expected: %X, got: %X", expectedPrefix, prefix)
	}
	uintID := sdk.BigEndianToUint64(bz[len(bz)-8:])
	chainID := string(bz[prefixL : len(bz)-8])
	return chainID, uintID, nil
}

// ConsAddrsToPruneKey returns the key for the consensus keys to prune, indexed by the
// chainID + vscID as the key.
func ConsAddrsToPruneKey(chainID string, vscID uint64) []byte {
	return utils.AppendMany(
		[]byte{ConsKeysToPruneBytePrefix},
		[]byte(chainID),
		sdk.Uint64ToBigEndian(vscID),
	)
}

// MaturityVscIDForChainIDConsAddrKey returns the key for the vsc id for chain cons addr
func MaturityVscIDForChainIDConsAddrKey(chainID string, consAddr sdk.ConsAddress) []byte {
	return utils.AppendMany(
		[]byte{MaturityVscIDForChainIDConsAddrBytePrefix},
		[]byte(chainID),
		consAddr.Bytes(),
	)
}

// UndelegationsToReleaseKey returns the key for the undelegations to release
func UndelegationsToReleaseKey(chainID string, vscID uint64) []byte {
	return utils.AppendMany(
		[]byte{UndelegationsToReleaseBytePrefix},
		[]byte(chainID),
		sdk.Uint64ToBigEndian(vscID),
	)
}
