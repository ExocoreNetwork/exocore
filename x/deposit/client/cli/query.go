package cli

import (
	deposittype "github.com/ExocoreNetwork/exocore/x/deposit/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the parent command for all incentives CLI query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        deposittype.ModuleName,
		Short:                      "Querying commands for the deposit module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand()
	return cmd
}
