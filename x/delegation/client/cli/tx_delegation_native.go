package cli

import (
	"errors"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/ExocoreNetwork/exocore/x/delegation/types"
)

func CmdDelagateNativeToken() *cobra.Command {
	cmd := &cobra.Command{
		// TODO: only support native token for now
		Use:   "delegate-native operator amount",
		Short: "Broadcast message delegate-native to delegate native token",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, ok := sdkmath.NewIntFromString(args[1])
			if !ok || amount.IsNegative() {
				return errors.New("amount invalid")
			}

			operatorAddrStr := args[0]

			msg := types.NewMsgDelegation(clientCtx.GetFromAddress().String(), map[string]sdkmath.Int{operatorAddrStr: amount})

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
