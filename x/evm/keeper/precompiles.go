package keeper

import (
	"fmt"

	assetsprecompile "github.com/ExocoreNetwork/exocore/precompiles/assets"
	avsManagerPrecompile "github.com/ExocoreNetwork/exocore/precompiles/avs"
	taskPrecompile "github.com/ExocoreNetwork/exocore/precompiles/avsTask"
	blsPrecompile "github.com/ExocoreNetwork/exocore/precompiles/bls"
	delegationprecompile "github.com/ExocoreNetwork/exocore/precompiles/delegation"
	rewardPrecompile "github.com/ExocoreNetwork/exocore/precompiles/reward"
	slashPrecompile "github.com/ExocoreNetwork/exocore/precompiles/slash"
	stakingStateKeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	avsManagerKeeper "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	taskKeeper "github.com/ExocoreNetwork/exocore/x/avstask/keeper"
	delegationKeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	rewardKeeper "github.com/ExocoreNetwork/exocore/x/reward/keeper"
	exoslashKeeper "github.com/ExocoreNetwork/exocore/x/slash/keeper"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	channelkeeper "github.com/cosmos/ibc-go/v7/modules/core/04-channel/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	ics20precompile "github.com/evmos/evmos/v14/precompiles/ics20"
	transferkeeper "github.com/evmos/evmos/v14/x/ibc/transfer/keeper"
	"golang.org/x/exp/maps"
)

const (
	BaseGas uint64 = 6000
)

// AvailablePrecompiles returns the list of all available precompiled contracts.
// NOTE: this should only be used during initialization of the Keeper.
func AvailablePrecompiles(
	authzKeeper authzkeeper.Keeper,
	transferKeeper transferkeeper.Keeper,
	channelKeeper channelkeeper.Keeper,
	delegationKeeper delegationKeeper.Keeper,
	assetskeeper stakingStateKeeper.Keeper,
	slashKeeper exoslashKeeper.Keeper,
	rewardKeeper rewardKeeper.Keeper,
	avsManagerKeeper avsManagerKeeper.Keeper,
	taskKeeper taskKeeper.Keeper,
) map[common.Address]vm.PrecompiledContract {
	// Clone the mapping from the latest EVM fork.
	precompiles := maps.Clone(vm.PrecompiledContractsBerlin)

	ibcTransferPrecompile, err := ics20precompile.NewPrecompile(
		transferKeeper,
		channelKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load ICS20 precompile: %w", err))
	}

	assetsPrecompile, err := assetsprecompile.NewPrecompile(
		assetskeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load deposit precompile: %w", err))
	}
	delegationPrecompile, err := delegationprecompile.NewPrecompile(
		assetskeeper,
		delegationKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load delegation precompile: %w", err))
	}

	slashPrecompile, err := slashPrecompile.NewPrecompile(
		assetskeeper,
		slashKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load slash precompile: %w", err))
	}
	rewardPrecompile, err := rewardPrecompile.NewPrecompile(
		assetskeeper,
		rewardKeeper,
		authzKeeper,
	)
	if err != nil {
		panic(fmt.Errorf("failed to load reward precompile: %w", err))
	}
	avsManagerPrecompile, err := avsManagerPrecompile.NewPrecompile(avsManagerKeeper, authzKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to load avsManager precompile: %w", err))
	}
	taskPrecompile, err := taskPrecompile.NewPrecompile(authzKeeper, taskKeeper, avsManagerKeeper)
	if err != nil {
		panic(fmt.Errorf("failed to load  reward precompile: %w", err))
	}
	blsPrecompile, err := blsPrecompile.NewPrecompile(BaseGas)
	if err != nil {
		panic(fmt.Errorf("failed to load bls precompile: %v", err))
	}
	precompiles[slashPrecompile.Address()] = slashPrecompile
	precompiles[rewardPrecompile.Address()] = rewardPrecompile
	precompiles[assetsPrecompile.Address()] = assetsPrecompile
	precompiles[delegationPrecompile.Address()] = delegationPrecompile
	precompiles[avsManagerPrecompile.Address()] = avsManagerPrecompile
	precompiles[taskPrecompile.Address()] = taskPrecompile
	precompiles[ibcTransferPrecompile.Address()] = ibcTransferPrecompile
	precompiles[blsPrecompile.Address()] = blsPrecompile
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
