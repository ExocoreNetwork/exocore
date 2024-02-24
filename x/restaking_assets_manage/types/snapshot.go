package types

import sdkmath "cosmossdk.io/math"

func (m *Snapshot) UpdateForDelegation(
	assetId string,
	stakerId string,
	changeAmount sdkmath.Int,
) {
	// find the asset within the snapshot
	perAssetId, ok := m.PerAssetId[assetId]
	if !ok {
		// if it doesn't exist, create one
		perAssetId = SnapshotPerAssetId{}
	}
	perStaker, ok := perAssetId.PerStaker[stakerId]
	if !ok {
		// if it doesn't exist, create one
		perStaker = SnapshotPerAssetIdPerStaker{}
	}
	perStaker.Delegated = perStaker.Delegated.Add(changeAmount)
	// update the snapshot
	perAssetId.PerStaker[stakerId] = perStaker
	m.PerAssetId[assetId] = perAssetId
}

// UpdateForUndelegation updates the snapshot for an undelegation
// it performs no error checking and assumes that the undelegation is valid
// and that the undelegation amount is not greater than the delegated amount
// if the amount undelegated is equal to the delegated amount, it performs clean up too
func (m *Snapshot) UpdateForUndelegation(
	assetId string, stakerId string,
	changeAmount sdkmath.Int,
) {
	x := m.PerAssetId[assetId].PerStaker[stakerId]
	x.Delegated = x.Delegated.Sub(changeAmount)
	if x.Delegated.IsZero() {
		delete(m.PerAssetId[assetId].PerStaker, stakerId)
		if len(m.PerAssetId[assetId].PerStaker) == 0 {
			delete(m.PerAssetId, assetId)
		}
	} else {
		m.PerAssetId[assetId].PerStaker[stakerId] = x
	}
}
