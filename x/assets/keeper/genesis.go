package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/ExocoreNetwork/exocore/x/assets/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ValidDepositByInitType(assetID string, info types.StakerAssetInfo, tokenState map[string]types.TotalSupplyAndStaking) error {
	if !k.isGeneralInit {
		// check that deposit amount does not exceed supply.
		if info.TotalDepositAmount.GT(tokenState[assetID].TotalSupply) {
			return errorsmod.Wrapf(
				types.ErrInvalidGenesisData,
				"deposit amount exceeds max supply for %s: %+v",
				assetID, info,
			)
		}
		// at genesis (not chain restart), there is no unbonding amount.
		if !info.WaitUnbondingAmount.IsZero() {
			return errorsmod.Wrapf(
				types.ErrInvalidGenesisData,
				"non-zero unbonding amount for %s: %s when initializing from the bootStrap contract",
				assetID, info.WaitUnbondingAmount,
			)
		}
		// check that the withdrawable amount and the deposited amount are equal.
		// this is because this module's genesis only sets up free deposits.
		// the delegation module bonds them, thereby altering the withdrawable amount.
		if !info.WithdrawableAmount.Equal(info.TotalDepositAmount) {
			return errorsmod.Wrapf(
				types.ErrInvalidGenesisData,
				"withdrawable amount is not equal to total deposit amount for %s: %+v when initializing from the bootStrap contract",
				assetID, info,
			)
		}
	} else {
		// check that deposit amount does not exceed the total staking amount
		// when initializing from the general exported genesis file
		if info.TotalDepositAmount.GT(tokenState[assetID].TotalStaking) {
			return errorsmod.Wrapf(
				types.ErrInvalidGenesisData,
				"deposit amount exceeds the total staking amount for %s: %+v",
				assetID, info,
			)
		}
		// the sum of `WaitUnbondingAmount` and `WithdrawableAmount` shouldn't be greater than the `TotalDepositAmount`
		// when initializing from the general exported genesis file
		if info.WaitUnbondingAmount.Add(info.WithdrawableAmount).GT(info.TotalDepositAmount) {
			return errorsmod.Wrapf(
				types.ErrInvalidGenesisData,
				"the sum of withdrawable amount and unbonding amount is greater than the total deposit amount for %s: %+v when initializing from the general exported genesis",
				assetID, info,
			)
		}
	}
	return nil
}

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	if err := k.SetParams(ctx, &data.Params); err != nil {
		panic(err)
	}
	// TODO(mm): is it possible to optimize / speed up this process?
	// client_chain.go
	for i := range data.ClientChains {
		info := data.ClientChains[i]
		if err := k.SetClientChainInfo(ctx, &info); err != nil {
			panic(err)
		}
	}
	// client_chain_asset.go
	tokenState := make(map[string]types.TotalSupplyAndStaking, 0)
	for i := range data.Tokens {
		info := data.Tokens[i]
		// the StakingTotalAmount should be zero when init from the bootStrap
		if !k.isGeneralInit && !info.StakingTotalAmount.IsZero() {
			panic(errorsmod.Wrapf(
				types.ErrInvalidGenesisData,
				"non-zero deposit amount for asset %s when initializing from the bootStrap contract",
				info.AssetBasicInfo.Address,
			))
		}
		if err := k.SetStakingAssetInfo(ctx, &info); err != nil {
			panic(err)
		}
		_, assetID := types.GetStakeIDAndAssetIDFromStr(
			info.AssetBasicInfo.LayerZeroChainID,
			"", info.AssetBasicInfo.Address,
		)
		tokenState[assetID] = types.TotalSupplyAndStaking{
			TotalSupply:  info.AssetBasicInfo.TotalSupply,
			TotalStaking: info.StakingTotalAmount,
		}
	}
	// staker_asset.go (deposits)
	// we simulate the behavior of the depositKeeper.Deposit call
	// it constructs the stakerID and the assetID, which we have validated previously.
	// it checks that the deposited amount is not negative, which we have already done.
	// and that the asset is registered, which we have also already done.
	for _, deposit := range data.Deposits {
		stakerID := deposit.StakerID
		for _, depositsByStaker := range deposit.Deposits {
			assetID := depositsByStaker.AssetID
			info := depositsByStaker.Info
			err := k.ValidDepositByInitType(assetID, info, tokenState)
			if err != nil {
				panic(err)
			}
			infoAsChange := types.DeltaStakerSingleAsset(info)
			// set the deposited and free values for the staker
			if err := k.UpdateStakerAssetState(
				ctx, stakerID, assetID, infoAsChange,
			); err != nil {
				panic(err)
			}
			// now for the asset, increase the deposit value
			// This should only be called when initializing from the bootStrap contract
			// because the `TotalDepositAmount` will be initialized when
			// initializing the tokens information from the general exported genesis file.
			if !k.isGeneralInit {
				if err := k.UpdateStakingAssetTotalAmount(
					ctx, assetID, info.TotalDepositAmount,
				); err != nil {
					panic(err)
				}
			}
		}
	}

	// initialize the operators assets from the genesis
	if !k.isGeneralInit && len(data.OperatorAssets) != 0 {
		panic(errorsmod.Wrap(
			types.ErrInvalidGenesisData,
			"the operator assets should be null when initializing from the bootStrap contract",
		))
	}
	for _, assets := range data.OperatorAssets {
		for _, assetInfo := range assets.AssetsState {
			// #nosec G703 // already validated
			accAddress, _ := sdk.AccAddressFromBech32(assets.Operator)
			infoAsChange := types.DeltaOperatorSingleAsset(assetInfo.Info)
			err := k.UpdateOperatorAssetState(ctx, accAddress, assetInfo.AssetID, infoAsChange)
			if err != nil {
				panic(err)
			}
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	res := types.GenesisState{}
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	res.Params = *params

	allClientChains, err := k.GetAllClientChainInfo(ctx)
	if err != nil {
		panic(err)
	}
	res.ClientChains = allClientChains

	allAssets, err := k.GetAllStakingAssetsInfo(ctx)
	if err != nil {
		panic(err)
	}
	res.Tokens = allAssets

	allDeposits, err := k.AllDeposits(ctx)
	if err != nil {
		panic(err)
	}
	res.Deposits = allDeposits

	operatorAssets, err := k.AllOperatorAssets(ctx)
	if err != nil {
		panic(err)
	}
	res.OperatorAssets = operatorAssets
	return &res
}
