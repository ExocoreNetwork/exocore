package keeper

import (
	"bytes"
	"fmt"
	assetsprecompile "github.com/ExocoreNetwork/exocore/precompiles/assets"
	avsManagerPrecompile "github.com/ExocoreNetwork/exocore/precompiles/avs"
	blsPrecompile "github.com/ExocoreNetwork/exocore/precompiles/bls"
	delegationprecompile "github.com/ExocoreNetwork/exocore/precompiles/delegation"
	rewardPrecompile "github.com/ExocoreNetwork/exocore/precompiles/reward"
	slashPrecompile "github.com/ExocoreNetwork/exocore/precompiles/slash"
	stakingStateKeeper "github.com/ExocoreNetwork/exocore/x/assets/keeper"
	avsManagerKeeper "github.com/ExocoreNetwork/exocore/x/avs/keeper"
	delegationKeeper "github.com/ExocoreNetwork/exocore/x/delegation/keeper"
	rewardKeeper "github.com/ExocoreNetwork/exocore/x/reward/keeper"
	exoslashKeeper "github.com/ExocoreNetwork/exocore/x/slash/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	channelkeeper "github.com/cosmos/ibc-go/v7/modules/core/04-channel/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	ics20precompile "github.com/evmos/evmos/v16/precompiles/ics20"
	transferkeeper "github.com/evmos/evmos/v16/x/ibc/transfer/keeper"
	"golang.org/x/exp/maps"
	"sort"
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
	blsPrecompile, err := blsPrecompile.NewPrecompile(BaseGas)
	if err != nil {
		panic(fmt.Errorf("failed to load bls precompile: %v", err))
	}
	precompiles[slashPrecompile.Address()] = slashPrecompile
	precompiles[rewardPrecompile.Address()] = rewardPrecompile
	precompiles[assetsPrecompile.Address()] = assetsPrecompile
	precompiles[delegationPrecompile.Address()] = delegationPrecompile
	precompiles[avsManagerPrecompile.Address()] = avsManagerPrecompile
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

// AddEVMExtensions adds the given precompiles to the list of active precompiles in the EVM parameters
// and to the available precompiles map in the Keeper. This function returns an error if
// the precompiles are invalid or duplicated.
func (k *Keeper) AddEVMExtensions(ctx sdk.Context, precompiles ...vm.PrecompiledContract) error {
	params := k.GetParams(ctx)

	addresses := make([]string, len(precompiles))
	precompilesMap := maps.Clone(k.precompiles)

	for i, precompile := range precompiles {
		// add to active precompiles
		address := precompile.Address()
		addresses[i] = address.String()

		// add to available precompiles, but check for duplicates
		if _, ok := precompilesMap[address]; ok {
			return fmt.Errorf("precompile already registered: %s", address)
		}
		precompilesMap[address] = precompile
	}

	params.ActivePrecompiles = append(params.ActivePrecompiles, addresses...)

	// NOTE: the active precompiles are sorted and validated before setting them
	// in the params
	if err := k.SetParams(ctx, params); err != nil {
		return err
	}

	// update the pointer to the map with the newly added EVM Extensions
	k.precompiles = precompilesMap
	return nil
}

// IsAvailablePrecompile returns true if the given precompile address is contained in the
// EVM keeper's available precompiles map.
func (k Keeper) IsAvailablePrecompile(address common.Address) bool {
	_, ok := k.precompiles[address]
	return ok
}

// GetAvailablePrecompileAddrs returns the list of available precompile addresses.
//
// NOTE: uses index based approach instead of append because it's supposed to be faster.
// Check https://stackoverflow.com/questions/21362950/getting-a-slice-of-keys-from-a-map.
func (k Keeper) GetAvailablePrecompileAddrs() []common.Address {
	addresses := make([]common.Address, len(k.precompiles))
	i := 0

	//#nosec G705 -- two operations in for loop here are fine
	for address := range k.precompiles {
		addresses[i] = address
		i++
	}

	sort.Slice(addresses, func(i, j int) bool {
		return bytes.Compare(addresses[i].Bytes(), addresses[j].Bytes()) == -1
	})

	return addresses
}
