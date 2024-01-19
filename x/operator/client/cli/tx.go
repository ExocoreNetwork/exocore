package cli

import (
	"fmt"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client/flags"
	operatortypes "github.com/exocore/x/operator/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

// NewTxCmd returns a root CLI command handler for deposit commands
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        operatortypes.ModuleName,
		Short:                      "delegation subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		RegisterOperator(),
	)
	return txCmd
}

// RegisterOperator register to be a operator
func RegisterOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "RegisterOperator EarningsAddr ApproveAddr OperatorMetaInfo clientChainLzID:ClientChainEarningsAddr",
		Short: "register to be a operator",
		Args:  cobra.MinimumNArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress()
			msg := &operatortypes.RegisterOperatorReq{
				FromAddress: sender.String(),
				Info: &operatortypes.OperatorInfo{
					EarningsAddr:     args[0],
					ApproveAddr:      args[1],
					OperatorMetaInfo: args[2],
				},
			}
			lastArgs := args[3:]
			clientChainEarningAddress := &operatortypes.ClientChainEarningAddrList{}
			clientChainEarningAddress.EarningInfoList = make([]*operatortypes.ClientChainEarningAddrInfo, 0)
			for _, arg := range lastArgs {
				strList := strings.Split(arg, ":")
				if len(strList) != 2 {
					return errorsmod.Wrap(operatortypes.ErrCliCmdInputArg, fmt.Sprintf("the error input arg is:%s", arg))
				}
				clientChainLzId, err := strconv.ParseUint(strList[0], 10, 64)
				if err != nil {
					return err
				}
				clientChainEarningAddress.EarningInfoList = append(clientChainEarningAddress.EarningInfoList,
					&operatortypes.ClientChainEarningAddrInfo{
						LzClientChainId: clientChainLzId, ClientChainEarningAddr: strList[1],
					})
			}
			msg.Info.ClientChainEarningsAddr = clientChainEarningAddress
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
