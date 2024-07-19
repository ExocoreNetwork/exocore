package cli

import (
	"fmt"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	return cmd
}

func RegisterAVS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "RegisterAVS: AvsName, AvsAddress, OperatorAddress, AvsOwnerAddress, AssetId",
		Short: "register to be an avs",
		Args:  cobra.MinimumNArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress()
			fromAddress := sender.String()
			msg := &types.RegisterAVSReq{
				FromAddress: fromAddress,
				Info: &types.AVSInfo{
					Name:            args[0],
					AvsAddress:      args[1],
					SlashAddr:       args[2],
					AvsOwnerAddress: []string{args[3]},
					AssetId:         []string{args[4]},
				},
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
