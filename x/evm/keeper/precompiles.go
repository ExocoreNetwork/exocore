package keeper

import (
	"fmt"

	stakingStateKeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	rewardKeeper "github.com/ExocoreNetwork/exocore/x/reward/keeper"
	withdrawKeeper "github.com/ExocoreNetwork/exocore/x/withdraw/keeper"

	delegationKeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	depositKeeper "github.com/ExocoreNetwork/exocore/x/deposit/keeper"

	"golang.org/x/exp/maps"

	delegationprecompile "github.com/ExocoreNetwork/exocore/precompiles/delegation"
	depositprecompile "github.com/ExocoreNetwork/exocore/precompiles/deposit"
	rewardPrecompile "github.com/ExocoreNetwork/exocore/precompiles/reward"
	slashPrecompile "github.com/ExocoreNetwork/exocore/precompiles/slash"
	withdrawPrecompile "github.com/ExocoreNetwork/exocore/precompiles/withdraw"
	exoslashKeeper "github.com/ExocoreNetwork/exocore/x/slash/keeper"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	channelkeeper "github.com/cosmos/ibc-go/v7/modules/core/04-channel/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	distprecompile "github.com/evmos/evmos/v14/precompiles/distribution"
	ics20precompile "github.com/evmos/evmos/v14/precompiles/ics20"
	stakingprecompile "github.com/evmos/evmos/v14/precompiles/staking"
	transferkeeper "github.com/evmos/evmos/v14/x/ibc/transfer/keeper"
)

// AvailablePrecompiles returns the list of all available precompiled contracts.
// NOTE: this should only be used during initialization of the Keeper.
func AvailablePrecompiles(
	stakingKeeper stakingkeeper.Keeper,
	distributionKeeper distributionkeeper.Keeper,
	authzKeeper authzkeeper.Keeper,
	transferKeeper transferkeeper.Keeper,
	channelKeeper channelkeeper.Keeper,
	depositKeeper depositKeeper.Keeper,
	delegationKeeper delegationKeeper.Keeper,
	stakingStateKeeper stakingStateKeeper.Keeper,
	withdrawKeeper withdrawKeeper.Keeper,
	slashKeeper exoslashKeeper.Keeper,
	rewardKeeper rewardKeeper.Keeper,
) map[common.Address]vm.PrecompiledContract {
	// Clone the mapping from the latest EVM fork.
	precompiles := maps.Clone(vm.PrecompiledContractsBerlin)

	stakingPrecompile, err := stakingprecompile.NewPrecompile(stakingKeeper, authzKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to load staking precompile: %w", err))
	}

	distributionPrecompile, err := distprecompile.NewPrecompile(distributionKeeper, authzKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to load distribution precompile: %w", err))
	}

	ibcTransferPrecompile, err := ics20precompile.NewPrecompile(
		transferKeeper,
		channelKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load ICS20 precompile: %w", err))
	}

	// add exoCore chain preCompiles
	depositPrecompile, err := depositprecompile.NewPrecompile(
		stakingStateKeeper,
		depositKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load deposit precompile: %w", err))
	}
	delegationPrecompile, err := delegationprecompile.NewPrecompile(
		stakingStateKeeper,
		delegationKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load delegation precompile: %w", err))
	}
	withdrawPrecompile, err := withdrawPrecompile.NewPrecompile(
		stakingStateKeeper,
		withdrawKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load  withdraw precompile: %w", err))
	}
	slashPrecompile, err := slashPrecompile.NewPrecompile(
		stakingStateKeeper,
		slashKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load  slash precompile: %w", err))
	}
	rewardPrecompile, err := rewardPrecompile.NewPrecompile(
		stakingStateKeeper,
		rewardKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load  reward precompile: %w", err))
	}
	precompiles[slashPrecompile.Address()] = slashPrecompile
	precompiles[rewardPrecompile.Address()] = rewardPrecompile
	precompiles[withdrawPrecompile.Address()] = withdrawPrecompile
	precompiles[depositPrecompile.Address()] = depositPrecompile
	precompiles[delegationPrecompile.Address()] = delegationPrecompile

	precompiles[stakingPrecompile.Address()] = stakingPrecompile
	precompiles[distributionPrecompile.Address()] = distributionPrecompile
	precompiles[ibcTransferPrecompile.Address()] = ibcTransferPrecompile
	return precompiles
}

// WithPrecompiles sets the available precompiled contracts.
func (k *Keeper) WithPrecompiles(
	precompiles map[common.Address]vm.PrecompiledContract,
) *Keeper {
	if k.precompiles != nil {
		panic("available precompiles map already set")
	}

	if len(precompiles) == 0 {
		panic("empty precompiled contract map")
	}

	k.precompiles = precompiles
	return k
}

// Precompiles returns the subset of the available precompiled contracts that
// are active given the current parameters.
func (k Keeper) Precompiles(
	activePrecompiles ...common.Address,
) map[common.Address]vm.PrecompiledContract {
	activePrecompileMap := make(map[common.Address]vm.PrecompiledContract)

	for _, address := range activePrecompiles {
		precompile, ok := k.precompiles[address]
		if !ok {
			panic(fmt.Sprintf("precompiled contract not initialized: %s", address))
		}

		activePrecompileMap[address] = precompile
	}

	return activePrecompileMap
}
