package keeper

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (*types.Params, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixParams)
	ifExist := store.Has(types.ParamsKey)
	if !ifExist {
		return nil, fmt.Errorf("params %s not found", types.KeyPrefixParams)
	}

	value := store.Get(types.ParamsKey)

	ret := &types.Params{}
	k.cdc.MustUnmarshal(value, ret)
	return ret, nil
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixParams)
	bz := k.cdc.MustMarshal(params)
	store.Set(types.ParamsKey, bz)
	return nil
}

type AVSRegisterOrDeregisterParams struct {
	AvsName             string
	AvsAddress          string
	MinStakeAmount      uint64
	TaskAddr            string
	SlashContractAddr   string
	RewardContractAddr  string
	AvsOwnerAddress     []string
	AssetID             []string
	UnbondingPeriod     uint64
	MinSelfDelegation   uint64
	EpochIdentifier     string
	MinOptInOperators   uint64
	MinTotalStakeAmount uint64
	StartingEpoch       []string
	CallerAddress       string
	ChainID             string
	AvsReward           uint64
	AvsSlash            uint64
	Action              uint64
}

type OperatorOptParams struct {
	Name            string
	BlsPublicKey    string
	IsRegistered    bool
	Action          uint64
	OperatorAddress string
	Status          string
	AvsAddress      string
}

type TaskParams struct {
	TaskContractAddress string
	TaskName            string
	StartingEpoch       string
	Data                []byte
	TaskID              string
	TaskResponsePeriod  uint64
	TaskChallengePeriod uint64
	ThresholdPercentage uint64
	CallerAddress       string
}
type BlsParams struct {
	Operator                      string
	Name                          string
	PubKey                        []byte
	PubkeyRegistrationSignature   []byte
	PubkeyRegistrationMessageHash []byte
}

type ProofParams struct {
	TaskID              string
	TaskContractAddress string
	AvsAddress          string
	Aggregator          string
	OperatorStatus      []OperatorStatusParams
	CallerAddress       string
}
type OperatorStatusParams struct {
	OperatorAddress string
	Status          string
	ProofData       string
}

const (
	RegisterAction   = 1
	DeRegisterAction = 2
	UpdateAction     = 3
)
