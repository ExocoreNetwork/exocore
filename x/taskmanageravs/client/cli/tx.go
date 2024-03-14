package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ExocoreNetwork/exocore/x/taskmanageravs/types"
	taskTypes "github.com/ExocoreNetwork/exocore/x/taskmanageravs/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
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

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		CreateTask(),
	)
	return cmd
}

func CreateTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Set task params to taskManager module",
		Short: "Set task params to taskManager module",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			taskId, err := strconv.Atoi(args[2])
			if err != nil {
				return err
			}

			quorumThresholdPercentage, err := strconv.ParseUint(args[3], 10, 32)
			if err != nil {
				return err
			}
			task := &taskTypes.TaskInfo{
				Name:                args[0],
				MetaInfo:            args[1],
				TaskId:              uint64(taskId),
				ThresholdPercentage: quorumThresholdPercentage,
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), task)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
