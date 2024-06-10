package cli

import (
	"context"
	"fmt"

	avstasktypes "github.com/ExocoreNetwork/exocore/x/avstask/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        avstasktypes.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", avstasktypes.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetTaskInfo())
	return cmd
}

// GetTaskInfo queries operator info
func GetTaskInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "GetTaskInfo avstask info",
		Short: "Get avstask info",
		Long:  "Get avstask info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := avstasktypes.NewQueryClient(clientCtx)
			req := &avstasktypes.GetAVSTaskInfoReq{
				TaskAddr: args[0],
			}
			res, err := queryClient.QueryAVSTaskInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
