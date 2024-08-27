package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultIBCTimeoutPeriod is the default timeout period for IBC packets.
	DefaultIBCTimeoutPeriod = 4 * 7 * 24 * time.Hour
)

// CalculateTrustPeriod calculates the trust period as the unbonding period multiplied by the trust period fraction, truncated to an integer.
func CalculateTrustPeriod(unbondingPeriod time.Duration, trustPeriodFraction string) (time.Duration, error) {
	trustDec, err := sdk.NewDecFromStr(trustPeriodFraction)
	if err != nil {
		return time.Duration(0), err
	}
	trustPeriod := time.Duration(trustDec.MulInt64(unbondingPeriod.Nanoseconds()).TruncateInt64())
	return trustPeriod, nil
}

// ValidateStringFraction validates that the given string is a valid fraction in the range [0, 1].
func ValidateStringFraction(str string) error {
	dec, err := sdk.NewDecFromStr(str)
	if err != nil {
		return err
	}
	if dec.IsNegative() {
		return fmt.Errorf("param cannot be negative, got %s", str)
	}
	if dec.Sub(sdk.NewDec(1)).IsPositive() {
		return fmt.Errorf("param cannot be greater than 1, got %s", str)
	}
	return nil
}

// ValidateDuration validates that the given duration is positive.
func ValidateDuration(period time.Duration) error {
	if period <= time.Duration(0) {
		return fmt.Errorf("duration must be positive")
	}
	return nil
}
