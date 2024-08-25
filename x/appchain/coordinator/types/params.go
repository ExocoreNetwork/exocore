package types

// DefaultParams returns the default parameters for the module.
func DefaultParams() Params {
	return Params{}
}

// NewParams creates a new Params object
func NewParams() Params {
	return Params{}
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	return nil
}
