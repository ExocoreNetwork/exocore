package keeper

import (
	"github.com/ExocoreNetwork/exocore/x/dogfood/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// interface guard
var _ types.DogfoodHooks = &MultiDogfoodHooks{}

// MultiDogfoodHooks is a collection of DogfoodHooks. It calls the hook for each element in the
// collection one-by-one. The hook is called in the order in which the collection is created.
type MultiDogfoodHooks []types.DogfoodHooks

// NewMultiDogfoodHooks is used to create a collective object of dogfood hooks from a list of
// the hooks. It follows the "accept interface, return concrete types" philosophy. Other modules
// may set the hooks by calling k := (*k).SetHooks(NewMultiDogfoodHooks(hookI))
func NewMultiDogfoodHooks(hooks ...types.DogfoodHooks) MultiDogfoodHooks {
	return hooks
}

// AfterValidatorBonded is the implementation of types.DogfoodHooks for MultiDogfoodHooks.
func (hooks MultiDogfoodHooks) AfterValidatorBonded(
	ctx sdk.Context,
	consAddr sdk.ConsAddress,
	operator sdk.ValAddress,
) error {
	for _, hook := range hooks {
		if err := hook.AfterValidatorBonded(ctx, consAddr, operator); err != nil {
			return err
		}
	}
	return nil
}
