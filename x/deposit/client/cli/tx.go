// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package cli

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	types2 "github.com/exocore/x/deposit/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

// NewTxCmd returns a root CLI command handler for deposit commands
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types2.ModuleName,
		Short:                      "deposit subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		UpdateParams(),
	)
	return txCmd
}

// UpdateParams todo: it should be a gov proposal command in future.
func UpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "UpdateParams ExoCoreLZAppAddr ExoCoreLzAppEventTopic",
		Short: "set ExoCoreLZAppAddr and ExoCoreLzAppEventTopic params to deposit module",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress()
			msg := &types2.MsgUpdateParams{
				Authority: sender.String(),
				Params: types2.Params{
					ExoCoreLzAppAddress:    args[0],
					ExoCoreLzAppEventTopic: args[1],
				},
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
