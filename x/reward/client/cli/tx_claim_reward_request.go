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

func CmdClaimRewardRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim-reward-request [id] [rewardaddress]",
		Short: "Broadcast message claim-reward-request",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
      		 argId, err := cast.ToUint64E(args[0])
            		if err != nil {
                		return err
            		}
             argRewardaddress := args[1]
            
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaimRewardRequest(
				clientCtx.GetFromAddress().String(),
				argId,
				argRewardaddress,
				
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