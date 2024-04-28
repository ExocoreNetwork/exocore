package cli

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ExocoreNetwork/exocore/x/operator/types"
)

const (
	FlagApproveAddr     = "approve-addr"
	FlagMetaInfo        = "meta-info"
	FlagClientChainData = "client-chain-data"
)

// NewTxCmd returns a root CLI command handler for deposit commands
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Operator transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		CmdRegisterOperator(),
		CmdOptIntoAVS(),
		CmdOptOutOfAVS(),
	)
	return txCmd
}

// CmdRegisterOperator returns a CLI command handler for creating a MsgRegisterOperator
// transaction.
func CmdRegisterOperator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-operator",
		Short: "register to become an operator",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			txf, msg, err := newBuildRegisterOperatorMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			// this calls ValidateBasic internally so we don't need to do that.
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	// EarningsAddr is the same as the sender's address, since the operator registration must be
	// done by the operators themselves.

	f := cmd.Flags()
	// ApproveAddr may be different from the sender's address.
	f.String(
		FlagApproveAddr, "", "The address which is used to approve the delegations made to "+
			"the operator. If not provided, it will default to the sender's address.",
	)
	// OperatorMetaInfo is the name of the operator.
	f.String(
		FlagMetaInfo, "", "The operator's meta info (like name)",
	)
	// clientChainLzID:ClientChainEarningsAddr
	f.StringArray(
		FlagClientChainData, []string{}, "The client chain's address to receive earnings; "+
			"can be supplied multiple times. "+
			"Format: <client-chain-id>:<client-chain-earnings-addr>",
	)
	f.AddFlagSet(stakingcli.FlagSetCommissionCreate())

	// transaction level flags from the SDK
	flags.AddTxFlagsToCmd(cmd)

	// required flags
	_ = cmd.MarkFlagRequired(FlagMetaInfo) // name of the operator

	return cmd
}

func newBuildRegisterOperatorMsg(
	clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet,
) (tx.Factory, *types.RegisterOperatorReq, error) {
	sender := clientCtx.GetFromAddress()
	// #nosec G703 // this only errors if the flag isn't defined.
	approveAddr, _ := fs.GetString(FlagApproveAddr)
	if approveAddr == "" {
		approveAddr = sender.String()
	}
	metaInfo, _ := fs.GetString(FlagMetaInfo)
	msg := &types.RegisterOperatorReq{
		FromAddress: sender.String(),
		Info: &types.OperatorInfo{
			EarningsAddr:     sender.String(),
			ApproveAddr:      approveAddr,
			OperatorMetaInfo: metaInfo,
		},
	}
	clientChainEarningAddress := &types.ClientChainEarningAddrList{}
	// #nosec G703
	ccData, _ := fs.GetStringArray(FlagClientChainData)
	clientChainEarningAddress.EarningInfoList = make(
		[]*types.ClientChainEarningAddrInfo, len(ccData),
	)
	for i, arg := range ccData {
		strList := strings.Split(arg, ":")
		if len(strList) != 2 {
			return txf, nil, errorsmod.Wrapf(
				types.ErrCliCmdInputArg, "the error input arg is:%s", arg,
			)
		}
		// note that this is not the hex value but the decimal number.
		clientChainLzID, err := strconv.ParseUint(strList[0], 10, 64)
		if err != nil {
			return txf, nil, errorsmod.Wrapf(
				types.ErrCliCmdInputArg, "the error input arg is:%s", arg,
			)
		}
		clientChainEarningAddress.EarningInfoList[i] = &types.ClientChainEarningAddrInfo{
			LzClientChainID: clientChainLzID, ClientChainEarningAddr: strList[1],
		}
	}
	msg.Info.ClientChainEarningsAddr = clientChainEarningAddress
	// get the initial commission parameters
	// #nosec G703
	rateStr, _ := fs.GetString(stakingcli.FlagCommissionRate)
	// #nosec G703
	maxRateStr, _ := fs.GetString(stakingcli.FlagCommissionMaxRate)
	// #nosec G703
	maxChangeRateStr, _ := fs.GetString(stakingcli.FlagCommissionMaxChangeRate)
	commission, err := buildCommission(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txf, nil, err
	}
	msg.Info.Commission = commission
	return txf, msg, nil
}

func buildCommission(rateStr, maxRateStr, maxChangeRateStr string) (
	commission stakingtypes.Commission, err error,
) {
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return commission, errorsmod.Wrap(
			types.ErrCliCmdInputArg, "must specify all validator commission parameters",
		)
	}

	rate, err := sdk.NewDecFromStr(rateStr)
	if err != nil {
		return commission, err
	}

	maxRate, err := sdk.NewDecFromStr(maxRateStr)
	if err != nil {
		return commission, err
	}

	maxChangeRate, err := sdk.NewDecFromStr(maxChangeRateStr)
	if err != nil {
		return commission, err
	}

	commission = stakingtypes.NewCommission(rate, maxRate, maxChangeRate)

	return commission, nil
}

// CmdOptIntoAVS returns a CLI command handler for creating a MsgOptIntoAVS transaction.
func CmdOptIntoAVS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opt-into-avs <avs-address>",
		Short: "opt into an AVS by specifying its address or the chain id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.OptIntoAVSReq{
				FromAddress: clientCtx.GetFromAddress().String(),
				AvsAddress:  args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}

func CmdOptOutOfAVS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "opt-out-of-avs <avs-address>",
		Short: "opt out of an AVS by specifying its address or the chain id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.OptOutOfAVSReq{
				FromAddress: clientCtx.GetFromAddress().String(),
				AvsAddress:  args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	return cmd
}
