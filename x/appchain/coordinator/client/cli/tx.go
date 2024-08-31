package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ExocoreNetwork/exocore/x/appchain/coordinator/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: types.ModuleName,
		Short: fmt.Sprintf(
			"%s transactions subcommands",
			types.ModuleName,
		),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdRegisterSubscriberChain(),
	)
	return cmd
}

// CmdRegisterSubscriberChain returns the command to register a subscriber chain.
func CmdRegisterSubscriberChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-subscriber-chain [json-args]",
		Short: "Register a subscriber chain",
		Long:  "Register a subscriber chain within the appchain coordinator module",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if !json.Valid([]byte(args[0])) {
				return fmt.Errorf("invalid JSON argument: %s", args[0])
			}

			msg := types.NewRegisterSubscriberChainRequest(clientCtx.GetFromAddress().String(), args[0])

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
