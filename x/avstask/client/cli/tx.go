package cli

import (
	"fmt"

	avstasktypes "github.com/ExocoreNetwork/exocore/x/avstask/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        avstasktypes.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", avstasktypes.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		RegisterAVSTask(),
	)
	return cmd
}

func RegisterAVSTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Registertask params to taskManager module",
		Short: "Registertask params to taskManager module",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}
			sender := cliCtx.GetFromAddress().String()

			msg := &avstasktypes.RegisterAVSTaskReq{
				FromAddress: sender,
				Task: &avstasktypes.TaskContractInfo{
					TaskContractAddress: args[0],
					Name:                args[1],
					MetaInfo:            args[2],
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
