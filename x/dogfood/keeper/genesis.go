package keeper

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(
	ctx sdk.Context,
	genState types.GenesisState,
) []abci.ValidatorUpdate {
	k.SetParams(ctx, genState.Params)
	// the `params` validator is not super useful to validate state level information
	// so, it must be done here. by extension, the `InitGenesis` of the epochs module
	// should be called before that of this module.
	epochID := genState.Params.EpochIdentifier
	epochInfo, found := k.epochsKeeper.GetEpochInfo(ctx, epochID)
	if !found {
		// the panic is suitable here because it is being done at genesis, when the node
		// is not running. it means that the genesis file is malformed.
		panic(fmt.Sprintf("epoch info not found %s", epochID))
	}
	// apply the same logic to the staking assets.
	for _, assetID := range genState.Params.AssetIDs {
		if !k.restakingKeeper.IsStakingAsset(ctx, assetID) {
			panic(fmt.Errorf("staking param %s not found in assets module", assetID))
		}
	}
	// at genesis, not chain restarts, each operator may not necessarily be an initial
	// validator. this is because the operator may not have enough minimum self delegation
	// to be considered, or may not be in the top N operators. so checking that count here
	// is meaningless as well.
	out := make([]abci.ValidatorUpdate, len(genState.ValSet))
	for _, val := range genState.ValSet {
		// #nosec G703 // already validated
		consKey, _ := operatortypes.HexStringToPubKey(val.PublicKey)
		// #nosec G703 // this only fails if the key is of a type not already defined.
		consAddr, _ := operatortypes.TMCryptoPublicKeyToConsAddr(consKey)
		// if GetOperatorAddressForChainIDAndConsAddr returns found, it means
		// that the operator is registered and also (TODO) that it has opted into
		// the dogfood AVS.
		found, _ := k.operatorKeeper.GetOperatorAddressForChainIDAndConsAddr(
			ctx, ctx.ChainID(), consAddr,
		)
		if !found {
			panic(fmt.Sprintf("operator not found: %s", consAddr))
		}
		out = append(out, abci.ValidatorUpdate{
			PubKey: *consKey,
			Power:  val.Power,
		})
	}
	for i := range genState.OptOutExpiries {
		obj := genState.OptOutExpiries[i]
		epoch := obj.Epoch
		if epoch < epochInfo.CurrentEpoch {
			panic(fmt.Sprintf("epoch %d is in the past", epoch))
		}
		for _, addr := range obj.OperatorAccAddrs {
			// #nosec G703 // already validated
			operatorAddr, _ := sdk.AccAddressFromBech32(addr)
			k.AppendOptOutToFinish(ctx, epoch, operatorAddr)
			k.SetOperatorOptOutFinishEpoch(ctx, operatorAddr, epoch)
		}
	}
	for i := range genState.ConsensusAddrsToPrune {
		obj := genState.ConsensusAddrsToPrune[i]
		epoch := obj.Epoch
		if epoch < epochInfo.CurrentEpoch {
			panic(fmt.Sprintf("epoch %d is in the past", epoch))
		}
		for _, addr := range obj.ConsAddrs {
			// #nosec G703 // already validated
			accAddr, _ := sdk.ConsAddressFromBech32(addr)
			k.AppendConsensusAddrToPrune(ctx, epoch, accAddr)
		}
	}
	for i := range genState.UndelegationMaturities {
		obj := genState.UndelegationMaturities[i]
		epoch := obj.Epoch
		if epoch < epochInfo.CurrentEpoch {
			panic(fmt.Sprintf("epoch %d is in the past", epoch))
		}
		for _, recordKey := range obj.UndelegationRecordKeys {
			// #nosec G703 // already validated
			recordKeyBytes, _ := hexutil.Decode(recordKey)
			k.AppendUndelegationToMature(ctx, epoch, recordKeyBytes)
			k.SetUndelegationMaturityEpoch(ctx, recordKeyBytes, epoch)
		}
	}
	// ApplyValidatorChanges only gets changes and hence the vote power must be set here.
	k.SetLastTotalPower(ctx, genState.LastTotalPower)

	// ApplyValidatorChanges will sort it internally
	return k.ApplyValidatorChanges(
		ctx, out,
	)
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetDogfoodParams(ctx)
	validators := []types.GenesisValidator{}
	k.IterateBondedValidatorsByPower(ctx, func(i int64, val stakingtypes.ValidatorI) bool {
		// #nosec G703 // already validated
		pubKey, _ := val.ConsPubKey()
		// #nosec G703 // already validated
		convKey, _ := cryptocodec.ToTmPubKeyInterface(pubKey)
		validators = append(validators,
			types.GenesisValidator{
				PublicKey: hexutil.Encode(convKey.Bytes()),
				Power:     val.GetConsensusPower(sdk.DefaultPowerReduction),
			},
		)
		return true
	})
	return types.NewGenesis(
		k.GetDogfoodParams(ctx),
		validators,
		k.GetAllOptOutsToFinish(ctx),
		k.GetAllConsAddrsToPrune(ctx),
		k.GetAllUndelegationsToMature(ctx),
		k.GetLastTotalPower(ctx),
	)
}
