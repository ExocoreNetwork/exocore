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
		GetOperatorConsKey(),
		GetAllOperatorsByChainID(),
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

// GetOperatorConsKey queries operator consensus key for the provided chain ID
func GetOperatorConsKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-operator-cons-key <chain_id>",
		Short: "Get operator consensus key",
		Long:  "Get operator consensus key for the provided chain ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := operatortypes.NewQueryClient(clientCtx)
			req := &operatortypes.QueryOperatorConsKeyRequest{
				ChainID: args[0],
			}
			res, err := queryClient.QueryOperatorConsKeyForChainID(
				context.Background(), req,
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetAllOperatorsByChainID queries all operators for the provided chain ID
func GetAllOperatorsByChainID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all-operators-by-chain-id <chain_id>",
		Short: "Get all operators for the provided chain ID",
		Long:  "Get all operators for the provided chain ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := operatortypes.NewQueryClient(clientCtx)
			req := &operatortypes.QueryAllOperatorsByChainIDRequest{
				ChainID:    args[0],
				Pagination: pageReq,
			}
			res, err := queryClient.QueryAllOperatorKeysByChainID(
				context.Background(), req,
			)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
