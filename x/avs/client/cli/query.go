package cli

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"golang.org/x/xerrors"

	"github.com/ExocoreNetwork/exocore/x/avs/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group avs queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(QueryAVSInfo())
	cmd.AddCommand(QueryAVSAddrByChainID())
	cmd.AddCommand(QueryTaskInfo())
	return cmd
}

func QueryAVSInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "AVSInfo query",
		Short: "AVSInfo query",
		Long:  "AVSInfo query for current registered AVS",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !common.IsHexAddress(args[0]) {
				return xerrors.Errorf("invalid avs  address,err:%s", types.ErrInvalidAddr)
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryAVSInfoReq{
				AVSAddress: args[0],
			}
			res, err := queryClient.QueryAVSInfo(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// QueryAVSAddrByChainID returns a command to query AVS address by chainID
func QueryAVSAddrByChainID() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "AVSAddrByChainID <chainID>",
		Short:   "AVSAddrByChainID <chainID>",
		Long:    "AVSAddrByChainID query for AVS address by chainID",
		Example: "exocored query avs AVSAddrByChainID exocoretestnet_233-1",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := &types.QueryAVSAddrByChainIDReq{
				ChainID: args[0],
			}
			res, err := queryClient.QueryAVSAddrByChainID(context.Background(), req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func QueryTaskInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "TaskInfo query",
		Short: "TaskInfo query",
		Long:  "TaskInfo query for current registered task",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !common.IsHexAddress(args[0]) {
				return xerrors.Errorf("invalid task  address,err:%s", types.ErrInvalidAddr)
			}
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			req := types.QueryAVSTaskInfoReq{
				TaskAddr: args[0],
				TaskId:   args[1],
			}
			res, err := queryClient.QueryAVSTaskInfo(context.Background(), &req)
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
