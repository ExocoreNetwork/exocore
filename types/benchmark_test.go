package types_test

import (
	"fmt"
	"testing"

	"github.com/evmos/evmos/v14/types"
)

func BenchmarkParseChainID(b *testing.B) {
	b.ReportAllocs()
	// Start at 1, for valid EIP155, see regexEIP155 variable.
	for i := 1; i < b.N; i++ {
		chainID := fmt.Sprintf("evmos_1-%d", i)
		if _, err := types.ParseChainID(chainID); err != nil {
			b.Fatal(err)
		}
	}
}
