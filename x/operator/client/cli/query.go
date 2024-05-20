package cli

import (
	"context"

	operatortypes "github.com/ExocoreNetwork/exocore/x/operator/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the parent command for all incentives CLI query commands.
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        operatortypes.ModuleName,
		Short:                      "Querying commands for the operator module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetOperatorInfo(),
		QueryOperatorUSDValue(),
		QueryAVSUSDValue(),
		QueryOperatorSlashInfo(),
	)
	return cmd
}

// GetOperatorInfo queries operator info
func GetOperatorInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "GetOperatorInfo operatorAddr",
		Short: "Get operator info",
		Long:  "Get operator info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := operatortypes.NewQueryClient(clientCtx)
			req := &operatortypes.GetOperatorInfoReq{
				OperatorAddr: args[0],
			}
			res, err := queryClient.GetOperatorInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueryOperatorUSDValue queries the opted-in USD value for the operator
func QueryOperatorUSDValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryOperatorUSDValue operatorAddr avsAddr",
		Short: "Get the opted-in USD value for the operator",
		Long:  "Get the opted-in USD value for the operator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := operatortypes.NewQueryClient(clientCtx)
			req := &operatortypes.QueryOperatorUSDValueRequest{
				OperatorAddr: args[0],
				AVSAddress:   args[1],
			}
			res, err := queryClient.QueryOperatorUSDValue(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueryAVSUSDValue queries the USD value for the avs
func QueryAVSUSDValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryAVSUSDValue avsAddr",
		Short: "Get the USD value for the avs",
		Long:  "Get the USD value for the avs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := operatortypes.NewQueryClient(clientCtx)
			req := &operatortypes.QueryAVSUSDValueRequest{
				AVSAddress: args[0],
			}
			res, err := queryClient.QueryAVSUSDValue(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueryOperatorSlashInfo queries the slash information for the specified operator and AVS
func QueryOperatorSlashInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryOperatorUSDValue operatorAddr avsAddr",
		Short: "Get the the slash information for the operator",
		Long:  "Get the the slash information for the operator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := operatortypes.NewQueryClient(clientCtx)
			req := &operatortypes.QueryOperatorSlashInfoRequest{
				OperatorAddr: args[0],
				AVSAddress:   args[1],
			}
			res, err := queryClient.QueryOperatorSlashInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
