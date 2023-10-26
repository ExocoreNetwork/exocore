package reward

import (
	"math/rand"

	"github.com/exocore/testutil/sample"
	rewardsimulation "github.com/exocore/x/reward/simulation"
	"github.com/exocore/x/reward/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = rewardsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
    opWeightMsgClaimRewardRequest = "op_weight_msg_claim_reward_request"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimRewardRequest int = 100

	opWeightMsgRewardDetail = "op_weight_msg_reward_detail"
	// TODO: Determine the simulation weight value
	defaultWeightMsgRewardDetail int = 100

	opWeightMsgClaimRewardResponse = "op_weight_msg_claim_reward_response"
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimRewardResponse int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	rewardGenesis := types.GenesisState{
		Params:	types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&rewardGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgClaimRewardRequest int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgClaimRewardRequest, &weightMsgClaimRewardRequest, nil,
		func(_ *rand.Rand) {
			weightMsgClaimRewardRequest = defaultWeightMsgClaimRewardRequest
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimRewardRequest,
		rewardsimulation.SimulateMsgClaimRewardRequest(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgRewardDetail int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgRewardDetail, &weightMsgRewardDetail, nil,
		func(_ *rand.Rand) {
			weightMsgRewardDetail = defaultWeightMsgRewardDetail
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgRewardDetail,
		rewardsimulation.SimulateMsgRewardDetail(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClaimRewardResponse int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgClaimRewardResponse, &weightMsgClaimRewardResponse, nil,
		func(_ *rand.Rand) {
			weightMsgClaimRewardResponse = defaultWeightMsgClaimRewardResponse
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimRewardResponse,
		rewardsimulation.SimulateMsgClaimRewardResponse(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
	    simulation.NewWeightedProposalMsg(
	opWeightMsgClaimRewardRequest,
	defaultWeightMsgClaimRewardRequest,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		rewardsimulation.SimulateMsgClaimRewardRequest(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgRewardDetail,
	defaultWeightMsgRewardDetail,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		rewardsimulation.SimulateMsgRewardDetail(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
simulation.NewWeightedProposalMsg(
	opWeightMsgClaimRewardResponse,
	defaultWeightMsgClaimRewardResponse,
	func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
		rewardsimulation.SimulateMsgClaimRewardResponse(am.accountKeeper, am.bankKeeper, am.keeper)
		return nil
	},
),
// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
