package cli

import (
    "strconv"
	
	 "github.com/spf13/cast"
	"github.com/spf13/cobra"
    "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/exocore/x/reward/types"
)

var _ = strconv.Itoa(0)

func CmdRewardDetail() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reward-detail [id]",
		Short: "Broadcast message reward-detail",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
      		 argId, err := cast.ToUint64E(args[0])
            		if err != nil {
                		return err
            		}
            
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgRewardDetail(
				clientCtx.GetFromAddress().String(),
				argId,
				
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

    return cmd
}