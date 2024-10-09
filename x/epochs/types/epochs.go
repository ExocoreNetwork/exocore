package types

// NewEpoch creates a new Epoch instance.
func NewEpoch(number uint64, identifier string) Epoch {
	return Epoch{
		EpochNumber:     number,
		EpochIdentifier: identifier,
	}
}
