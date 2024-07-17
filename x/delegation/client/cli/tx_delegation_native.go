package cli

import (
	"errors"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	assetstypes "github.com/ExocoreNetwork/exocore/x/assets/types"
	"github.com/ExocoreNetwork/exocore/x/delegation/types"
)

func CmdDelegate() *cobra.Command {
	cmd := &cobra.Command{
		// TODO: only support native token for now
		Use:   "delegate asset-id operator amount approve-signature, approve-salt",
		Short: "Broadcast a transaction to delegate amount of native token to the operator",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetID := args[0]
			approveSignature, approveSalt := "", ""
			if assetID != assetstypes.NativeAssetID {
				approveSignature = args[3]
				approveSalt = args[4]
			}

			operatorAddrStr := args[1]

			amount, ok := sdkmath.NewIntFromString(args[2])
			if !ok || amount.LTE(sdkmath.ZeroInt()) {
				return errors.New("amount invalid")
			}
			msg := types.NewMsgDelegation(assetID, clientCtx.GetFromAddress().String(), approveSignature, approveSalt, []types.KeyValue{{Key: operatorAddrStr, Value: &types.ValueField{Amount: amount}}})

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdUndelegate() *cobra.Command {
	cmd := &cobra.Command{
		// TODO: only support native token for now
		Use:   "undelegate asset-id operator amount",
		Short: "Broadcast a transaction to undelegate amount of native token from the operator",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			assetID := args[0]

			operatorAddrStr := args[1]

			amount, ok := sdkmath.NewIntFromString(args[2])
			if !ok || amount.LTE(sdkmath.ZeroInt()) {
				return errors.New("amount invalid")
			}
			msg := types.NewMsgUndelegation(assetID, clientCtx.GetFromAddress().String(), []types.KeyValue{{Key: operatorAddrStr, Value: &types.ValueField{Amount: amount}}})

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
