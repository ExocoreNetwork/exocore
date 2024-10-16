package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/ExocoreNetwork/exocore/x/appchain/subscriber/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(string) *cobra.Command {
	// Group dogfood queries under a subcommand
	cmd := &cobra.Command{
		Use: types.ModuleName,
		Short: fmt.Sprintf(
			"Querying commands for the %s module",
			types.ModuleName,
		),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return cmd
}
