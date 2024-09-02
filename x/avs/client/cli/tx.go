package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/pflag"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

const (
	FlagOperatorAddress     = "operator-address"
	FlagTaskResponseHash    = "task-response-hash"
	FlagTaskResponse        = "task-response"
	FlagBlsSignature        = "bls-signature"
	FlagTaskContractAddress = "task-contract-address"
	FlagTaskID              = "task-id"
	FlagStage               = "stage"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		CmdSubmitTaskResult(),
	)
	return txCmd
}

// CmdSubmitTaskResult returns a CLI command handler for submit  a TaskResult
// transaction.
func CmdSubmitTaskResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-task-result",
		Short: "submit task result",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			msg := newBuildMsg(clientCtx, cmd.Flags())

			// this calls ValidateBasic internally so we don't need to do that.
			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	f := cmd.Flags()
	f.String(
		FlagOperatorAddress, "", "The address of the operator being queried "+
			" If not provided, it will default to the sender's address.",
	)
	f.String(
		FlagTaskResponseHash, "", "The task response msg hash",
	)
	f.String(
		FlagTaskResponse, "", "The task response data",
	)
	f.String(
		FlagBlsSignature, "", "The operator bls sig info",
	)
	f.String(
		FlagTaskContractAddress, "", "The contract address of task",
	)
	f.Uint64(
		FlagTaskID, 1, "The  task id",
	)
	f.String(
		FlagStage, "", "The stage is a two-stage submission with two values, 1 and 2",
	)
	// #nosec G703 // this only errors if the flag isn't defined.
	_ = cmd.MarkFlagRequired(FlagTaskID)
	_ = cmd.MarkFlagRequired(FlagBlsSignature)
	_ = cmd.MarkFlagRequired(FlagTaskContractAddress)
	_ = cmd.MarkFlagRequired(FlagStage)

	// transaction level flags from the SDK
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func newBuildMsg(
	clientCtx client.Context, fs *pflag.FlagSet,
) *types.SubmitTaskResultReq {
	sender := clientCtx.GetFromAddress()
	operatorAddress, _ := fs.GetString(FlagOperatorAddress)
	if operatorAddress == "" {
		operatorAddress = sender.String()
	}
	taskResponseHash, _ := fs.GetString(FlagTaskResponseHash)

	taskResponse, _ := fs.GetString(FlagTaskResponse)
	taskRes, _ := hex.DecodeString(taskResponse)
	blsSignature, _ := fs.GetString(FlagBlsSignature)
	sig, _ := hex.DecodeString(blsSignature)
	taskContractAddress, _ := fs.GetString(FlagTaskContractAddress)

	taskID, _ := fs.GetUint64(FlagTaskID)
	stage, _ := fs.GetString(FlagStage)

	msg := &types.SubmitTaskResultReq{
		FromAddress: sender.String(),
		Info: &types.TaskResultInfo{
			OperatorAddress:     operatorAddress,
			TaskResponseHash:    taskResponseHash,
			TaskResponse:        taskRes,
			BlsSignature:        sig,
			TaskContractAddress: taskContractAddress,
			TaskId:              taskID,
			Stage:               stage,
		},
	}
	return msg
}
