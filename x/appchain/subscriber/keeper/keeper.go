package keeper

import (
	"fmt"
	"time"

	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
)

type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         storetypes.StoreKey
	scopedKeeper     commontypes.ScopedKeeper
	portKeeper       commontypes.PortKeeper
	clientKeeper     commontypes.ClientKeeper
	connectionKeeper commontypes.ConnectionKeeper
	channelKeeper    commontypes.ChannelKeeper
	ibcCoreKeeper    commontypes.IBCCoreKeeper
}

// NewKeeper creates a new subscriber keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	scopedKeeper commontypes.ScopedKeeper,
	portKeeper commontypes.PortKeeper,
	clientKeeper commontypes.ClientKeeper,
	connectionKeeper commontypes.ConnectionKeeper,
	channelKeeper commontypes.ChannelKeeper,
	ibcCoreKeeper commontypes.IBCCoreKeeper,
) Keeper {
	return Keeper{
		cdc:              cdc,
		storeKey:         storeKey,
		scopedKeeper:     scopedKeeper,
		portKeeper:       portKeeper,
		clientKeeper:     clientKeeper,
		connectionKeeper: connectionKeeper,
		channelKeeper:    channelKeeper,
		ibcCoreKeeper:    ibcCoreKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetPort returns the portID for the IBC app module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey()))
}

// SetPort sets the portID for the IBC app module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey(), []byte(portID))
}

// IsBound checks if the IBC app module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	capability := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, capability, host.PortPath(portID))
}

// ClaimCapability allows the IBC app module to claim a capability that core IBC
// passes to it
func (k Keeper) ClaimCapability(
	ctx sdk.Context,
	cap *capabilitytypes.Capability,
	name string,
) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

// GetPendingChanges gets the pending validator set changes that will be applied
// at the end of this block.
func (k Keeper) GetPendingChanges(
	ctx sdk.Context,
) (*commontypes.ValidatorSetChangePacketData, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PendingChangesKey())
	if bz == nil {
		return nil, false
	}
	res := &commontypes.ValidatorSetChangePacketData{}
	k.cdc.MustUnmarshal(bz, res)
	return res, true
}

// SetPendingChanges sets the pending validator set changes that will be applied
// at the end of this block.
func (k Keeper) SetPendingChanges(
	ctx sdk.Context,
	data *commontypes.ValidatorSetChangePacketData,
) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(data)
	store.Set(types.PendingChangesKey(), bz)
}

// SetPacketMaturityTime sets the maturity time for a given received VSC packet id
func (k Keeper) SetPacketMaturityTime(
	ctx sdk.Context, vscID uint64, maturityTime time.Time,
) {
	store := ctx.KVStore(k.storeKey)
	maturingVSCPacket := &types.MaturingVSCPacket{
		ValidatorSetChangeID: vscID,
		MaturityTime:         maturityTime,
	}
	store.Set(
		types.PacketMaturityTimeKey(vscID, maturityTime),
		k.cdc.MustMarshal(maturingVSCPacket),
	)
}

// DeleteOutstandingDowntime deletes the outstanding downtime flag for the given validator
// consensus address
func (k Keeper) DeleteOutstandingDowntime(
	ctx sdk.Context, consAddress sdk.ConsAddress,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.OutstandingDowntimeKey(consAddress))
}
