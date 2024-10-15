package types

import (
	commontypes "github.com/ExocoreNetwork/exocore/x/appchain/common/types"
)

// SubscriberPacketDataWithIdx is a wrapper struct for SubscriberPacketData with an index.
type SubscriberPacketDataWithIdx struct {
	commontypes.SubscriberPacketData
	Idx uint64
}
