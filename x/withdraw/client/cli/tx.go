package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	paramstypes "github.com/exocore/x/deposit/types"
	"github.com/exocore/x/withdraw/types"
)

var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        paramstypes.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", paramstypes.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		UpdateParams(),
	)

	return cmd
}

// UpdateParams todo: it should be a gov proposal command in future.
func UpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "UpdateParams ExoCoreLZAppAddr ExoCoreLzAppEventTopic",
		Short: "set ExoCoreLZAppAddr and ExoCoreLzAppEventTopic params to withdraw module",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress()
			msg := &types.MsgUpdateParams{
				Authority: sender.String(),
				Params: paramstypes.Params{
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
