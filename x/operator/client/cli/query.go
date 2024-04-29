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
		GetOperatorConsAddress(),
		GetAllOperatorKeys(),
		GetAllOperatorConsAddrs(),
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
			res, err := queryClient.QueryOperatorInfo(context.Background(), req)
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
				Chain: args[0],
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

// GetAllOperatorKeys queries all operators for the provided chain ID and their
// consensus keys
func GetAllOperatorKeys() *cobra.Command {
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
			req := &operatortypes.QueryAllOperatorConsKeysByChainIDRequest{
				Chain:      args[0],
				Pagination: pageReq,
			}
			res, err := queryClient.QueryAllOperatorConsKeysByChainID(
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

// GetOperatorConsAddress queries operator consensus address for the provided chain ID
func GetOperatorConsAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-operator-cons-address <chain_id>",
		Short: "Get operator consensus address",
		Long:  "Get operator consensus address for the provided chain ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := operatortypes.NewQueryClient(clientCtx)
			req := &operatortypes.QueryOperatorConsAddressRequest{
				Chain: args[0],
			}
			res, err := queryClient.QueryOperatorConsAddressForChainID(
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

// GetAllOperatorConsAddrs queries all operators for the provided chain ID and their
// consensus addresses
func GetAllOperatorConsAddrs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all-operator-cons-addrs <chain_id>",
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
			req := &operatortypes.QueryAllOperatorConsAddrsByChainIDRequest{
				Chain:      args[0],
				Pagination: pageReq,
			}
			res, err := queryClient.QueryAllOperatorConsAddrsByChainID(
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
