package types

import (
	"strings"

	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// ChainIDCode is the "fake" code used to mark a generated AVS address as occupied by a contract.
	ChainIDCode = []byte("chain-id-code")
	// ChainIDCodeHash is the hash of the ChainIDCode.
	ChainIDCodeHash = crypto.Keccak256Hash(ChainIDCode)
	// ChainIDPrefix
	ChainIDPrefix = []byte("chain-id-prefix")
)

type AVSRegisterOrDeregisterParams struct {
	// AvsName is the name of the AVS as an arbitrary string.
	AvsName string
	// AvsAddress is the hex address of the AVS.
	AvsAddress string
	// MinStakeAmount is the minimum amount of stake for a task to be considered valid.
	MinStakeAmount uint64
	// TaskAddr is the hex address of the task contract.
	TaskAddr string
	// SlashContractAddr is the hex address of the slash contract.
	SlashContractAddr string
	// RewardContractAddr is the hex address of the reward contract.
	RewardContractAddr string
	// AvsOwnerAddress is the list of bech32 addresses of the AVS owners.
	AvsOwnerAddress []string
	// AssetID is the list of asset IDs that the AVS is allowed to use.
	AssetID             []string
	UnbondingPeriod     uint64
	MinSelfDelegation   uint64
	EpochIdentifier     string
	MinOptInOperators   uint64
	MinTotalStakeAmount uint64
	// CallerAddress is the bech32 address of the precompile caller.
	CallerAddress string
	ChainID       string
	AvsReward     uint64
	AvsSlash      uint64
	Action        uint64
}

// ChainIDWithoutRevision returns the chainID without the revision number.
// For example, "exocoretestnet_233-1" returns "exocoretestnet_233".
func ChainIDWithoutRevision(chainID string) string {
	if !ibcclienttypes.IsRevisionFormat(chainID) {
		return chainID
	}
	splitStr := strings.Split(chainID, "-")
	return splitStr[0]
}

// GenerateAVSAddr generates a hex AVS address based on the chainID.
// It returns a hex address as a string.
func GenerateAVSAddr(chainID string) string {
	return common.BytesToAddress(
		crypto.Keccak256(
			append(
				ChainIDPrefix,
				[]byte(chainID)...,
			),
		),
	).String()
}
